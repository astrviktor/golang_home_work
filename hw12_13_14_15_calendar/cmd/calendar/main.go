package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.TimeFormat)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	server.Start(ctx)
}
