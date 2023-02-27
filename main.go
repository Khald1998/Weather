package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type htmlvar struct {
	city string
}

func main() {
	router := gin.Default()
	const ContentTypeHTML = "./site/*"
	router.LoadHTMLGlob(ContentTypeHTML)
	router.GET("/", home)

	router.Run("localhost:8080")
}

func home(c *gin.Context) {
	varToPass := gin.H{
		"city": "al dhahran",
	}
	c.HTML(http.StatusOK, "index.html", varToPass)
	fmt.Println(c.Writer)

}
