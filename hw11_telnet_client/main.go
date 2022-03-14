package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("not enough arguments host port")
	}

	client := NewTelnetClient(net.JoinHostPort(flag.Args()[0], flag.Args()[1]), timeout, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		return
	}

	done := make(chan int)
	go Done(done)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go SendOrReceive(wg, done, "send", client)
	go SendOrReceive(wg, done, "receive", client)

	wg.Wait()
}

func Done(done chan int) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	<-sigChan
	close(sigChan)
	close(done)
}

func SendOrReceive(wg *sync.WaitGroup, done chan int, sr string, client TelnetClient) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			if sr == "send" {
				if err := client.Send(); err != nil {
					return
				}
			}
			if sr == "receive" {
				if err := client.Receive(); err != nil {
					return
				}
			}
		}
	}
}
