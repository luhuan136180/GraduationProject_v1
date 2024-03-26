package i18n

import (
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

const (
	LangZH = "zh"
	LangEN = "en"
)

type locales map[language.Tag]catalog.Message

func init() {
	translates := make(map[string]locales)
	translates["internal server error"] = locales{
		language.Chinese: catalog.String("服务器内部错误"),
		language.English: catalog.String("internal server error"),
	}
	translates["illegal parameter"] = locales{
		language.Chinese: catalog.String("非法参数"),
		language.English: catalog.String("illegal parameter"),
	}
	translates["not found"] = locales{
		language.Chinese: catalog.String("数据未找到"),
		language.English: catalog.String("not found"),
	}
	translates["unauthorized"] = locales{
		language.Chinese: catalog.String("会话过期，请重新登录"),
		language.English: catalog.String("unauthorized"),
	}
	translates["permission denied"] = locales{
		language.Chinese: catalog.String("权限不足"),
		language.English: catalog.String("permission denied"),
	}
	translates["illegal operation"] = locales{
		language.Chinese: catalog.String("非法操作"),
		language.English: catalog.String("illegal operation"),
	}
	translates["delete user failed"] = locales{
		language.Chinese: catalog.String("删除用户失败"),
		language.English: catalog.String("delete user failed"),
	}
	translates["create user failed"] = locales{
		language.Chinese: catalog.String("创建用户失败"),
		language.English: catalog.String("create user failed"),
	}
	translates["change user Info failed"] = locales{
		language.Chinese: catalog.String("修改用户信息失败"),
		language.English: catalog.String("change user Info failed"),
	}
	translates["select userList failed"] = locales{
		language.Chinese: catalog.String("查询用户列表失败"),
		language.English: catalog.String("select userList failed"),
	}
	translates["reset user password failed"] = locales{
		language.Chinese: catalog.String("重置用户密码失败"),
		language.English: catalog.String("reset user password failed"),
	}
	translates["get audit logs failed"] = locales{
		language.Chinese: catalog.String("获取审计日志失败"),
		language.English: catalog.String("get audit logs failed"),
	}
	translates["edit credential info failed"] = locales{
		language.Chinese: catalog.String("修改云帐号信息失败"),
		language.English: catalog.String("edit credential info failed"),
	}
	translates["captcha value is wrong"] = locales{
		language.Chinese: catalog.String("验证码错误"),
		language.English: catalog.String("captcha value is wrong"),
	}
	translates["scan benchmark not found"] = locales{
		language.Chinese: catalog.String("扫描基线未找到"),
		language.English: catalog.String("scan benchmark not found"),
	}
	translates["scan task is running, please stop it and then try again"] = locales{
		language.Chinese: catalog.String("已有扫描任务正在进行中，请结束后再发起扫描"),
		language.English: catalog.String("scan task is running, please stop it and then try again"),
	}
	translates["no cloud account found, please add cloud account and then scan"] = locales{
		language.Chinese: catalog.String("没有找到云帐号，请添加云帐号后再扫描"),
		language.English: catalog.String("no cloud account found, please add cloud account and then scan"),
	}
	translates["scan task not found"] = locales{
		language.Chinese: catalog.String("扫描任务不存在"),
		language.English: catalog.String("scan task not found"),
	}
	translates["scan task not finished"] = locales{
		language.Chinese: catalog.String("扫描任务未完成"),
		language.English: catalog.String("scan task not finished"),
	}
	translates["risk data not found"] = locales{
		language.Chinese: catalog.String("风险数据未找到"),
		language.English: catalog.String("risk data not found"),
	}
	translates["scan assets not found"] = locales{
		language.Chinese: catalog.String("资产数据未找到"),
		language.English: catalog.String("scan assets not found"),
	}
	translates["failed to export risks info"] = locales{
		language.Chinese: catalog.String("导出风险信息失败"),
		language.English: catalog.String("failed to export risks info"),
	}
	translates["unknown compliance key"] = locales{
		language.Chinese: catalog.String("未知的合规类型"),
		language.English: catalog.String("unknown compliance type"),
	}
	translates["account name exist"] = locales{
		language.Chinese: catalog.String("账户名称重复"),
		language.English: catalog.String("account name exist"),
	}
	translates["Invalid access key"] = locales{
		language.Chinese: catalog.String("无效的访问密钥"),
		language.English: catalog.String("Invalid access key"),
	}
	translates["failed to get credential info"] = locales{
		language.Chinese: catalog.String("获取云账户信息失败"),
		language.English: catalog.String("failed to get credential info"),
	}
	translates["failed to export credential info"] = locales{
		language.Chinese: catalog.String("导出云账户详情失败"),
		language.English: catalog.String("failed to export credential info"),
	}
	translates["account already exists"] = locales{
		language.Chinese: catalog.String("账号已存在"),
		language.English: catalog.String("account already exists"),
	}
	translates["username or password error"] = locales{
		language.Chinese: catalog.String("用户名或密码错误"),
		language.English: catalog.String("username or password error"),
	}
	translates["reset password failed"] = locales{
		language.Chinese: catalog.String("重置密码失败"),
		language.English: catalog.String("reset password failed"),
	}
	translates["license invalid"] = locales{
		language.Chinese: catalog.String("无效 License"),
		language.English: catalog.String("License invalid"),
	}
	translates["create captcha failed"] = locales{
		language.Chinese: catalog.String("创建验证码失败"),
		language.English: catalog.String("create captcha failed"),
	}
	translates["删除失败，该云帐号号正在执行扫描"] = locales{
		language.Chinese: catalog.String("删除失败，该云帐号正在执行扫描"),
		language.English: catalog.String("delete failed, the cloud account is performing scanning"),
	}

	registerCatalog(translates)
}

// register 注册多语言
func registerCatalog(translates map[string]locales) {
	fallback := language.MustParse(LangZH)
	cl := catalog.NewBuilder(catalog.Fallback(fallback))
	for k, v := range translates {
		for tag, msg := range v {
			_ = cl.Set(tag, k, msg)
		}
	}
	message.DefaultCatalog = cl
}

// T translate
func T(lang string, key string, a ...any) string {
	langTag := message.MatchLanguage(lang)

	zap.L().Debug("i18n.T", zap.String("key", key),
		zap.String("langTag", langTag.String()))

	p := message.NewPrinter(langTag)
	return p.Sprintf(key, a...)
}
