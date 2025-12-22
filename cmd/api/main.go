package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-mongo-railway-api/internal/config"
	"go-mongo-railway-api/internal/db"
	"go-mongo-railway-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	client, err := db.ConnectMongo(cfg.MongoURL)
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	database := client.Database(cfg.MongoDBName)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Public routes
	r.GET("/health", handlers.Health(client, cfg.MongoDBName))

	// Simple Items CRUD (Mongo collection: items)
	r.POST("/items", handlers.CreateItem(database))
	r.GET("/items", handlers.ListItems(database))
	r.GET("/items/:id", handlers.GetItem(database))
	r.PUT("/items/:id", handlers.UpdateItem(database))
	r.DELETE("/items/:id", handlers.DeleteItem(database))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server started on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("shutdown complete")
}
