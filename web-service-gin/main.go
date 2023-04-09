package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type album struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Artist string `json:"artist"`
	Price float32 `json:"price_in_euro"`
}

var albums = []album {
	{Id: "1", Title: "Test1", Artist: "Test1", Price: 42.42},
	{Id: "2", Title: "Test2", Artist: "Test2", Price: 42.42},
	{Id: "3", Title: "Test3", Artist: "Test3", Price: 42.42},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); nil != err {
		c.IndentedJSON(http.StatusBadRequest, "Could not decode to album!")
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.Id != id {
			continue
		}

		c.IndentedJSON(http.StatusOK, a)
		return
	}

	c.Status(http.StatusNotFound)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumById)

	err := router.Run("localhost:8080")
	if nil != err {
		panic(err)
	}
}
