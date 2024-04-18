package system

import "v1/pkg/model"

type (
	deleteUserReq struct {
		Id string `form:"id"`
	}

	createUserReq struct {
		Account  string `json:"account"`  // 用戶名
		Password string `json:"password"` // 密码
		Username string `json:"username"` // 昵称

		Role             string `json:"role"`
		ProfessionHashID string `json:"profession_hash_id"`
		ClassHashID      string `json:"class_hash_id"`

		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	editUserReq struct {
		Id       string `json:"id"`
		Username string `json:"username"`

		Role             string `json:"role"`
		ProfessionHashID string `json:"profession_hash_id"`
		ClassHashID      string `json:"class_hash_id"`

		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	getUserListReq struct {
		Page int `json:"page"`
		Size int `json:"size"`
		model.UserOption
	}

	userListItem struct {
		Id             string `json:"id"`
		UserName       string `json:"username"`
		Name           string `json:"name"`
		Role           string `json:"role"`
		ProfessionName string `json:"profession_name"`
		ClassName      string `json:"class_name"`
	}

	getUserDetailResp struct {
		ID               int64          `json:"id"`
		UID              string         `json:"uid"` // hash_id
		Username         string         `json:"username"`
		Name             string         `json:"name"` // 昵称
		Role             model.RoleType `json:"role"`
		Password         string         `json:"password"`
		ProfessionHashID string         `json:"profession_hash_id"`
		ProfessionName   string         `json:"profession_name"`
		ClassHashID      string         `json:"class_hash_id"`
		ClassName        string         `json:"class_name"`

		Phone string `json:"phone"`
		Emial string `json:"emial"`
		Head  string `json:"head"`
	}

	changePwdReq struct {
		OldPwd string `json:"old"`
		NewPwd string `json:"new"`
	}

	resetPwdReq struct {
		NewPwd string `json:"password"`
	}

	ceateCollegeReq struct {
		CollegeName string      `json:"college_name"`
		CollegeInfo collegeInfo `json:"college_info"`
	}

	deleteCollegeReq struct {
		CollegeHashID string `json:"college_hash_id"`
	}

	getCollegeTreeResp struct {
		HashID      string `json:"hash_id"`
		CollegeName string `json:"college_name"`
	}

	ceateProfessionReq struct {
		CollegeHashID  string         `json:"college_hash_id"`
		ProfessionName string         `json:"profession_name"`
		ProfessionInfo professionInfo `json:"profession_info"`
	}
	ceateProfessionResp struct {
		ProfessionHashID string `json:"profession_hash_id"`
	}

	collegeInfo struct {
		Info string `json:"info"`
	}
	professionInfo struct {
		Info string `json:"info"`
	}

	deleteProfessionrReq struct {
		HashID string `form:"hash_id"`
	}
	getProfessionTreeResp struct {
		HashID         string `json:"hash_id"`
		ProfessionName string `json:"profession_name"`
	}

	createClassReq struct {
		ProfessionHashID string `json:"profession_hash_id"`
		ClassName        string `json:"class_name"`
		ClassID          int    `json:"class_id"`
	}
	ceateClassResp struct {
		ClassHashID string `json:"class_hash_id"`
	}

	deleteClassReq struct {
		HashID string `form:"hash_id"`
	}

	getClassTreeResp struct {
		ClassHashID string `json:"class_hash_id"`
		ClassInfo   string `json:"class_name"`
	}
)
