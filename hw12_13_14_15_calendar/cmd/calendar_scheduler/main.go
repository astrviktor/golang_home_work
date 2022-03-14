package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/queue/producer"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	config := config.NewSchedulerConfig(configFile)

	var st storage.Storage
	if config.Storage.Mode == storage.SQLMode {
		st = sqlstorage.New(config.Storage.DSN, config.Storage.MaxConnectAttempts)
	} else {
		st = memorystorage.New()
	}

	if err := st.Connect(context.Background()); err != nil {
		log.Fatal("error connecting to storage")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	repeat := time.Duration(config.AMQPScheduler.RepeatSecond) * time.Second
	ticker := time.NewTicker(repeat)

	go func() {
		for t := range ticker.C {
			log.Println("tick...", t)

			p, err := producer.NewProducer("scheduler", config.AMQPScheduler.URI, config.AMQPScheduler.Exchange)
			if err != nil {
				log.Println("error connecting to AMQP:", err.Error())
				continue
			}

			if err = p.Produce(st); err != nil {
				log.Println("error producing:", err.Error())
			}

			if err = p.Close(); err != nil {
				log.Println("error closing producer:", err.Error())
			}
		}
	}()

	<-ctx.Done()
	if err := st.Close(ctx); err != nil {
		log.Println("error closing storage:", err.Error())
	}
	log.Println("scheduler done")
}
