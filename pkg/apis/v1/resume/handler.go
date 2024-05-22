package resume

import (
	"context"
	"errors"
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

	_, profession, err := dao.GetProfessionByHashID(ctx, h.db, user.ProfessionHashID)
	if err != nil {
		zap.L().Error("dao.GetProfessionByHashID", zap.Error(err))
		encoding.HandleError(c, errors.New("get profession failed"))
		return
	}

	req.ResumeInfo.PorfessionName = profession.ProfessionName
	req.ResumeInfo.CollegeName = profession.CollegeName
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

func (h *resumeHandler) projectTreeList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	found, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil || !found {
		zap.L().Error("get user info failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	projectList := make([]model.Project, 0)
	db := h.db.WithContext(ctx).Model(&model.Project{}).Where("participator_id = ?", user.UID)
	err = db.Find(&projectList).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	data := []projectBasicInfo{}
	for _, project := range projectList {
		data = append(data, projectBasicInfo{
			ID:          project.ID,
			ProjectName: project.ProjectName,
		})
	}

	encoding.HandleSuccess(c, data)
}

func (h *resumeHandler) deleteResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := deleteResumeReq{}
	err := c.ShouldBind(&req)
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
	if request.GetRoleTypeFromCtx(ctx) == model.RoleTypeStudent { // 学生只能看到自己的
		db = db.Where("creator = ?", user.Username)
	} else {
		if len(req.UserName) != 0 {
			db = db.Where("user_name = (?)", req.UserName)
		}
	}

	if req.Title != "" {
		db = db.Where("resume_name like ?", "%"+req.Title+"%")
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
	var data []resumeListRespData
	for _, resume := range resumeList {
		data = append(data, resumeListRespData{
			ID:           resume.ID,
			ResumeName:   resume.ResumeName,
			UserUID:      resume.UserUid,
			UserName:     resume.UserName,
			ContractFlag: resume.Flag,
		})
	}

	encoding.HandleSuccess(c, resumeListResp{Count: count, ResumeList: data})
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

		Flag:           resume.Flag,
		BlockHash:      resume.BlockHash,
		ContractHashID: resume.ContractHashID,
		ContractKeyID:  resume.ContractKeyID,
	}

	encoding.HandleSuccess(c, data)
}

func (h *resumeHandler) resumeListByUid(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	idStr := c.Param("uid")
	if idStr == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, idStr)
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errors.New("this user is not exit"))
		return
	}

	if user.Role != model.RoleTypeStudent {
		zap.L().Info("this user not have resumes")
		encoding.HandleSuccess(c)
		return
	}

	db := h.db.Model(&model.Resume{}).Where("user_uid = ?", user.UID)

	resumeList := make([]model.Resume, 0)
	var count int64

	err = db.Count(&count).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	err = db.Limit(5).Find(&resumeList).Error
	if err != nil {
		zap.L().Error("get project failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	var data []resumeListRespData
	for _, resume := range resumeList {
		data = append(data, resumeListRespData{
			ID:           resume.ID,
			ResumeName:   resume.ResumeName,
			UserUID:      resume.UserUid,
			UserName:     resume.UserName,
			ContractFlag: resume.Flag,
		})
	}
	encoding.HandleSuccess(c, resumeListResp{Count: count, ResumeList: data})
}

func (h *resumeHandler) getListOnBlock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errors.New("this user is not exit"))
		return
	}

	if user.Role != model.RoleTypeStudent {
		zap.L().Info("this user not have resumes")
		encoding.HandleSuccess(c)
		return
	}

	resumeList := make([]model.Resume, 0)
	err = h.db.Model(&model.Resume{}).Where("user_uid = ?", user.UID).Where("flag = ?", true).Find(&resumeList).Error
	if err != nil {
		zap.L().Error("GET RESUME failed:", zap.Error(err))
		encoding.HandleError(c, errors.New("获取简历失败"))
		return
	}

	var data []resumeListRespData
	for _, resume := range resumeList {
		data = append(data, resumeListRespData{
			ID:           resume.ID,
			ResumeName:   resume.ResumeName,
			UserUID:      resume.UserUid,
			UserName:     resume.UserName,
			ContractFlag: resume.Flag,
		})
	}

	encoding.HandleSuccess(c, resumeListResp{ResumeList: data})
}

func (h *resumeHandler) SendResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := SendResumeResp{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	data := model.RecruitResume{
		RecruitID: req.RecruitID,
		ResumeID:  req.ResumeID,
	}
	if err := h.db.WithContext(ctx).Create(&data).Error; err != nil {
		zap.L().Error("insert failed:", zap.Error(err))
		encoding.HandleError(c, errors.New("发送简历失败"))
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *resumeHandler) getResumeByRecruitID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getResumeByRecruitIDReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Size <= 0 {
		req.Size = 1
	}

	var resumelist []model.RecruitResume
	var count int64

	if err := h.db.WithContext(ctx).Model(&model.RecruitResume{}).Where("recruit_id = ?", req.RecruitID).Count(&count).Error; err != nil {
		zap.L().Error("get  failed:", zap.Error(err))
		encoding.HandleError(c, errors.New("获取简历失败"))
		return
	}

	if err := h.db.WithContext(ctx).Limit(req.Size).Offset((req.Page-1)*req.Size).Model(&model.RecruitResume{}).Where("recruit_id = ?", req.RecruitID).Find(&resumelist).Error; err != nil {
		zap.L().Error("get failed:", zap.Error(err))
		encoding.HandleError(c, errors.New("获取简历失败"))
		return
	}

	resumeIDs := make([]int64, 0)
	for _, value := range resumelist {
		resumeIDs = append(resumeIDs, value.ResumeID)
	}

	resumes := make([]model.Resume, 0)
	if err := h.db.Model(&model.Resume{}).Where("id in (?)", resumeIDs).Find(&resumes).Error; err != nil {
		zap.L().Error("get failed:", zap.Error(err))
		encoding.HandleError(c, errors.New("获取简历失败"))
		return
	}

	encoding.HandleSuccess(c, getResumeByRecruitIDResp{Count: count, Items: resumes})
}
