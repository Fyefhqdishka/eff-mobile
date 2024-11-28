package main

import (
	"fmt"
	"github.com/Fyefhqdishka/eff-mobile/internal/app"
	"github.com/Fyefhqdishka/eff-mobile/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Song Library API
// @version 1.0
// @description API для управления библиотекой песен
// @contact.url https://github.com/Fyefhqdishka
// @contact.email anuar.nassipov@gmail.com
// @BasePath /
// @schemes http
// @Accept json
// @Produce json
func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("can't load config, err: %v", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("can't load server, err: %v", err)
	}

	go func() {
		if err = app.Run(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()
	log.Println("shutting down...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := app.Stop(); err != nil {
		fmt.Errorf("error during shutdown: %v", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't load env file, err=%v", err)
	}
}
