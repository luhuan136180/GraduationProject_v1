package interview

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
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

	// 暂时屏蔽 用户角色导致的 权限不足
	// if creator.Role != model.RoleTypeFirm {
	// 	zap.L().Error("user role is not for firm")
	// 	encoding.HandleError(c, errutil.ErrPermissionDenied)
	// 	return
	// }

	// get interviewee info
	_, interviewee, err := dao.GetUserByUID(ctx, h.db, req.IntervieweeUid)
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

func (h *interviewHandler) interviewList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := interviewListReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 10 {
		req.Size = 10
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error(" dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	var interviews []model.Interview
	var count int64
	if req.IntervieweeUID != "" || req.CreatorUID != "" {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:          req.Title,
			IntervieweeUID: req.IntervieweeUID,
			CreatorUID:     req.CreatorUID,
		})
	} else if user.Role == model.RoleTypeStudent {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:          req.Title,
			IntervieweeUID: user.UID,
		})
	} else if user.Role == model.RoleTypeFirm {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:      req.Title,
			CreatorUID: user.UID,
		})
	}
	if err != nil {
		zap.L().Error(" dao.FindInterviewByOption", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, interviewListResp{Total: count, Data: interviews})
}

func (h interviewHandler) interviewChangeStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := interviewChangeStatusRep{}
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

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if interview.CreatorUID != user.UID && interview.IntervieweeUID != user.UID {
		zap.L().Error("permission denied")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	err = dao.UpdateInterviewStatus(ctx, h.db, req.ID, req.Status)
	if err != nil {
		zap.L().Error("dao.UpdateInterviewStatus", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, interviewChangeStatusResp{ID: interview.ID, Title: interview.Ttile, info: interview.Info, Interviewee: interview.Interviewee, Status: req.Status})
}

func (h *interviewHandler) interviewDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	idStr := c.Param("id")

	if idStr == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("permission denied")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	interview, err := dao.GetInterviewByID(ctx, h.db, id)
	if err != nil {
		zap.L().Error("dao.GetInterviewByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, interviewDetailResp{ID: interview.ID, Title: interview.Ttile, Info: interview.Info, Interviewee: interview.Interviewee, Status: interview.Status})
}
