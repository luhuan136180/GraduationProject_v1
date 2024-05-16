package contract

import "v1/pkg/model"

type (
	blockListReq struct {
		Page      int    `form:"page"`
		Size      int    `form:"size"`
		SaveType  string `form:"type"`
		BlodkHash string `form:"block_hash"`
	}

	blockListResp struct {
		Count int64                `json:"count"`
		Items []model.BlockSaveLog `json:"items"`
	}

	getContentReq struct {
		KeyHash string `form:"key_hash"`
	}

	getContentResp struct {
		Value []string `json:"value"`
	}
)
