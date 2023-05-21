package main

import (
	"net/http"

	"context"

	"github.com/ell1jah/linkcompress/internal/http/middleware"
	"github.com/ell1jah/linkcompress/internal/http/service"
	"github.com/ell1jah/linkcompress/internal/http/transport"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const port = ":8081"

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	ctx, finish := context.WithCancel(context.Background())
	defer func() {
		finish()
	}()

	service := service.NewMicroserviceClient(logger, ctx)
	handler := transport.NewLinkHandler(service, logger)

	r := mux.NewRouter()

	r.HandleFunc("/", handler.Post).Methods("POST")
	r.HandleFunc("/", handler.Get).Methods("GET")

	mux := middleware.AccessLog(logger, r)
	mux = middleware.Panic(logger, mux)

	logger.Infow("starting server",
		"port", port,
	)

	logger.Errorln(http.ListenAndServe(port, mux))
}
