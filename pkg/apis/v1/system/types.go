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
		Account        string `json:"account"`
		Role           string `json:"role"`
		Username       string `json:"username"`
		ProfessionName string `json:"profession_name"`
		ClassName      string `json:"class_name"`
		Phone          string `json:"phone"`
		Email          string `json:"email"`
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

	createClassReq struct {
		ProfessionHashID string `json:"profession_hash_id"`
		ClassName        string `json:"class_name"`
		ClassID          int    `json:"class_id"`
	}
	ceateClassResp struct {
		classHashID string `json:"profession_hash_id"`
	}

	deleteClassReq struct {
		HashID string `form:"hash_id"`
	}
)
