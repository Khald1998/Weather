package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	const ContentTypeHTML = "./site/*"
	router.LoadHTMLGlob(ContentTypeHTML)
	router.GET("/", home)

	router.Run("localhost:8080")
}

func home(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", map[string]string{"title": "home page"})
	fmt.Println(c.Writer)

}
