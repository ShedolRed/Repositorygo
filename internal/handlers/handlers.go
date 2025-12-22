package handlers

import (
    "repositorygo/internal/config"

    "go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
    db  *mongo.Database
    cfg config.Config
}

func New(db *mongo.Database, cfg config.Config) *Handler {
    return &Handler{db: db, cfg: cfg}
}

func (h *Handler) coll() *mongo.Collection {
    return h.db.Collection("items")
}
