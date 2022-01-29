package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandr-lakeev/otus-final-project/internal/app/usecase"
	"github.com/alexandr-lakeev/otus-final-project/internal/cache"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	"github.com/alexandr-lakeev/otus-final-project/internal/image"
	"github.com/alexandr-lakeev/otus-final-project/internal/logger"
	internalhttp "github.com/alexandr-lakeev/otus-final-project/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar.dev.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	cache := cache.New()
	loader := image.NewLoader()
	uc := usecase.New(loader, cache, logger)

	server := internalhttp.NewServer(config.Server, uc, logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
