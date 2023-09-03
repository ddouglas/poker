package internal

import (
	"context"
	"poker"
)

type contextKey uint

const (
	userCtxKey contextKey = iota
)

func ContextWithUser(ctx context.Context, user *poker.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func UserFromContext(ctx context.Context) *poker.User {

	userInf := ctx.Value(userCtxKey)
	if userInf == nil {
		return nil
	}

	if user, ok := userInf.(*poker.User); ok {
		return user
	}

	return nil

}
