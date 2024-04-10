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
		Total int64                   `json:"total"`
		Data  []interviewListRespData `json:"data"`
	}

	interviewListRespData struct {
		ID             int64  `json:"id"`
		Ttile          string `json:"ttile"`
		Interviewee    string `json:"interviewee"`     // 面试者
		IntervieweeUID string `json:"interviewee_uid"` // uid
		Creator        string `json:"creator"`         // 面试创建者
		CreatorUID     string `json:"creator_uid"`
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

		Flag           bool   `json:"flag""` // 是否上链; false:没有;true:上链
		ContractHashID string `json:"contract_hash_id"`
		ContractKeyID  string `json:"contract_key_id"`
	}
)
