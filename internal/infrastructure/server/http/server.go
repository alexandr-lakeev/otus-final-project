package internalhttp

import (
	"context"
	"net/http"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	deliveryhttp "github.com/alexandr-lakeev/otus-final-project/internal/app/delivery/http"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	"github.com/gorilla/mux"
)

func NewServer(cfg config.ServerConf, usecase app.UseCase, logger app.Logger) *http.Server {
	handler := deliveryhttp.NewHandler(usecase, logger)

	router := mux.NewRouter()
	router.Use(newLoggingMiddleware(logger))
	router.PathPrefix("/fill").Handler(handler.Fill(context.Background())).Methods("GET")

	return &http.Server{
		Handler:      router,
		Addr:         cfg.BindAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
