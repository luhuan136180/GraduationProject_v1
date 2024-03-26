package auth

import "v1/pkg/model"

type (
	loginReq struct {
		Account      string `json:"account"`       // 用户名
		Password     string `json:"password"`      // 密码
		CaptchaID    string `json:"captcha_id"`    // 验证码id
		CaptchaValue string `json:"captcha_value"` // 验证码
	}

	loginResp struct {
		ID       int64          `json:"id"`
		Username string         `json:"username"`
		Name     string         `json:"name"`
		Role     model.RoleType `json:"role"`
		Token    string         `json:"token"`
	}

	createCaptchaResp struct {
		CapthchaID string
		Image      string
	}
)
