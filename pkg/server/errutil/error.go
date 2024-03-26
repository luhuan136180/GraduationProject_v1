package errutil

import "net/http"

type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewError returns a ServiceError using the code and reason
func NewError(code int, message string, data ...any) ServiceError {
	err := ServiceError{Code: code, Message: message}
	if len(data) > 0 {
		err.Data = data[0]
	}
	return err
}

// Error returns a text representation of the service error
func (s ServiceError) Error() string {
	return s.Message
}

var (
	ErrInternalServer   = NewError(http.StatusInternalServerError, "internal server error") // 服务器内部错误
	ErrIllegalParameter = NewError(http.StatusBadRequest, "illegal parameter")              // 非法参数
	ErrNotFound         = NewError(http.StatusNotFound, "not found")                        // 数据未找到
	ErrUnauthorized     = NewError(http.StatusUnauthorized, "unauthorized")                 // 会话过期，请重新登录
	ErrPermissionDenied = NewError(http.StatusForbidden, "permission denied")               // 权限不足
	ErrIllegalOperation = NewError(http.StatusBadRequest, "illegal operation ")             // 非法操作
	ErrDeleteUser       = NewError(http.StatusBadRequest, "delete user failed")             // 删除用户失败
	ErrCreateUser       = NewError(http.StatusBadRequest, "create user failed")             // 创建用户失败
	ErrEditUserInfo     = NewError(http.StatusBadRequest, "change user Info failed")        // 修改用户信息失败
	ErrSelectUserList   = NewError(http.StatusBadRequest, "select userList failed")         // 查询用户列表失败
	ErrChangeUserPWD    = NewError(http.StatusBadRequest, "reset user password failed")     // 重置用户密码失败
	ErrGetAuditLogs     = NewError(http.StatusBadRequest, "get audit logs failed")          // 获取审计日志失败
	ErrEditCredential   = NewError(http.StatusBadRequest, "edit credential info failed")    // 修改云帐号信息失败
	ErrInvalidLicense   = NewError(http.StatusBadRequest, "license invalid")                // 验证码错误
)
