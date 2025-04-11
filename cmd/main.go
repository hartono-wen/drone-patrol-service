package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hartono-wen/drone-patrol-service/config"
	"github.com/hartono-wen/drone-patrol-service/generated"
	"github.com/hartono-wen/drone-patrol-service/handler"
	"github.com/hartono-wen/drone-patrol-service/repository"
	"github.com/hartono-wen/drone-patrol-service/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Validator = validator.NewRequestValidator()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 seconds wait for graceful shutdown
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func newServer() *handler.Server {
	config, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	dbDsn := config.DatabaseURL
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	log.Println("Successfully initialized repository")
	opts := handler.NewServerOptions{
		Repository: repo,
		Config:     config,
	}
	return handler.NewServer(opts)
}
