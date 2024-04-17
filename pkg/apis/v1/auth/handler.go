package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	v1 "v1/pkg/apis/v1"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/apiserver/imsystem"
	"v1/pkg/apiserver/request"
	"v1/pkg/captcha"
	"v1/pkg/client/cache"
	"v1/pkg/dao"
	"v1/pkg/model"
	"v1/pkg/server/errutil"
	"v1/pkg/token"
	"v1/pkg/utils"
)

type authHandlerOption struct {
	tokenManager token.Manager
	db           *gorm.DB
	cacheClient  cache.Interface
}

type authHandler struct {
	authHandlerOption
}

func newAuthHandler(option authHandlerOption) *authHandler {
	return &authHandler{
		authHandlerOption: option,
	}
}

func (h *authHandler) login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := loginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Account == "" || req.Password == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	captchaService := captcha.GetService()

	if !captchaService.VerifyCaptcha(req.CaptchaID, strings.ToLower(req.CaptchaValue)) {
		zap.L().Error("captcha value is wrong")
		encoding.HandleError(c, errutil.NewError(400, "captcha value is wrong"))
		return
	}

	u, err := dao.GetUserByUsername(ctx, h.db, req.Account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			encoding.HandleError(c, errutil.NewError(http.StatusBadRequest, "用户不存在"))
			return
		}

		encoding.HandleError(c, err)
		return
	}

	if u.Status == model.UserStatusDisabled {
		encoding.HandleError(c, errutil.NewError(http.StatusBadRequest, "用户已经被停用，请联系管理员"))
		return
	}

	if u.Password != utils.MD5Hex(req.Password) {
		encoding.HandleError(c, errutil.NewError(http.StatusBadRequest, "密码错误"))
		return
	}

	// 签发token
	tokenStr, err := h.tokenManager.IssueTo(token.Payload{
		ID:       u.ID,
		UID:      u.UID,
		Username: u.Username,
		Name:     u.Name,
		Role:     u.Role,
	}, time.Hour*24)
	if err != nil {
		encoding.HandleSuccess(c, err)
		return
	}

	if err = h.cacheClient.Set(ctx, "token:"+tokenStr, u.Username, time.Minute*30); err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleSuccess(c, errutil.ErrInternalServer)
		return
	}

	err = imsystem.AddChatClient(imsystem.ChatServerIp, imsystem.ChatServerPort, u.UID)
	if err != nil {
		zap.L().Error(" imsystem.AddChatClient ", zap.Error(err))
		encoding.HandleSuccess(c, errutil.ErrInternalServer)
		return
	}

	result := loginResp{
		ID:       u.ID,
		UID:      u.UID,
		Username: u.Username,
		Role:     u.Role,
		Name:     u.Name,
		Token:    tokenStr,
		Head:     u.Head,
	}
	encoding.HandleSuccess(c, &result)
}

func (h *authHandler) createCaptcha(c *gin.Context) {
	if !utils.Lmt.AllowKey(utils.MD5Hex(c.Request.UserAgent())) {
		zap.L().Error("The captcha request is too fast. Please try again later")
		encoding.HandleError(c, errutil.NewError(http.StatusBadRequest, "The captcha request is too fast"))
		return
	}

	serveice := captcha.GetService()
	captchaId, captchaValue, answer, err := serveice.CreateCaptcha()
	if err != nil {
		zap.L().Error("create captcha failed", zap.Error(err))
		encoding.HandleError(c, errutil.NewError(400, "create captcha failed"))
		return
	}

	fmt.Println("answer:", answer)

	encoding.HandleSuccess(c, createCaptchaResp{captchaId, captchaValue})
}

func (h *authHandler) logout(c *gin.Context) {

	imsystem.DeleteClient(request.GetUserUIDFromCtx(c))
	encoding.HandleSuccess(c)
}
