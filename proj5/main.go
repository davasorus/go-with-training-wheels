package main

import (
	database "example/web-service-gin/Database"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
	albums, err := database.AllAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Error fetching albums", "error", err)
		return
	}
	slog.Info("Fetched albums successfully", "count", len(albums))
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := database.AddAlbum(database.Album{
		Title:  newAlbum.Title,
		Artist: newAlbum.Artist,
		Price:  newAlbum.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Error adding album", "error", err)
		return
	}
	slog.Info("Added album successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		slog.Error("Invalid album ID", "id", id, "error", err)
		return
	}
	album, err := database.AlbumByID(intID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		slog.Error("Album not found", "id", intID, "error", err)
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func main() {

	database.ConnectDatabase()
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}
