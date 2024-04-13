package project

import (
	"context"
	"encoding/json"
	"fmt"
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
	if Role == model.RoleTypeStudent {
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
		ProfessionHashID: user.ProfessionHashID,
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

	if len(req.Professions) == 1 && (req.Professions[0] == "string" || req.Professions[0] == "") {
		req.Professions = nil
	}

	user, err := dao.GetUserByUsername(ctx, h.db, request.GetUsernameFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUsername", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	projectList := make([]model.Project, 0)
	var count int64

	// 查询通过审核的
	db := h.db.Model(&model.Project{}).Where("status != ?", 1)

	if request.GetRoleTypeFromCtx(ctx) == model.RoleTypeStudent {
		db = db.Where("status = ?", model.ProjectStatusPASS)
	}

	if req.ProjectName != "" {
		db = db.Where("project_name like ?", "%"+req.ProjectName+"%")
	}
	if req.Title != "" {
		db = db.Where("title like ?", "%"+req.Title+"%")
	}
	if len(req.Professions) != 0 {
		db = db.Where("profession_hash_id in (?)", req.Professions)
	}

	// 非超管只能看自己学院
	if user.Role != model.RoleTypeSuperAdmin && user.Role != model.RoleTypeCollegeAdmin {
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

	data := []projectBasicInfo{}
	for _, project := range projectList {
		data = append(data, projectBasicInfo{
			ID:           project.ID,
			ProjectName:  project.ProjectName,
			ProjectTtile: project.Title,
			Status:       project.Status,
		})
	}

	encoding.HandleSuccess(c, projectListResp{Count: count, Projects: data})
}

func (h *projectHandler) getProjects(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	found, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil || !found {
		zap.L().Error("get user info failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	req := getProjectReq{}
	err = c.ShouldBindJSON(&req)
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

	db := h.db.WithContext(ctx).Model(&model.Project{})

	switch request.GetRoleTypeFromCtx(ctx) {
	case model.RoleTypeTeacher:
		db = db.Where("creator_uid = ?", user.UID)
	case model.RoleTypeCollegeAdmin:
		db = db.Where("profession_hash_id = ?", user.ProfessionHashID).Where("status = ?", model.ProjectStatusAudit)
	case model.RoleTypeStudent:
		db = db.Where("participator_id = ?", user.UID)
	}

	if request.GetRoleTypeFromCtx(ctx) != model.RoleTypeSuperAdmin {
		db = db.Where("profession_hash_id = ?", user.ProfessionHashID)
	}

	projectList := make([]model.Project, 0)
	var count int64

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

	data := []projectBasicInfo{}
	for _, project := range projectList {
		data = append(data, projectBasicInfo{
			ID:           project.ID,
			ProjectName:  project.ProjectName,
			ProjectTtile: project.Title,
			Status:       project.Status,
		})
	}

	encoding.HandleSuccess(c, projectListResp{Count: count, Projects: data})
}

func (h *projectHandler) projectDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := projectDetailReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	found, project, err := dao.GetProjectByID(ctx, h.db, req.ID)
	if err != nil || !found {
		zap.L().Error("dao.GetProjectByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// get external info
	found, profession, err := dao.GetProfessionByHashID(ctx, h.db, project.ProfessionHashID)
	if err != nil || !found {
		zap.L().Error("dao.GetProfessionByHashID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, project.ParticipatorID)
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	found, class, err := dao.GetClassByHashID(ctx, h.db, user.ClassHashID)
	if err != nil {
		zap.L().Error("dao.GetClassByHashID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	result := projectDetailResp{
		ID:               project.ID,
		ProjectName:      project.ProjectName,
		ProjectBasicInfo: project.ProjectBasicInfo,
		Title:            project.Title,
		Status:           project.Status,
		ProfessionHashID: project.ProfessionHashID,

		Creator:               project.Creator,
		Auditor:               project.Auditor,
		Participator:          project.Participator,
		ProfessionName:        profession.ProfessionName,
		CollegeName:           profession.CollegeName,
		ParticipatorClassName: class.ClassName,
		ParticipatorClassID:   class.ClassID,
	}

	encoding.HandleSuccess(c, result)
}

func (h *projectHandler) chooseProject(c *gin.Context) {
	// 先抢先得
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := chooseProjectReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if request.GetRoleTypeFromCtx(ctx) != model.RoleTypeStudent && request.GetRoleTypeFromCtx(ctx) != model.RoleTypeSuperAdmin {
		zap.L().Error("not student")
		encoding.HandleError(c, errutil.ErrIllegalOperation)
	}

	found, project, err := dao.GetProjectByID(ctx, h.db, req.ProjectID)
	if err != nil || !found {
		zap.L().Error("dao.GetProjectByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// == 5 表示未被选择
	if project.Status != 5 {
		zap.L().Error("project status is not : PASS")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	// 已经被选择
	if project.Participator != "" {
		zap.L().Error("this project is choose by other")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))

	err = dao.UpdateProjectParticipator(ctx, h.db, project.ID, *user)
	if err != nil {
		zap.L().Error("dao.UpdateProjectParticipator", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *projectHandler) auditProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := auditProjectReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if request.GetRoleTypeFromCtx(ctx) != model.RoleTypeCollegeAdmin && request.GetRoleTypeFromCtx(ctx) != model.RoleTypeSuperAdmin {
		zap.L().Error("not college admin")
		encoding.HandleError(c, errutil.ErrIllegalOperation)
		return
	}

	found, project, err := dao.GetProjectByID(ctx, h.db, req.ProjectID)
	if err != nil || !found {
		zap.L().Error("dao.GetProjectByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if project.Status != model.ProjectStatusAudit {
		zap.L().Error("the project status is not auditing")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))

	err = h.db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", req.ProjectID).Updates(map[string]interface{}{"audit_uid": user.UID, "auditor": user.Username, "status": 5}).Error
	if err != nil {
		zap.L().Error("change project status failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	encoding.HandleSuccess(c, "success")
}

func (h *projectHandler) changeStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := changeStatusReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, project, err := dao.GetProjectByID(ctx, h.db, req.ProjectID)
	if err != nil {
		zap.L().Error("dao.GetProjectByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if request.GetRoleTypeFromCtx(ctx) == model.RoleTypeStudent {
		zap.L().Error("illegal parameter, user is student")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = dao.UpdateProjectStatus(ctx, h.db, project.ID, req.Status)
	if err != nil {
		zap.L().Error("dao.UpdateProjectStatus", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	encoding.HandleSuccess(c, "success")
}

func (h *projectHandler) uploadFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Error("strconv.Atoi", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, project, err := dao.GetProjectByID(ctx, h.db, int64(id))
	if err != nil {
		zap.L().Error("dao.GetProjectByID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("c.FormFile ERR:", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	dst := fmt.Sprintf("./file/%s", file.Filename)

	// 上传
	c.SaveUploadedFile(file, dst)

	fileInfo, err := dao.InsterFile(ctx, h.db, model.File{FileName: file.Filename, FilePath: dst})
	if err != nil {
		zap.L().Error("dao.InsterFile", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = dao.UpdateProjectFiles(ctx, h.db, project.ID, fileInfo.ID)
	if err != nil {
		zap.L().Error("dao.UpdateProjectFiles", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	//
	encoding.HandleSuccess(c, "upload success")
}
