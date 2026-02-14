package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vector-ops/chisai/cmd"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	level := &slog.LevelVar{}

	env := os.Getenv("ENV")

	if env == "prod" {
		level.Set(slog.LevelInfo)
	}

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	logger := slog.New(logHandler)

	slog.SetDefault(logger)

	if err = cmd.InitDB(); err != nil {
		log.Fatal(err)
	}

	cancelCh := make(chan os.Signal, 1)

	signal.Notify(cancelCh, syscall.SIGINT, syscall.SIGTERM)

	server := cmd.NewServer(cmd.GetDB())
	go server.Start()

	<-cancelCh

	server.Close()
	cmd.Close(context.TODO())

}
