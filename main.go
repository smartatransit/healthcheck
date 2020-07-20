package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	flags "github.com/jessevdk/go-flags"
	"github.com/smartatransit/healthcheck/pkg/health"
	"go.uber.org/zap"
)

type options struct {
	ConfigPath string `long:"config-path" env:"CONFIG_PATH" description:"An optional file that overrides the default configuration of sources and targets." required:"true"`
	Debug      bool   `long:"debug" env:"DEBUG" description:"enabled debug logging"`
}

func main() {
	fmt.Println("Starting healthcheck service")
	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	var logger *zap.Logger
	if opts.Debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer func() {
		_ = logger.Sync() // flushes buffer, if any
	}()

	cfg, err := health.NewConfig(opts.ConfigPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	cc := health.NewCheckClient(cfg, logger)
	ctx, cancelFunc := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cc.Start(ctx)

	select {
	case <-quit:
		cancelFunc()
		logger.Info("interrupt signal received")
		logger.Info("shutting down...")
	}
}
