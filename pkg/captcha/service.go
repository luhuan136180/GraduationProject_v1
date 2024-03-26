package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
	"time"
)

type StringConfig struct {
	config *base64Captcha.DriverString
	store  base64Captcha.Store
}

var (
	GCLimitNumber = 10240
	Expiration    = 1 * time.Minute
	Config        = StringConfig{config: stringConfig(), store: base64Captcha.NewMemoryStore(GCLimitNumber, Expiration)}
)

func stringConfig() *base64Captcha.DriverString {
	stringType := &base64Captcha.DriverString{
		Height: 100,
		Width:  200,
		Length: 4,
		Source: "abcdefghijklmnopqrstuvwxyz",
		// Source: "QWERTYUIOPASDFGHJKLZXCVBNM",
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return stringType
}

func (config StringConfig) CreateCaptcha() (string, string, string, error) {
	c := base64Captcha.NewCaptcha(config.config, config.store)
	id, b64s, answer, err := c.Generate()
	return id, b64s, answer, err
}

func (config StringConfig) VerifyCaptcha(id, VerifyValue string) bool {
	// result 为步骤1 创建的图片验证码存储对象
	return config.store.Verify(id, VerifyValue, true)
}

func GetService() *StringConfig {
	return &Config
}
