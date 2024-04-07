package resume

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	v1 "v1/pkg/apis/v1"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/apiserver/request"
	"v1/pkg/dao"
	"v1/pkg/model"
	"v1/pkg/server/errutil"
)

type resumeHandlerOption struct {
	db *gorm.DB
}

type resumeHandler struct {
	resumeHandlerOption
}

func newResumeHandler(option resumeHandlerOption) *resumeHandler {
	return &resumeHandler{
		resumeHandlerOption: option,
	}
}

func (h *resumeHandler) createResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := createReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error(" dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	found, _, err := dao.FoundResumesByname(ctx, h.db, req.ResumeName)
	if err != nil {
		zap.L().Error(" dao.GetResumesByUserUid", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	if found {
		zap.L().Error(" the name is exit", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	_, err = dao.InsertResume(ctx, h.db, model.Resume{
		UserUid:    user.UID,
		UserName:   user.Username,
		ResumeName: req.ResumeName,
		BasicInfo:  req.ResumeInfo,
		ProjectIDs: req.ProjectIDs,
		Creator:    user.Username,
	})
	if err != nil {
		zap.L().Error("dao.InsertResume", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *resumeHandler) deleteResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := deleteResumeReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, resume, err := dao.GetResumeByID(ctx, h.db, req.ResumeID)
	if err != nil {
		zap.L().Error("dao.GetResumeByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if request.GetUsernameFromCtx(ctx) != resume.UserName {
		zap.L().Error("this resume is not create by account")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	err = dao.DeleteResumeByID(ctx, h.db, req.ResumeID)
	if err != nil {
		zap.L().Error("dao.DeleteResumeByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *resumeHandler) resumeList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := resumeListReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 10 {
		req.Size = 10
	}

	user, err := dao.GetUserByUsername(ctx, h.db, request.GetUsernameFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUsername", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	resumeList := make([]model.Resume, 0)
	var count int64

	db := h.db.Model(&model.Resume{})
	if request.GetRoleTypeFromCtx(ctx) == model.RoleTypeStudent {
		db = db.Where("creator = ?", user.Username)
	} else {
		db = db.Where("user_name = (?)", req.UserName)
	}

	if req.Title != "" {
		db = db.Where("resume_name = ?", req.Title)
	}

	if req.StartAt != "" {
		db = db.Where("created_at > ?", req.StartAt)
	}

	if req.EndAt != "" {
		db = db.Where("created_at < ?", req.EndAt)
	}

	err = db.Count(&count).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = db.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&resumeList).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, resumeListResp{Count: count, ResumeList: resumeList})
}

func (h *resumeHandler) resumeDetail(c *gin.Context) {
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

	_, resume, err := dao.GetResumeByID(ctx, h.db, id)
	if err != nil {
		zap.L().Error("dao.GetInterviewByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	data := resumeDetailResp{
		ID:              resume.ID,
		UserUid:         resume.UserUid,
		UserName:        resume.UserName,
		ResumeName:      resume.ResumeName,
		ResumeBasicInfo: resume.BasicInfo,
		ProjectIDs:      resume.ProjectIDs,
	}

	encoding.HandleSuccess(c, data)
}
