package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"repositorygo/internal/config"
	"repositorygo/internal/db"
	"repositorygo/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	// Чтобы по ссылке Railway сразу было видно, что сервис живой
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Repositorygo API",
			"status":  "ok",
			"health":  "/health",
			"items":   "/items",
			"version": "1.0",
		})
	})

	r.GET("/health", handlers.Health)

	client, database, err := db.ConnectMongo(context.Background(), cfg)
	if err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	h := handlers.NewItems(database)

	r.GET("/items", h.List)
	r.POST("/items", h.Create)
	r.GET("/items/:id", h.GetByID)
	r.PUT("/items/:id", h.UpdateByID)
	r.DELETE("/items/:id", h.DeleteByID)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
