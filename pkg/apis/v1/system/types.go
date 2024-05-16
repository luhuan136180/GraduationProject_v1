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
		FirmHashID       string `json:"firm_hash_id"`

		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	editUserReq struct {
		Id       string `json:"id"`
		Username string `json:"name"`

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
	getUserListResp struct {
		Total int64          `json:"total"`
		Items []userListItem `json:"items"`
	}

	userListItem struct {
		Id             string `json:"id"`
		Uid            string `json:"uid"`
		UserName       string `json:"username"`
		Name           string `json:"name"`
		Role           string `json:"role"`
		ProfessionName string `json:"profession_name"`
		ClassName      string `json:"class_name"`
	}

	getUserDetailResp struct {
		ID               int64  `json:"id"`
		UID              string `json:"uid"` // hash_id
		Username         string `json:"username"`
		Name             string `json:"name"` // 昵称
		Role             string `json:"role"`
		Password         string `json:"password"`
		ProfessionHashID string `json:"profession_hash_id"`
		ProfessionName   string `json:"profession_name"`
		ClassHashID      string `json:"class_hash_id"`
		ClassName        string `json:"class_name"`
		EmploymentStatus string `json:"employment_status"`

		Phone string `json:"phone"`
		Emial string `json:"email"`
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
		CollegeHashID string `form:"college_hash_id"`
	}

	getCollegeTreeResp struct {
		HashID      string `json:"hash_id"`
		CollegeName string `json:"college_name"`
	}

	collegeListReq struct {
		Page int `form:"page"`
		Size int `form:"size"`
	}

	collegeListResp struct {
		ID          int64  `json:"id"`
		HashID      string `json:"college_hash_id"`
		CollegeName string `json:"college_name"`
		CreatAt     int64  `json:"creat_at"`
	}

	getCollegeDetailReq struct {
		Uid string `form:"uid"`
	}

	getCollegeDetailResp struct {
		ID          int64  `json:"id"`
		HashID      string `json:"college_hash_id"`
		CollegeName string `json:"college_name"`
		CollegeInfo string `json:"college_info"`
		Creator     string `json:"creator"`
		CreatedAt   int64  `json:"created_at"`
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
		CollegeHashID  string `json:"college_hash_id"`
	}

	professionListReq struct {
		HashID string `form:"college_hash_id"`
	}

	professionListResp struct {
		Count int                      `json:"count"`
		Items []professionListRespItem `json:"items"`
	}
	professionListRespItem struct {
		HashID         string `json:"hash_id"`
		ProfessionName string `json:"profession_name"`
		CreatedAt      int64  `json:"created_at"`
	}

	professionDetailReq struct {
		HashID string `form:"hash_id"`
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

	classListReq struct {
		Page             int    `form:"page"`
		Size             int    `form:"size"`
		Name             string `form:"name"`
		ProfessionHashID string `form:"profession_hash_id"`
	}

	classListResp struct {
		Count int64               `json:"count"`
		Items []classListRespItem `json:"items"`
	}
	classListRespItem struct {
		ClassID     int    `json:"class_id"`
		ClassName   string `json:"class_name"`
		ClassHashID string `json:"class_hash_id"`
		CreatAt     int64  `json:"creat_at"`
	}

	firmListReq struct {
		Page int    `form:"page"`
		Size int    `form:"size"`
		Name string `form:"firm_name"`
	}

	firmListResp struct {
		Count int64         `json:"count"`
		Items []*model.Firm `json:"items"`
	}

	createFirmReq struct {
		FirmName string `form:"firm_name"`
		FirmInfo string `form:"firm_info"`
	}

	deleteFirmReq struct {
		HashID string `form:"hash_id"`
	}

	firmDetailReq struct {
		HashID string `form:"hash_id"`
	}

	firmTreeResp struct {
		Items []firmTreeRespItem `json:"items"`
	}

	firmTreeRespItem struct {
		FirmHashID string `json:"firm_hash_id"`
		FirmName   string `json:"firm_name"`
	}
)
