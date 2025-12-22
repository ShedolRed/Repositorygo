package main

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"

    "repositorygo/internal/config"
    "repositorygo/internal/db"
    "repositorygo/internal/handlers"
)

func main() {
    cfg := config.Load()

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    mongoClient, mongoDB, err := db.ConnectMongo(ctx, cfg)
    if err != nil {
        // Railway перезапустит контейнер, если приложение упало.
        panic(err)
    }
    defer func() { _ = mongoClient.Disconnect(context.Background()) }()

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    h := handlers.New(mongoDB, cfg)

    r.GET("/health", h.Health)

    r.GET("/items", h.ListItems)
    r.POST("/items", h.CreateItem)
    r.GET("/items/:id", h.GetItem)
    r.PUT("/items/:id", h.UpdateItem)
    r.DELETE("/items/:id", h.DeleteItem)

    srv := &http.Server{
        Addr:              ":" + cfg.Port,
        Handler:           r,
        ReadHeaderTimeout: 10 * time.Second,
    }

    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}
