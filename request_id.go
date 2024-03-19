/**
 * @author Jose Nidhin
 */
package main

import (
	"net/http"

	"github.com/google/uuid"
)

const DefaultRequestIdHeader = "X-Request-Id"

func RequestId() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return requestContextFn(next)
	}
}

func requestContextFn(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		requestId := "Ooops"
		id, err := uuid.NewRandom()
		if err == nil {
			requestId = id.String()
		}
		ctx := SetCtxString(req.Context(), ReqIdCtxKey, requestId)

		req = req.WithContext(ctx)
		res.Header().Add(DefaultRequestIdHeader, requestId)

		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)
}
