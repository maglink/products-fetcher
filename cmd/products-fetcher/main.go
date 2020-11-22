package main

import (
	"context"
	"flag"
	"github.com/maglink/products-fetcher/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "./configs/default.yaml", "config path")

func main() {
	flag.Parse()
	cfg, err := server.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed on load config %v", err)
	}

	ctx := initContext()
	srv := server.New(cfg, ctx)
	srv.Run()
}

func initContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigquit:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}
