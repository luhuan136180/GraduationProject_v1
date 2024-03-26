package license

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"v1/pkg/server/errutil"
	"v1/pkg/utils"
)

type licenseVerifyReq struct {
	License string `json:"license"`
}

type AuthorityLicenseResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	License   string `json:"license"`
	CreatedAt int64  `json:"created_at"`
	ExpireAt  int64  `json:"expire_at"`
	Valid     bool   `json:"valid"`
}

type response struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    AuthorityLicenseResp `json:"data,omitempty"`
}

// -X cspm/pkg/license.host=https://delivery-platform-sit.tensorsecurity.cn
var host = "https://delivery-platform-release.tensorsecurity.cn"

func GetAuthorityLicense(license string) (*AuthorityLicenseResp, error) {
	r := licenseVerifyReq{License: license}
	b, err := json.Marshal(r)
	if err != nil {
		zap.L().Error("marshal license info failed", zap.Error(err))
		return nil, errutil.ErrInternalServer
	}

	body, err := utils.RequestWithData(http.MethodPost, host+"/api/v1/cspm/license/verify", nil, bytes.NewReader(b))
	if err != nil {
		zap.L().Error("verify license failed", zap.Error(err))
		return nil, errutil.ErrInternalServer
	}

	resp := response{}
	if err = json.Unmarshal(body, &resp); err != nil {
		zap.L().Error("unmarshal license info failed", zap.Error(err))
		return nil, errutil.ErrInternalServer
	}

	if resp.Code != 0 {
		if resp.Code == 404 {
			return nil, errutil.ErrInvalidLicense
		}

		zap.L().Error("get license error", zap.Any("resp", resp))
		return nil, errutil.ErrInternalServer
	}

	return &resp.Data, nil
}
