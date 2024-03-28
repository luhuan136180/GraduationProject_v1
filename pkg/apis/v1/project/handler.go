package project

import (
	"context"
	"encoding/json"
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

type projectHandlerOption struct {
	db *gorm.DB
}

type projectHandler struct {
	projectHandlerOption
}

func newProjectHandler(option projectHandlerOption) *projectHandler {
	return &projectHandler{
		projectHandlerOption: option,
	}
}

func (h *projectHandler) createProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := createReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role != model.RoleTypeTeacher {
		zap.L().Error("the operator's authority is illegal")
	}

	// 验证老师合法（x专业的老师只能创建x专业的项目）--搁置

	//
	projectBadicInfo, _ := json.Marshal(req.ProjectBasicInfo)

	_, user, err := dao.GetUserByID(ctx, h.db, strconv.FormatInt(request.GetUserIdFromCtx(ctx), 10))
	if err != nil {
		zap.L().Error(" dao.GetUserByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	project, err := dao.InsertProject(ctx, h.db, model.Project{
		ProjectName:      req.ProjectName,
		ProjectBasicInfo: projectBadicInfo,
		Title:            req.Title,
		CreatorUID:       user.UID,
		Creator:          user.Username,
	})
	if err != nil {
		zap.L().Error("dao.InsertProject", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	encoding.HandleSuccess(c, strconv.FormatInt(project.ID, 10))
}

func (h *projectHandler) deleteProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role != model.RoleTypeTeacher && Role != model.RoleTypeSuperAdmin && Role != model.RoleTypeCollegeAdmin {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	req := deleteProjectReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	found, _, err := dao.GetProjectByID(ctx, h.db, req.ID)
	if err != nil {
		zap.L().Error("the project not found", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if !found {
		zap.L().Error("the project not found", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	// 软删除
	err = dao.DeleteProjectByID(ctx, h.db, req.ID)
	if err != nil {
		zap.L().Error("close project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *projectHandler) projectList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := projectListReq{}
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

	projectList := make([]model.Project, 0)
	var count int64

	db := h.db.Model(&model.Project{})
	if request.GetRoleTypeFromCtx(ctx) == model.RoleTypeStudent {
		db = db.Where("status = ?", model.ProjectStatusPASS)
	}

	if req.ProjectName != "" {
		db = db.Where("project_name like ?", "%"+req.ProjectName+"%")
	}
	if req.Title != "" {
		db = db.Where("title like ?", "%"+req.Title+"%")
	}
	if req.Creator != "" {
		db = db.Where("creator like ?", "%"+req.Creator+"%")
	}
	if req.Auditor != "" {
		db = db.Where("auditor like ?", "%"+req.Auditor+"%")
	}
	if len(req.Status) > 0 && request.GetRoleTypeFromCtx(ctx) != model.RoleTypeStudent {
		db = db.Where("status in (?)", req.Status)
	}

	if request.GetRoleTypeFromCtx(ctx) != model.RoleTypeSuperAdmin {
		db = db.Where("profession_hash_id = ?", user.ProfessionHashID)
	}

	err = db.Count(&count).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = db.Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&projectList).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, projectListResp{Count: count, Projects: projectList})
}
