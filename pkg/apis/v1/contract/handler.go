package contract

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "v1/pkg/apis/v1"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/contract"
	"v1/pkg/dao"
	"v1/pkg/server/errutil"
)

type contractHandlerOption struct {
	db          *gorm.DB
	bolckClient *contract.Client
}

type contractHandler struct {
	contractHandlerOption
}

func newInterviewHandler(option contractHandlerOption) *contractHandler {
	return &contractHandler{
		contractHandlerOption: option,
	}
}

func (h *contractHandler) blockList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := blockListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	count, list, err := dao.GetBlocks(ctx, h.db, req.BlodkHash, req.SaveType, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetBlocks", zap.Error(err))
		encoding.HandleError(c, errors.New("get log failed"))
		return
	}

	encoding.HandleSuccess(c, blockListResp{Count: count, Items: list})
}

func (h *contractHandler) getContent(c *gin.Context) {
	_, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getContentReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if h.bolckClient == nil {
		encoding.HandleError(c, errors.New("合约未部署成功，无法连接私链"))
		return
	}

	value, _, err := h.bolckClient.GetValue(req.KeyHash)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errors.New("获取内容失败"))
		return
	}

	encoding.HandleSuccess(c, getContentResp{Value: value})
}

func (h *contractHandler) blockDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	hash := c.Param("block_hash")

	log, err := dao.GetBlockByBlockHash(ctx, h.db, hash)
	if err != nil {
		zap.L().Error("dao.GetBlockByBlockHash", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到对应日志"))
		return
	}

	encoding.HandleSuccess(c, log)
}
