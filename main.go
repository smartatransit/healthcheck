package main

import (
	"fmt"
	"log"

	flags "github.com/jessevdk/go-flags"
	"github.com/smartatransit/healthcheck/pkg/config"
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

	cfg, err := config.NewConfig(opts.ConfigPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("%+v", cfg))
}
