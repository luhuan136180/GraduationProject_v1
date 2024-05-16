package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"v1/pkg/utils"
)

type systemHandlerOption struct {
	db *gorm.DB
}

type systemHandler struct {
	systemHandlerOption
}

func newSystemHandler(option systemHandlerOption) *systemHandler {
	return &systemHandler{
		systemHandlerOption: option,
	}
}

func (s *systemHandler) userListTest(c *gin.Context) {
	type resp struct {
		Total int64        `json:"total"`
		Items []model.User `json:"items"`
	}

	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	total, users, err := dao.GETUserList(ctx, s.db, 1, 10, model.UserOption{})
	if err != nil {
		zap.L().Error("dao.GETUserList", zap.Error(err))
		encoding.HandleError(c, errors.New("get users list err"))
		return
	}
	encoding.HandleSuccess(c, resp{Total: total, Items: users})
}

// API
func (h *systemHandler) deleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(c)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role == model.RoleTypeStudent || Role == model.RoleTypeNormal {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
	}

	req := deleteUserReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 检测用户是否存在
	ok, user, err := dao.GetUserByID(ctx, h.db, req.Id)
	if err != nil {
		zap.L().Error("find user by id failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	if !ok {
		zap.L().Error("failed to delete,the user is not found")
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	// 删除用户
	if err = dao.DeleteUserByID(ctx, h.db, req.Id); err != nil {
		zap.L().Error("delete user failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	// 删除其他相关数据
	// 简历的数据库信息
	err = h.db.WithContext(ctx).Where("user_uid = ?", user.UID).Delete(&model.Resume{}).Error
	if err != nil {
		zap.L().Error("delete user related assets failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	// 项目
	err = h.db.WithContext(ctx).Where("creator = ?", user.Username).Delete(&model.Project{}).Error
	if err != nil {
		zap.L().Error("delete user related benchmarks failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	// 面试记录
	err = h.db.WithContext(ctx).Where("interviewee = ?", user.Username).Delete(&model.Interview{}).Error
	if err != nil {
		zap.L().Error("delete user related credentials failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	encoding.HandleSuccess(c)
}

func (s *systemHandler) createUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := createUserReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 检验密码合规性
	if !utils.CheckPWD(req.Password) {
		zap.L().Error("password illegal")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 验证学院和班级合理性
	if req.Role == string(model.RoleTypeStudent) {
		if _, _, err := dao.GetProfessionByHashID(ctx, s.db, req.ProfessionHashID); err != nil {
			zap.L().Error("not found real profession Info", zap.Error(err))
			encoding.HandleError(c, errutil.ErrIllegalParameter)
			return
		}

		if _, _, err := dao.GetClassByHashID(ctx, s.db, req.ClassHashID); err != nil {
			zap.L().Error("not found real class Info", zap.Error(err))
			encoding.HandleError(c, errutil.ErrIllegalParameter)
			return
		}
	} else {
		if _, _, err := dao.GetProfessionByHashID(ctx, s.db, req.ProfessionHashID); err != nil {
			admin, _ := dao.GetSuperProfession(ctx, s.db)
			req.ProfessionHashID = admin.HashID
		}

		if found, _, err := dao.GetClassByHashID(ctx, s.db, req.ClassHashID); err != nil || !found {
			req.ClassHashID = ""
		}
	}

	// 账号重复性验证
	ok, _, err := dao.GetUserByAccount(ctx, s.db, req.Account)
	if ok {
		zap.L().Error("this user account is already exists")
		encoding.HandleError(c, errutil.NewError(400, "account already exists"))
		return
	}
	if err != nil {
		zap.L().Error("dao.GetUserByAccount", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	user, err := dao.InsertUser(ctx, s.db, model.User{
		UID:              utils.NextID(),
		Username:         req.Account,
		Name:             req.Username,
		Role:             model.RoleType(req.Role),
		Password:         utils.MD5Hex(req.Password),
		ProfessionHashID: req.ProfessionHashID,
		ClassHashID:      req.ClassHashID,
		Creator:          request.GetUsernameFromCtx(ctx),
		Updater:          request.GetUsernameFromCtx(ctx),
		Status:           model.UserStatusNormal,
		Phone:            req.Phone,
		Emial:            req.Email,
	})
	if err != nil {
		zap.L().Error("dao.InsertUser", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	encoding.HandleSuccess(c, strconv.FormatInt(user.ID, 10))
}

func (h *systemHandler) getUserList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getUserListReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if len(req.Status) == 1 && req.Status[0] == "" {
		req.Status = nil
	}
	if len(req.ClassHashIDs) == 1 && req.ClassHashIDs[0] == "" {
		req.ClassHashIDs = nil
	}
	if len(req.ProfessionHashIDs) == 1 && req.ProfessionHashIDs[0] == "" {
		req.ProfessionHashIDs = nil
	}
	if len(req.RoleTypes) == 1 && string(req.RoleTypes[0]) == "" {
		req.RoleTypes = nil
	}

	num, userList, err := dao.GETUserList(ctx, h.db, req.Page, req.Size, req.UserOption)
	if err != nil {
		zap.L().Error("dao.GETUserList error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrSelectUserList)
		return
	}

	var pisd, cids []string
	for _, user := range userList {
		if user.ProfessionHashID != "" {
			pisd = append(pisd, user.ProfessionHashID)
		}
		if user.ClassHashID != "" {
			cids = append(cids, user.ClassHashID)
		}
	}
	// get profession
	professions, err := dao.GetProfessionsByHashIDs(ctx, h.db, pisd)
	if err != nil {
		zap.L().Error("dao.GetProfessionsByHashIDs error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrNotFound)
		return
	}

	professionMap := make(map[string]model.Profession) // hash_id : value
	for _, profession := range professions {
		if _, found := professionMap[profession.HashID]; !found {
			professionMap[profession.HashID] = profession
		}
	}

	// get class info
	classes, err := dao.GetClassByHashIDs(ctx, h.db, cids)
	if err != nil {
		zap.L().Error("dao.GetClassByHashIDs error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrNotFound)
		return
	}
	classMap := make(map[string]model.Class) // hash_id : value
	for _, class := range classes {
		if _, found := classMap[class.ClassHashID]; !found {
			classMap[class.ClassHashID] = class
		}
	}

	items := make([]userListItem, 0)
	for _, user := range userList {
		profession, _ := professionMap[user.ProfessionHashID]

		class, _ := classMap[user.ClassHashID]

		items = append(items, userListItem{
			UserName: user.Username,
			Name:     user.Name,
			Uid:      user.UID,
			Id:       strconv.FormatInt(user.ID, 10),
			Role:     string(user.Role),

			ClassName:      class.ClassName,
			ProfessionName: profession.ProfessionName,
		})
	}

	encoding.HandleSuccess(c, getUserListResp{Total: num, Items: items})
}

func (h *systemHandler) getUserDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	idStr := c.Param("id")
	if idStr == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	found, user, err := dao.GetUserByID(ctx, h.db, idStr)
	if err != nil || !found {
		zap.L().Error("dao.GetUserByID error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrEditUserInfo)
		return
	}

	var profession model.Profession
	var class model.Class
	if user.ProfessionHashID != "" {
		_, profession, err = dao.GetProfessionByHashID(ctx, h.db, user.ProfessionHashID)
		if err != nil {
			zap.L().Error("dao.GetProfessionByHashID error", zap.Error(err))
			encoding.HandleError(c, errutil.ErrEditUserInfo)
			return
		}
	} else {
		profession.ProfessionName = ""
		profession.HashID = ""
	}

	if profession.ProfessionName == "admin" && profession.CollegeName == "admin" {
		profession.ProfessionName = ""
		profession.HashID = ""
	}

	if user.ClassHashID == "" {
		class.ClassHashID = ""
		class.ClassName = ""
	} else {
		_, class, err = dao.GetClassByHashID(ctx, h.db, user.ClassHashID)
		if err != nil {
			zap.L().Error("dao.GetProfessionByHashID error", zap.Error(err))
			encoding.HandleError(c, errutil.ErrEditUserInfo)
			return
		}
	}

	data := getUserDetailResp{
		ID:               user.ID,
		UID:              user.UID,
		Username:         user.Username,
		Name:             user.Name,
		Role:             model.RoleMap[user.Role],
		ProfessionHashID: profession.HashID,
		ProfessionName:   profession.ProfessionName,
		ClassHashID:      class.ClassHashID,
		ClassName:        class.ClassName,
		EmploymentStatus: user.EmploymentStatus,

		Phone: user.Phone,
		Emial: user.Emial,
		Head:  user.Head,
	}

	encoding.HandleSuccess(c, data)
}

func (h *systemHandler) editUserInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := editUserReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 检查该用户是否存在
	ok, _, err := dao.GetUserByID(ctx, h.db, req.Id)
	if err != nil {
		zap.L().Error("find user by id failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	if !ok {
		zap.L().Error("failed to delete,the user is not found")
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	updater := request.GetUsernameFromCtx(ctx)

	err = dao.UpdateUserInfo(ctx, h.db, req.Id, model.User{
		Name: req.Username,

		ProfessionHashID: req.ProfessionHashID,
		ClassHashID:      req.ClassHashID,

		Phone:   req.Phone,
		Emial:   req.Email,
		Updater: updater,
	})
	if err != nil {
		zap.L().Error("dao.ChangeUserInfo error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrEditUserInfo)
		return
	}

	encoding.HandleSuccess(c)
}

func (h *systemHandler) changeUserPwd(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := changePwdReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	id := request.GetUserIdFromCtx(ctx)
	ok, u, err := dao.GetUserByID(ctx, h.db, strconv.FormatInt(id, 10))
	if !ok { // 不存在
		zap.L().Error("this user not  exists")
		encoding.HandleError(c, errutil.NewError(400, "username or password error"))
		return
	}
	if err != nil {
		zap.L().Error("dao.GetUserByAccount", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	if u.Password != utils.MD5Hex(req.OldPwd) {
		zap.L().Error("oldpassword error")
		encoding.HandleError(c, errutil.NewError(400, "username or password error"))
		return
	}

	if !utils.CheckPWD(req.NewPwd) {
		zap.L().Error("password illegal")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	updated := map[string]interface{}{
		"password":   utils.MD5Hex(req.NewPwd),
		"updated_at": time.Now().UnixMilli(),
		"updater":    request.GetUsernameFromCtx(ctx),
	}
	err = h.db.WithContext(ctx).Model(model.User{}).Where("id = ?", id).Updates(updated).Error
	if err != nil {
		zap.L().Error("update pwd failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrChangeUserPWD)
		return
	}

	encoding.HandleSuccess(c)
}

func (h *systemHandler) resetUserPWD(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role != model.RoleTypeSuperAdmin && Role != model.RoleTypeCollegeAdmin {
		zap.L().Error("the operator's authority is illegal")
	}

	req := resetPwdReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	id := c.Param("id")
	ok, _, err := dao.GetUserByID(ctx, h.db, id)
	if !ok {
		zap.L().Error("this user not  exists")
		encoding.HandleError(c, errutil.NewError(400, "username or password error"))
		return
	}
	if err != nil {
		zap.L().Error("dao.GetUserByAccount", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	// 检验密码合规性
	if !utils.CheckPWD(req.NewPwd) {
		zap.L().Error("password illegal")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	updated := map[string]interface{}{
		"password":   utils.MD5Hex(req.NewPwd),
		"updated_at": time.Now().UnixMilli(),
		"updater":    request.GetUsernameFromCtx(ctx),
	}
	ID, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("strconv.Atoi", zap.Error(err))
		encoding.HandleError(c, errutil.NewError(400, "reset password failed"))
		return
	}

	err = h.db.WithContext(ctx).Model(model.User{}).Where("id = ?", ID).Updates(updated).Error
	if err != nil {
		zap.L().Error("update pwd failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrChangeUserPWD)
		return
	}

	encoding.HandleSuccess(c)
}

func (s *systemHandler) createCollege(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := ceateCollegeReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 专业重复性验证
	found, _, err := dao.GetCollegeByHashID(ctx, s.db, utils.HashCollegeID(req.CollegeName))
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("not found real profession Info", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	if found {
		zap.L().Error("this profession account is already exists")
		encoding.HandleError(c, errutil.NewError(400, "college already exists"))
		return
	}

	CollegeInfoJson, err := json.Marshal(req.CollegeInfo)
	if err != nil {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	college, err := dao.InsertCollege(ctx, s.db, model.College{
		HashID:      utils.HashCollegeID(req.CollegeName),
		CollegeName: req.CollegeName,
		CollegeInfo: string(CollegeInfoJson),

		Creator: request.GetUsernameFromCtx(ctx),
		Updater: request.GetUsernameFromCtx(ctx),
	})
	if err != nil {
		zap.L().Error("dao.InsertUser", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	encoding.HandleSuccess(c, ceateProfessionResp{ProfessionHashID: college.HashID})
}

func (h *systemHandler) deleteCollege(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(c)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role == model.RoleTypeStudent || Role == model.RoleTypeNormal {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
	}

	req := deleteCollegeReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 检测
	ok, _, err := dao.GetCollegeByHashID(ctx, h.db, req.CollegeHashID)
	if err != nil {
		zap.L().Error("find profession by id failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	if !ok {
		zap.L().Error("failed to delete,the profession is not found")
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	// 删除
	if err = dao.DeleteCollege(ctx, h.db, req.CollegeHashID); err != nil {
		zap.L().Error("delete profession failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	encoding.HandleSuccess(c)
}

func (h *systemHandler) getCollegeTree(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	_, colleges, err := dao.GetColleges(ctx, h.db, 0, 0)
	if err != nil {
		zap.L().Error("dao.GetProfessions", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	var data []getCollegeTreeResp

	for _, college := range colleges {
		if college.CollegeName == "admin" {
			continue
		}
		data = append(data, getCollegeTreeResp{
			HashID:      college.HashID,
			CollegeName: college.CollegeName,
		})
	}

	encoding.HandleSuccess(c, data)
}

func (s *systemHandler) collegeList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		zap.L().Error("login user not hava enough permission")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := collegeListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if Role == model.RoleTypeStudent || Role == model.RoleTypeTeacher || Role == model.RoleTypeFirm {
		zap.L().Error("login user not hava enough permission")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	count, colleges, err := dao.GetColleges(ctx, s.db, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetProfessions", zap.Error(err))
		encoding.HandleError(c, errors.New(fmt.Sprintf("get college failed :%v", err.Error())))
		return
	}

	var resp []collegeListResp
	for _, college := range colleges {
		if college.CollegeName == "admin" {
			continue
		}
		data := collegeListResp{
			ID:          college.ID,
			HashID:      college.HashID,
			CollegeName: college.CollegeName,
			CreatAt:     college.CreatedAt,
		}
		resp = append(resp, data)
	}

	encoding.HandleSuccessList(c, count, resp)
}

func (s *systemHandler) getCollegeDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := getCollegeDetailReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, college, err := dao.GetCollegeByHashID(ctx, s.db, req.Uid)
	if err != nil {
		zap.L().Error("dao.GetCollegeByHashID", zap.Error(err))
		encoding.HandleError(c, errors.New("get college failed"))
		return
	}

	_, user, err := dao.GetUserByUID(ctx, s.db, college.Creator)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error(" dao.GetUserByUID", zap.Error(err))
			encoding.HandleError(c, errutil.ErrIllegalParameter)
			return
		}
		user = &model.User{}
	}
	collegeInfo := model.CollegeInfo{}
	_ = json.Unmarshal([]byte(college.CollegeInfo), &collegeInfo)

	resp := getCollegeDetailResp{
		ID:          college.ID,
		HashID:      college.HashID,
		CollegeName: college.CollegeName,
		CollegeInfo: collegeInfo.Info,
		Creator:     user.Username,
		CreatedAt:   college.CreatedAt,
	}
	encoding.HandleSuccess(c, resp)
}

func (s *systemHandler) createProfession(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := ceateProfessionReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 专业重复性验证
	found, _, err := dao.GetProfessionByHashID(ctx, s.db, utils.HashProfessionID(req.CollegeHashID, req.ProfessionName))
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("not found real profession Info", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	if found {
		zap.L().Error("this profession account is already exists")
		encoding.HandleError(c, errutil.NewError(400, "account already exists"))
		return
	}

	_, collegeItem, err := dao.GetCollegeByHashID(ctx, s.db, req.CollegeHashID)
	if err != nil {
		zap.L().Error("not found real college Info", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	professionInfoJson, err := json.Marshal(req.ProfessionInfo)
	if err != nil {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	profession, err := dao.InsertProfession(ctx, s.db, model.Profession{
		HashID:         utils.HashProfessionID(req.CollegeHashID, req.ProfessionName),
		CollegeHashID:  req.CollegeHashID,
		CollegeName:    collegeItem.CollegeName,
		ProfessionName: req.ProfessionName,
		ProfessionInfo: string(professionInfoJson),

		Creator: request.GetUsernameFromCtx(ctx),
		Updater: request.GetUsernameFromCtx(ctx),
	})
	if err != nil {
		zap.L().Error("dao.InsertUser", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	encoding.HandleSuccess(c, ceateProfessionResp{ProfessionHashID: profession.HashID})
}

func (h *systemHandler) deleteProfession(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(c)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role == model.RoleTypeStudent || Role == model.RoleTypeNormal {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := deleteProfessionrReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 删除班级
	if err := dao.DeleteClassByProfession(ctx, h.db, req.HashID); err != nil {
		zap.L().Error("delete class failed", zap.Error(err))
		encoding.HandleError(c, errors.New("删除班级失败"))
		return
	}

	// 删除专业
	if err = dao.DeleteProfession(ctx, h.db, req.HashID); err != nil {
		zap.L().Error("delete profession failed", zap.Error(err))
		encoding.HandleError(c, errors.New("删除专业失败"))
		return
	}

	encoding.HandleSuccess(c)
}

func (h *systemHandler) getProfessionTree(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	professions, err := dao.GetProfessions(ctx, h.db)
	if err != nil {
		zap.L().Error("dao.GetProfessions", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	var data []getProfessionTreeResp

	for _, profession := range professions {
		if profession.CollegeName == "admin" {
			continue
		}
		data = append(data, getProfessionTreeResp{
			HashID:         profession.HashID,
			ProfessionName: profession.ProfessionName,
			CollegeHashID:  profession.CollegeHashID,
		})
	}

	encoding.HandleSuccess(c, data)
}

func (s *systemHandler) professionList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := professionListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	professions, err := dao.GetProfessionByCollegeashID(ctx, s.db, req.HashID)
	if err != nil {
		zap.L().Error("dao.GetProfessionByCollegeashID", zap.Error(err))
		encoding.HandleError(c, errors.New("该学院为创建专业"))
		return
	}

	respdata := []professionListRespItem{}
	for _, profession := range professions {
		respdata = append(respdata, professionListRespItem{
			HashID:         profession.HashID,
			ProfessionName: profession.ProfessionName,
			CreatedAt:      profession.CreatedAt,
		})
	}

	encoding.HandleSuccess(c, professionListResp{Count: len(professions), Items: respdata})
}

func (s *systemHandler) professionDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := professionDetailReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, profession, err := dao.GetProfessionByHashID(ctx, s.db, req.HashID)
	if err != nil {
		zap.L().Error("dao.GetProfessionByHashID", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到该专业的详情信息"))
		return
	}

	_, user, err := dao.GetUserByUID(ctx, s.db, profession.Creator)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			zap.L().Error("dao.GetUserByUID", zap.Error(err))
			encoding.HandleError(c, errors.New("未找到用户"))
			return
		}
		user = &model.User{}
	}

	ProfessionInfo := model.ProfessionInfo{}
	_ = json.Unmarshal([]byte(profession.ProfessionInfo), &ProfessionInfo)

	profession.ProfessionInfo = ProfessionInfo.Info
	profession.Creator = user.Username

	encoding.HandleSuccess(c, profession)
}

func (s *systemHandler) createClass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := createClassReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 重复性验证
	found, _, err := dao.GetClassByHashID(ctx, s.db, utils.HashClassID(req.ProfessionHashID, req.ClassName, req.ClassID))
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("not found real class Info", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}
	if found {
		zap.L().Error("this class  is already exists")
		encoding.HandleError(c, errutil.NewError(400, "class already exists"))
		return
	}

	class, err := dao.InsertClass(ctx, s.db, model.Class{
		ProfessionHashID: req.ProfessionHashID,
		ClassHashID:      utils.HashClassID(req.ProfessionHashID, req.ClassName, req.ClassID),
		ClassName:        req.ClassName,
		ClassID:          req.ClassID,

		Creator: request.GetUsernameFromCtx(ctx),
		Updater: request.GetUsernameFromCtx(ctx),
	})
	if err != nil {
		zap.L().Error("dao.InsertUser", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	encoding.HandleSuccess(c, ceateClassResp{ClassHashID: class.ClassHashID})
}

func (h *systemHandler) deleteClass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(c)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}
	if Role == model.RoleTypeStudent || Role == model.RoleTypeNormal {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
	}

	req := deleteClassReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	ok, _, err := dao.GetClassByHashID(ctx, h.db, req.HashID)
	if err != nil {
		zap.L().Error("find class by id failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}
	if !ok {
		zap.L().Error("failed to delete,the class is not found")
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	// 删除班级
	if err = dao.DeleteClass(ctx, h.db, req.HashID); err != nil {
		zap.L().Error("delete class failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrDeleteUser)
		return
	}

	encoding.HandleSuccess(c)
}

func (h *systemHandler) getClassTree(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	PHashID := c.Param("profession_hash_id")
	if PHashID == "" {
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	classes, err := dao.GetClassesByPID(ctx, h.db, PHashID)
	if err != nil {
		zap.L().Error("dao.GetClassesByPID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	var data []getClassTreeResp

	for _, class := range classes {
		if class.ClassName == "admin" {
			continue
		}
		data = append(data, getClassTreeResp{
			ClassHashID: class.ClassHashID,
			ClassInfo:   class.ClassName + strconv.Itoa(class.ClassID),
		})
	}

	encoding.HandleSuccess(c, data)
}

func (h *systemHandler) uploaduserTouxiang(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("c.FormFile ERR:", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	dst := fmt.Sprintf("./file/touxiang/%s", file.Filename)
	dstUrl := fmt.Sprintf("http://localhost:9090/api/v1/auth/fs/%s", file.Filename)
	// 上传
	c.SaveUploadedFile(file, dst)

	err = dao.InsterTouxiang(ctx, h.db, request.GetUserUIDFromCtx(ctx), dstUrl)
	if err != nil {
		zap.L().Error("dao.InsterTouxiang", zap.Error(err))
		encoding.HandleError(c, errors.New("save failed"))
		return
	}
	encoding.HandleSuccess(c, dstUrl)
}

func (s *systemHandler) classList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := classListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	count, classes, err := dao.GetClasses(ctx, s.db, req.Name, req.ProfessionHashID, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetClasses", zap.Error(err))
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	items := make([]classListRespItem, 0)
	for _, class := range classes {
		items = append(items, classListRespItem{
			ClassID:     class.ClassID,
			ClassName:   class.ClassName,
			ClassHashID: class.ClassHashID,
			CreatAt:     class.CreatedAt,
		})
	}

	encoding.HandleSuccess(c, classListResp{Count: count, Items: items})
}

func (h *systemHandler) firmList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		zap.L().Error("login user not hava enough permission")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := firmListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	count, firms, err := dao.GetFirmList(ctx, h.db, req.Name, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetFirmList", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到对应的企业"))
		return
	}

	encoding.HandleSuccess(c, firmListResp{Count: count, Items: firms})
}

func (h *systemHandler) createFirm(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := createFirmReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到登录的用户数据"))
		return
	}

	_, err = dao.InsertFirm(ctx, h.db, model.Firm{
		FirmHashID: utils.HashFirmID(req.FirmName, req.FirmInfo),
		FirmName:   req.FirmName,
		FirmInfo:   req.FirmInfo,
		Creator:    user.Username,
		CreatorUID: user.UID,
	})
	if err != nil {
		zap.L().Error("dao.InsertFirm", zap.Error(err))
		encoding.HandleError(c, errors.New("创建企业失败"))
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *systemHandler) deleteFirm(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	// defer cancel()
	//
	// req := deleteFirmReq{}
	// err := c.ShouldBind(&req)
	// if err != nil {
	// 	zap.L().Error("c.ShouldBindQuery", zap.Error(err))
	// 	encoding.HandleError(c, errutil.ErrIllegalParameter)
	// 	return
	// }

}

func (s *systemHandler) createFirmUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := createUserReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 检验密码合规性
	if !utils.CheckPWD(req.Password) {
		zap.L().Error("password illegal")
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	// 验证企业的合理性
	firm, err := dao.GetFirmByHashID(ctx, s.db, req.FirmHashID)
	if err != nil {
		zap.L().Error("dao.GetFirmByHashID", zap.Error(err))
		encoding.HandleError(c, errors.New("该企业未找到"))
		return
	}

	// 账号重复性验证
	ok, _, err := dao.GetUserByAccount(ctx, s.db, req.Account)
	if ok {
		zap.L().Error("this user account is already exists")
		encoding.HandleError(c, errutil.NewError(400, "account already exists"))
		return
	}
	if err != nil {
		zap.L().Error("dao.GetUserByAccount", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	user, err := dao.InsertUser(ctx, s.db, model.User{
		UID:        utils.NextID(),
		Username:   req.Account,
		Name:       req.Username,
		Role:       model.RoleType(req.Role),
		Password:   utils.MD5Hex(req.Password),
		FirmHashID: firm.FirmHashID,
		Creator:    request.GetUsernameFromCtx(ctx),
		Updater:    request.GetUsernameFromCtx(ctx),
		Status:     model.UserStatusNormal,
		Phone:      req.Phone,
		Emial:      req.Email,
	})
	if err != nil {
		zap.L().Error("dao.InsertUser", zap.Error(err))
		encoding.HandleError(c, errutil.ErrCreateUser)
		return
	}

	encoding.HandleSuccess(c, strconv.FormatInt(user.ID, 10))
}

func (h *systemHandler) FirmUserList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getUserListReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	if req.Size > 10 || req.Size <= 0 {
		req.Size = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if len(req.Status) == 1 && req.Status[0] == "" {
		req.Status = nil
	}
	if len(req.ClassHashIDs) == 1 && req.ClassHashIDs[0] == "" {
		req.ClassHashIDs = nil
	}
	if len(req.ProfessionHashIDs) == 1 && req.ProfessionHashIDs[0] == "" {
		req.ProfessionHashIDs = nil
	}
	if len(req.RoleTypes) == 1 && string(req.RoleTypes[0]) == "" {
		req.RoleTypes = nil
	}

	num, userList, err := dao.GETUserList(ctx, h.db, req.Page, req.Size, req.UserOption)
	if err != nil {
		zap.L().Error("dao.GETUserList error", zap.Error(err))
		encoding.HandleError(c, errutil.ErrSelectUserList)
		return
	}

	items := make([]userListItem, 0)
	for _, user := range userList {

		items = append(items, userListItem{
			UserName: user.Username,
			Name:     user.Name,
			Uid:      user.UID,
			Id:       strconv.FormatInt(user.ID, 10),
			Role:     string(user.Role),
		})
	}

	encoding.HandleSuccess(c, getUserListResp{Total: num, Items: items})
}

func (s *systemHandler) firmDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	req := firmDetailReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("c.ShouldBindQuery", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	firm, err := dao.GetFirmByHashID(ctx, s.db, req.HashID)
	if err != nil {
		zap.L().Error("dao.GetFirmByHashID", zap.Error(err))
		encoding.HandleError(c, errors.New("该企业未找到"))
		return
	}

	encoding.HandleSuccess(c, firm)
}

func (h *systemHandler) firmTree(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
		return
	}

	if Role == model.RoleTypeNormal || Role == model.RoleTypeStudent {
		zap.L().Error("the operator's authority is illegal")
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	firms, err := dao.GetFirmTree(ctx, h.db)
	if err != nil {
		zap.L().Error("dao.GetFirmTree", zap.Error(err))
		encoding.HandleError(c, errors.New("未获取到企业列表"))
		return
	}

	items := make([]firmTreeRespItem, 0)
	for _, firm := range firms {
		items = append(items, firmTreeRespItem{
			FirmHashID: firm.FirmHashID,
			FirmName:   firm.FirmName,
		})
	}

	encoding.HandleSuccess(c, firmTreeResp{Items: items})
}
