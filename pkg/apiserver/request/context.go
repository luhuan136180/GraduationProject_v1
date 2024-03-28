package request

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"v1/pkg/model"
	"v1/pkg/token"
)

type ctxKey string

const (
	ctxUserInfoKey ctxKey = "meta"
	ctxLanguageKey ctxKey = "language"
)

func WithTokenPayloadToCtx(ctx context.Context, info *token.Payload) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	return context.WithValue(ctx, ctxUserInfoKey, info)
}

func TokenPayloadFromCtx(ctx context.Context) (*token.Payload, error) {
	if ctx == nil {
		return nil, fmt.Errorf("ctx is nil")
	}

	val := ctx.Value(ctxUserInfoKey)
	if val == nil {
		return nil, fmt.Errorf("ctx meta info not found")
	}

	info, ok := val.(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("ctx meta info damaged")
	}

	return info, nil
}

func GetUsernameFromCtx(ctx context.Context) string {
	payload, err := TokenPayloadFromCtx(ctx)
	if err != nil {
		zap.L().Warn("GetUsernameFromCtx 解析token错误", zap.Error(err))
		return ""
	}

	return payload.Username
}

func GetNameFromCtx(ctx context.Context) string {
	payload, err := TokenPayloadFromCtx(ctx)
	if err != nil {
		zap.L().Warn("GetNameFromCtx 解析token错误", zap.Error(err))
		return ""
	}
	return payload.Name
}
func GetUserIdFromCtx(ctx context.Context) int64 {
	payload, err := TokenPayloadFromCtx(ctx)
	if err != nil {
		zap.L().Warn("GetUsernameFromCtx 解析token错误", zap.Error(err))
		return -1
	}
	return payload.ID
}

func GetUserUIDFromCtx(ctx context.Context) string {
	payload, err := TokenPayloadFromCtx(ctx)
	if err != nil {
		zap.L().Warn("GetUsernameFromCtx 解析token错误", zap.Error(err))
		return ""
	}
	return payload.UID
}

func GetRoleTypeFromCtx(ctx context.Context) model.RoleType {
	payload, err := TokenPayloadFromCtx(ctx)
	if err != nil {
		zap.L().Warn("GetRoleTypeFromCtx 解析token错误", zap.Error(err))
		return ""
	}

	return payload.Role
}

func WithLanguageToCtx(ctx context.Context, lang string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	return context.WithValue(ctx, ctxLanguageKey, lang)
}

func LanguageFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxLanguageKey)
	if v != nil {
		return v.(string)
	}

	return ""
}
