package resume

import "v1/pkg/model"

type (
	createReq struct {
		ResumeName       string  `json:"resume_name"`
		ProjectIDs       []int64 `json:"project_ids"`
		model.ResumeInfo `json:"resume_info"`
	}

	deleteResumeReq struct {
		ResumeID int64 `json:"resume_id"`
	}

	resumeListReq struct {
		Page               int `json:"page"`
		Size               int `json:"size"`
		model.ResumeOption `json:"resume_option"`
	}
	resumeListResp struct {
		Count      int64          `json:"count,omitempty"`
		resumeList []model.Resume `json:"resume_list"`
	}
)
