package interview

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
	v1 "v1/pkg/apis/v1"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/apiserver/request"
	"v1/pkg/dao"
	"v1/pkg/model"
	"v1/pkg/server/errutil"
)

type interviewHandlerOption struct {
	db *gorm.DB
}

type interviewHandler struct {
	interviewHandlerOption
}

func newInterviewHandler(option interviewHandlerOption) *interviewHandler {
	return &interviewHandler{
		interviewHandlerOption: option,
	}
}

func (h *interviewHandler) createInterview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := createInterviewReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// get creator
	_, creator, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error(" dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if creator.Role != model.RoleTypeFirm {
		zap.L().Error("user role is not for firm")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	// get interviewee info
	_, interviewee, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error(" dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	//
	if interviewee.Role != model.RoleTypeStudent {
		zap.L().Error("role not student")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	info := model.InterviewInfo{
		Title:       req.Title,
		Content:     req.Content,
		Date:        time.UnixMilli(req.Date).Format(time.DateTime),
		Location:    req.Location,
		Position:    req.Position,
		Creator:     creator.Username,
		ContactInfo: creator.Phone,
		Interviewee: interviewee.Username,
	}

	_, err = dao.InsertInterview(ctx, h.db, model.Interview{
		Ttile:          req.Title,
		Info:           info,
		Interviewee:    interviewee.Username,
		IntervieweeUID: interviewee.UID,
		Creator:        creator.Username,
		CreatorUID:     creator.UID,
	})
	if err != nil {
		zap.L().Error("dao.InsertInterview", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *interviewHandler) deleteInterview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := deleteInterviewReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	interview, err := dao.GetInterviewByID(ctx, h.db, req.ID)
	if err != nil {
		zap.L().Error("dao.GetInterviewByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = dao.DeleteInterviewByID(ctx, h.db, interview.ID)
	if err != nil {
		zap.L().Error("dao.DeleteInterviewByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	encoding.HandleSuccess(c, "success")
}
