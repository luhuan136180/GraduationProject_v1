package project

import "v1/pkg/model"

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
		Count    int64           `json:"count"`
		Projects []model.Project `json:"projects"`
	}
)
