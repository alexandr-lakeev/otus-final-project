package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandr-lakeev/otus-final-project/internal/app/usecase"
	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	internalcache "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/cache"
	internalimage "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/image"
	internalloger "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/logger"
	internalhttp "github.com/alexandr-lakeev/otus-final-project/internal/infrastructure/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/previewer/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := internalloger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	cache := internalcache.NewCache(config.Previewer.CacheSize, config.Previewer.CacheDir)

	httpClient := &http.Client{
		Timeout: config.Previewer.RequestTimeout,
	}

	uc := usecase.New(
		internalimage.NewLoader(httpClient),
		internalimage.NewResizer(),
		cache,
		logger,
	)

	server := internalhttp.NewServer(config.Server, uc, logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("previewer is running...")

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
