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
		professionHashID string `json:"profession_hash_id"`
	}

	deleteProjectReq struct {
		ID int64 `json:"id"`
	}
	projectListReq struct {
		model.ProjectOption
		Size int `json:"size"`
		Page int `json:"page"`
	}

	projectListResp struct {
		Count    int64              `json:"count"`
		Projects []projectBasicInfo `json:"projects"`
	}

	projectBasicInfo struct {
		ID           int64  `json:"id"` // 项目id
		ProjectName  string `json:"project_name"`
		ProjectTtile string `json:"title"`

		Status model.ProjectStatus `json:"status"` // 项目状态
	}

	getProjectReq struct {
		Size int `json:"size"`
		Page int `json:"page"`
	}

	projectDetailReq struct {
		ID int64 `json:"id"`
	}

	projectDetailResp struct {
		ID               int64               `json:"id"`
		ProjectName      string              `json:"projectName"`
		ProjectBasicInfo datatypes.JSON      `json:"projectBasicInfo"`
		ProjectFile      []byte              `json:"projectFile"`
		Title            string              `json:"title"`
		Status           model.ProjectStatus `json:"status"`
		ProfessionHashID string              `json:"professionHashID"`
		ProfessionName   string              `json:"professionName"`
		CollegeName      string              `json:"collegeName"`

		Creator               string `json:"creator"`
		Auditor               string `json:"auditor"`
		Participator          string `json:"participator"`
		ParticipatorClassName string `json:"participatorClassName"`
		ParticipatorClassID   int    `json:"participatorClassID"`
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
)
