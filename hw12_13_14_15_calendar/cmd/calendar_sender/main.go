package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/queue/consumer"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/version"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.yaml", "Path to configuration file")
}

// http://localhost:15672/ guest:guest
func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		version.PrintVersion()
		return
	}

	config := config.NewSenderConfig(configFile)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	c, err := consumer.NewConsumer("sender", config.AMQPSender.URI, config.AMQPSender.Queue)
	for i := 0; i < config.AMQPSender.Retry && err != nil; i++ {
		time.Sleep(5 * time.Second)
		log.Println("sender connection error:", err.Error())
		log.Println("reconnect...")
		c, err = consumer.NewConsumer("sender", config.AMQPSender.URI, config.AMQPSender.Queue)
	}
	if err != nil {
		log.Fatal("sender connection failed")
	}

	msgs, err := c.Consume(ctx)

	if err == nil {
		log.Println("sender start consuming...")

		for m := range msgs {
			log.Println("receive new message: ", string(m.Data))
		}

		if err = c.Close(); err != nil {
			log.Println("error closing sender:", err.Error())
		}
		log.Println("sender done")
	}
}
