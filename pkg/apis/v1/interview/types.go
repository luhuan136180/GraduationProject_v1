package interview

import "v1/pkg/model"

type (
	createInterviewReq struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Date     int64  `json:"date"`
		Location string `json:"location"`
		Position string `json:"position"`

		IntervieweeUid string `json:"interviewee_uid"`
	}

	deleteInterviewReq struct {
		ID int64 `json:"id"`
	}

	interviewListReq struct {
		model.InterviewOption
	}

	interviewListResp struct {
		Total int64             `json:"total"`
		Data  []model.Interview `json:"data"`
	}

	interviewChangeStatusRep struct {
		ID     int64                 `json:"id"`
		Status model.InterviewStatus `json:"status"`
	}
	interviewChangeStatusResp struct {
		ID          int64                 `json:"id"`
		Title       string                `json:"title"`
		info        interface{}           `json:"info"`
		Interviewee string                `json:"interviewee"`
		Status      model.InterviewStatus `json:"status"`
	}
	interviewDetailResp struct {
		ID          int64                 `json:"id"`
		Title       string                `json:"title"`
		Info        interface{}           `json:"info"`
		Interviewee string                `json:"interviewee"`
		Status      model.InterviewStatus `json:"status"`
	}
)
