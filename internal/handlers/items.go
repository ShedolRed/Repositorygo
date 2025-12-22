package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"

    "repositorygo/internal/models"
)

func (h *Handler) ListItems(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
    defer cancel()

    cur, err := h.coll().Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cur.Close(ctx)

    items := make([]models.Item, 0)
    for cur.Next(ctx) {
        var it models.Item
        if err := cur.Decode(&it); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        items = append(items, it)
    }

    c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateItem(c *gin.Context) {
    var req models.CreateItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    now := time.Now().UTC()
    doc := models.Item{
        ID:        primitive.NewObjectID(),
        Name:      req.Name,
        Value:     req.Value,
        CreatedAt: now,
        UpdatedAt: now,
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
    defer cancel()

    _, err := h.coll().InsertOne(ctx, doc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, doc)
}

func (h *Handler) GetItem(c *gin.Context) {
    oid, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
    defer cancel()

    var it models.Item
    if err := h.coll().FindOne(ctx, bson.M{"_id": oid}).Decode(&it); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.JSON(http.StatusOK, it)
}

func (h *Handler) UpdateItem(c *gin.Context) {
    oid, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req models.UpdateItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    set := bson.M{}
    if req.Name != nil {
        set["name"] = *req.Name
    }
    if req.Value != nil {
        set["value"] = *req.Value
    }
    if len(set) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
        return
    }
    set["updatedAt"] = time.Now().UTC()

    ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
    defer cancel()

    res, err := h.coll().UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if res.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    var it models.Item
    _ = h.coll().FindOne(ctx, bson.M{"_id": oid}).Decode(&it)
    c.JSON(http.StatusOK, it)
}

func (h *Handler) DeleteItem(c *gin.Context) {
    oid, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
    defer cancel()

    res, err := h.coll().DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if res.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    c.Status(http.StatusNoContent)
}
