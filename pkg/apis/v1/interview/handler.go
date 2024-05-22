package interview

import (
	"context"
	"errors"
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

	// get firm
	firm, err := dao.GetFirmByHashID(ctx, h.db, creator.FirmHashID)
	if err != nil {
		zap.L().Error(" dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

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
		FirmName:       firm.FirmName,
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
			Status:         req.Status,
		})
	} else if user.Role == model.RoleTypeStudent {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:          req.Title,
			IntervieweeUID: user.UID,
			Status:         req.Status,
		})
	} else if user.Role == model.RoleTypeFirm || user.Role == model.RoleTypeSuperAdmin {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:      req.Title,
			CreatorUID: user.UID,
			Status:     req.Status,
		})
	}
	if err != nil {
		zap.L().Error(" dao.FindInterviewByOption", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	data := []interviewListRespData{}

	for _, interview := range interviews {
		data = append(data, interviewListRespData{
			ID:             interview.ID,
			Ttile:          interview.Ttile,
			Interviewee:    interview.Interviewee,
			IntervieweeUID: interview.IntervieweeUID,
			Creator:        interview.Creator,
			CreatorUID:     interview.CreatorUID,
		})
	}

	encoding.HandleSuccess(c, interviewListResp{Total: count, Data: data})
}

func (h *interviewHandler) interviewListAccept(c *gin.Context) {
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
			Status:         req.Status,
		})
	} else if user.Role == model.RoleTypeStudent {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:          req.Title,
			IntervieweeUID: user.UID,
			Status:         req.Status,
		})
	} else if user.Role == model.RoleTypeFirm || user.Role == model.RoleTypeSuperAdmin {
		count, interviews, err = dao.FindInterviewByOption(ctx, h.db, model.InterviewOption{
			Size: req.Size,
			Page: req.Page,

			Title:      req.Title,
			CreatorUID: user.UID,
			Status:     req.Status,
		})
	}
	if err != nil {
		zap.L().Error(" dao.FindInterviewByOption", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	data := []interviewListRespOnlyAccpetData{}

	for _, interview := range interviews {
		basic := interview.Info.(map[string]interface{})

		data = append(data, interviewListRespOnlyAccpetData{
			ID:             interview.ID,
			Ttile:          interview.Ttile,
			Interviewee:    interview.Interviewee,
			IntervieweeUID: interview.IntervieweeUID,
			Creator:        interview.Creator,
			CreatorUID:     interview.CreatorUID,
			Date:           basic["date"].(string),
			FirmName:       interview.FirmName,
		})
	}

	encoding.HandleSuccess(c, interviewListOnlyAccpetResp{Total: count, Data: data})
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

	if req.Comment != "" {
		changeInfo := map[string]interface{}{
			"comment": req.Comment,
		}

		if err := h.db.WithContext(ctx).Model(&model.Interview{}).Where("id = ?", req.ID).Updates(changeInfo).Error; err != nil {
			zap.L().Error("dao.UpdateInterviewStatus", zap.Error(err))
			encoding.HandleError(c, errutil.ErrIllegalParameter)
			return
		}
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

	switch interview.Status {
	case "Proceed":
		interview.Status = "进行中"
	case "Accept":
		interview.Status = "未开始"
	case "End":
		interview.Status = "面试结束-通过"
	case "Failed":
		interview.Status = "面试结束-未通过"
	default:
		interview.Status = "拒绝"

	}

	encoding.HandleSuccess(c, interviewDetailResp{
		ID: interview.ID, Title: interview.Ttile, Info: interview.Info, Interviewee: interview.Interviewee, Status: interview.Status,
		Comment: interview.Comment,
		Flag:    interview.Flag, BlockHash: interview.BlockHash, ContractHashID: interview.ContractHashID, ContractKeyID: interview.ContractKeyID,
	})
}

func (h *interviewHandler) getMyRecruitList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getMyRecruitListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	if user.Role != model.RoleTypeFirm {
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	count, list, err := dao.GetRecruitList(ctx, h.db, user.UID, req.JobName, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetRecruitList", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到招聘信息"))
		return
	}

	encoding.HandleSuccess(c, getMyRecruitListResp{Count: count, Items: list})
}

func (h *interviewHandler) createRecruit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := createRecruitReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	if user.Role != model.RoleTypeFirm {
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	firm, err := dao.GetFirmByHashID(ctx, h.db, user.FirmHashID)
	if err != nil {
		zap.L().Error("dao.GetFirmByHashID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	data := model.Recruit{
		FirmHashID:   user.FirmHashID,
		FirmName:     firm.FirmName,
		JobName:      req.JobName,
		JobIntroduce: req.JobIntroduce,
		JobCondition: req.JobCondition,
		JobSalary:    req.JobSalary,

		CreatorUID: user.UID,
		Creator:    user.Username,
	}

	_, err = dao.InsertRecruit(ctx, h.db, data)
	if err != nil {
		zap.L().Error("dao.InsertRecruit", zap.Error(err))
		encoding.HandleError(c, errors.New("创建招聘信息失败"))
		return
	}

	encoding.HandleSuccess(c, "success")
}

func (h *interviewHandler) deleteRecruit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.L().Error("get recruit id failed", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
	}

	// 鉴权
	_, user, err := dao.GetUserByUID(ctx, h.db, request.GetUserUIDFromCtx(ctx))
	if err != nil {
		zap.L().Error("dao.GetUserByUID", zap.Error(err))
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	if user.Role != model.RoleTypeFirm {
		encoding.HandleError(c, errutil.ErrPermissionDenied)
		return
	}

	err = dao.DeleteRecruit(ctx, h.db, id)
	if err != nil {
		zap.L().Error("dao.DeleteRecruit", zap.Error(err))
		encoding.HandleError(c, errors.New("删除招聘信息失败"))
		return
	}

}

func (h *interviewHandler) getRecruitList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	req := getMyRecruitListReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		encoding.HandleError(c, errutil.ErrIllegalParameter)
		return
	}

	count, list, err := dao.GetRecruitList(ctx, h.db, "", req.JobName, req.Page, req.Size)
	if err != nil {
		zap.L().Error("dao.GetRecruitList", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到招聘信息"))
		return
	}

	encoding.HandleSuccess(c, getMyRecruitListResp{Count: count, Items: list})
}

func (s *interviewHandler) firmDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	// 验证权限
	Role := request.GetRoleTypeFromCtx(ctx)
	if Role == "" {
		encoding.HandleError(c, errutil.ErrInternalServer)
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

func (h *interviewHandler) getRecruitDeatial(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, v1.DefaultTimeout)
	defer cancel()

	id, _ := strconv.Atoi(c.Param("id"))

	recruit, err := dao.GetRecruitByID(ctx, h.db, id)
	if err != nil {
		zap.L().Error("dao.GetRecruitByID", zap.Error(err))
		encoding.HandleError(c, errors.New("未找到招聘信息"))
		return
	}

	encoding.HandleSuccess(c, recruit)
}
