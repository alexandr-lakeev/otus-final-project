package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexandr-lakeev/otus-final-project/internal/app"
	deliveryhttp "github.com/alexandr-lakeev/otus-final-project/internal/app/delivery/http"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	"github.com/gorilla/mux"
)

type Server struct {
	server  *http.Server
	router  *mux.Router
	usecase app.UseCase
	logg    app.Logger
}

func NewServer(cfg config.ServerConf, usecase app.UseCase, logger app.Logger) *Server {
	handler := deliveryhttp.NewHandler(usecase, logger)

	router := mux.NewRouter()
	router.Use(newLoggingMiddleware(logger))
	router.PathPrefix("/fill").Handler(handler.Fill(context.Background())).Methods("GET")

	return &Server{
		server: &http.Server{
			Handler:      router,
			Addr:         cfg.BindAddress,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		logg:    logger,
		usecase: usecase,
		router:  router,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info(fmt.Sprintf("server is starting on %s", s.server.Addr))

	return s.server.ListenAndServe()
}

// for tests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
