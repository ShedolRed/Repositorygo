package handlers

import (
	"context"
	"net/http"
	"time"

	"go-mongo-railway-api/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateItem(db *mongo.Database) gin.HandlerFunc {
	col := db.Collection("items")
	return func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		now := time.Now().Unix()
		item := models.Item{
			Data:      payload,
			CreatedAt: now,
			UpdatedAt: now,
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := col.InsertOne(ctx, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
			return
		}
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			item.ID = oid
		}

		c.JSON(http.StatusCreated, item)
	}
}

func ListItems(db *mongo.Database) gin.HandlerFunc {
	col := db.Collection("items")
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		cur, err := col.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
			return
		}
		defer cur.Close(ctx)

		var out []models.Item
		if err := cur.All(ctx, &out); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decode failed"})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

func GetItem(db *mongo.Database) gin.HandlerFunc {
	col := db.Collection("items")
	return func(c *gin.Context) {
		idHex := c.Param("id")
		oid, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var item models.Item
		if err := col.FindOne(ctx, bson.M{"_id": oid}).Decode(&item); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func UpdateItem(db *mongo.Database) gin.HandlerFunc {
	col := db.Collection("items")
	return func(c *gin.Context) {
		idHex := c.Param("id")
		oid, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
			return
		}

		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"data":            payload,
				"updated_at_unix": time.Now().Unix(),
			},
		}

		res, err := col.UpdateOne(ctx, bson.M{"_id": oid}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func DeleteItem(db *mongo.Database) gin.HandlerFunc {
	col := db.Collection("items")
	return func(c *gin.Context) {
		idHex := c.Param("id")
		oid, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := col.DeleteOne(ctx, bson.M{"_id": oid})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed"})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
