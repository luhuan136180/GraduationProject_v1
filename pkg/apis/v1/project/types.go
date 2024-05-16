package project

import (
	"gorm.io/datatypes"
	"v1/pkg/model"
)

type (
	createReq struct {
		ProjectName string `json:"project_name"`
		model.ProjectBasicInfo
		Title            string `json:"title"`
		ProfessionHashID string `json:"profession_hash_id"`
	}

	deleteProjectReq struct {
		ID int64 `form:"id"`
	}
	projectListReq struct {
		model.ProjectOption
		Size int `json:"size"`
		Page int `json:"page"`
	}

	projectListResp struct {
		Count    int64              `json:"total"`
		Projects []projectBasicInfo `json:"projects"`
	}

	projectBasicInfo struct {
		ID           int64  `json:"id"` // 项目id
		ProjectName  string `json:"project_name"`
		ProjectTtile string `json:"title"`

		Status       string `json:"status"` // 项目状态
		ContractFlag bool   `json:"contract_flag"`
	}

	getProjectReq struct {
		Size int     `json:"size"`
		Page int     `json:"page"`
		IDs  []int64 `json:"ids"`
	}

	projectDetailReq struct {
		ID int64 `json:"id"`
	}

	projectDetailResp struct {
		ID               int64          `json:"id"`
		ProjectName      string         `json:"projectName"`
		ProjectBasicInfo datatypes.JSON `json:"projectBasicInfo"`
		ProjectFile      []int          `json:"projectFile"`
		Title            string         `json:"title"`
		Status           string         `json:"status"`
		ProfessionHashID string         `json:"professionHashID"`
		ProfessionName   string         `json:"professionName"`
		CollegeName      string         `json:"collegeName"`

		Creator               string `json:"creator"`
		Auditor               string `json:"auditor"`
		Participator          string `json:"participator"`
		ParticipatorClassName string `json:"participatorClassName"`
		ParticipatorClassID   int    `json:"participatorClassID"`

		Flag           bool   `json:"flag""` // 是否上链; false:没有;true:上链
		ContractHashID string `json:"contract_hash_id"`
		ContractKeyID  string `json:"contract_key_id"`
	}

	chooseProjectReq struct {
		ProjectID int64 `json:"project_id"`
	}

	auditProjectReq struct {
		ProjectID int64 `json:"project_id"`
	}

	changeStatusReq struct {
		Status    model.ProjectStatus `json:"status"`
		ProjectID int64               `json:"project_id"`
	}

	fileListReq struct {
		Ids []int `json:"ids"`
	}

	fileListResp struct {
		Count int          `json:"count"`
		Items []model.File `json:"items"`
	}
)
