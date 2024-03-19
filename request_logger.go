/**
 * @author Jose Nidhin
 */
package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return requestLoggerFn(next, logger)
	}
}

func requestLoggerFn(next http.Handler, logger *slog.Logger) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		reqLogger := logger.With(slog.String("requestId", GetCtxString(req.Context(), ReqIdCtxKey)))

		w := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

		start := time.Now()
		defer func() {
			duration := time.Since(start)
			routePatten := chi.RouteContext(req.Context()).RoutePattern()
			status := w.Status()

			reqLogger.Info("HTTP request log",
				slog.String("httpVersion", req.Proto),
				slog.String("userAgent", req.Header.Get("User-Agent")),
				slog.String("method", req.Method),
				slog.String("route", routePatten),
				slog.String("url", req.RequestURI),
				slog.Int("status", status),
				slog.Duration("duration", duration))
		}()

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
