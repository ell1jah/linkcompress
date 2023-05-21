package middleware

import (
	"net/http"
	"time"
)

type MiddlewareLogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

func Panic(logger MiddlewareLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorw("Server paniced",
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
					"url", r.URL.Path,
					"error", err,
				)

				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AccessLog(logger MiddlewareLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}
