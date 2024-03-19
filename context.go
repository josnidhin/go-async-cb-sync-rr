/**
 * @author Jose Nidhin
 */
package main

import (
	"context"
)

type CtxKey int8

const (
	UnknownCtxKey CtxKey = iota
	ReqIdCtxKey
)

func SetCtxString(ctx context.Context, key CtxKey, val string) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetCtxString(ctx context.Context, key CtxKey) string {
	id := ctx.Value(key)
	str, _ := id.(string)
	return str
}
