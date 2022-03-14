package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	config := config.NewCalendarConfig(configFile)
	logg := logger.New(config.Logger.Level, config.Logger.TimeFormat)

	var st storage.Storage
	if config.Storage.Mode == storage.SQLMode {
		st = sqlstorage.New(config.Storage.DSN, config.Storage.MaxConnectAttempts)
	} else {
		st = memorystorage.New()
	}

	if err := st.Connect(context.Background()); err != nil {
		logg.Fatal("error connecting to storage")
	}

	calendar := app.New(logg, st)

	httpServer := internalhttp.NewServer(logg, calendar, st, config.HTTPServer.Host, config.HTTPServer.Port)
	grpcServer := internalgrpc.NewServer(logg, calendar, st, config.GRPCServer.Host, config.GRPCServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go httpServer.Start(ctx)
	go grpcServer.Start(ctx)

	<-ctx.Done()
	logg.Info("calendar is done...")
}
