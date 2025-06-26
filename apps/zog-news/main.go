package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"zog-news/config"
	"zog-news/database"
	"zog-news/domain"
	"zog-news/internal/repository/postgres"
	"zog-news/internal/rest"
	"zog-news/internal/rest/middleware"
	"zog-news/internal/validator"
	"zog-news/service"
    _ "zog-news/docs"

	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
    "github.com/swaggo/echo-swagger"
)

func init() {
	config.LoadEnv()
}

//	@title			Zero One Group News
//	@version		1.0
//	@description	Example Go API using Zero One Group's monorepo template

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {

	env := os.Getenv("APP_ENVIRONMENT")
	var handler slog.Handler

	w := os.Stdout
	if env == "local" {
		handler = tint.NewHandler(w, &tint.Options{
			ReplaceAttr: middleware.ColorizeLogging,
		})
	} else {
		// or continue setup log for another env
		handler = slog.NewTextHandler(w, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	dbPool, err := database.SetupPgxPool()
	if err != nil {
		slog.Error("Failed to set up database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	e := echo.New()
	e.HideBanner = true

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(0)

    e.Validator = validator.NewValidator()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	tp, shutdown := config.InitTracer(ctx)
	defer shutdown(ctx)

	e.Use(middleware.AttachTraceProvider(tp))
	e.Use(middleware.SlogLoggerMiddleware())
	e.Use(middleware.Cors())

	// Register the routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, domain.Response{
			Code:    200,
			Status:  "Succes",
			Message: "All is well!",
		})
	})

    e.GET("/swagger/*", echoSwagger.WrapHandler)

	articleRepo := postgres.NewArticleRepository(dbPool)
	articleService := service.NewArticleService(articleRepo)

	topicRepo := postgres.NewTopicRepository(dbPool)
	topicService := service.NewTopicService(topicRepo)

	apiV1 := e.Group("/api/v1")
	articlesGroup := apiV1.Group("")
    topicsGroup := apiV1.Group("")

	rest.NewArticleHandler(articlesGroup, articleService)
    rest.NewTopicHandler(topicsGroup, topicService)

	// Get host from environment variable, default to 127.0.0.1 if not set
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Get port from environment variable, default to 8000 if not set
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Server address and port to listen on
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	go func() {
		slog.Info("Server starting", "address", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Shutting down server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
	}
}
