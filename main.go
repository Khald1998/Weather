package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	City string `json:"City"`
}

func main() {
	router := gin.Default()
	const ContentTypeHTML = "./site/*"
	router.LoadHTMLGlob(ContentTypeHTML)
	router.GET("/", home)
	router.POST("/search", search)

	router.Run("localhost:8080")
}

func home(x *gin.Context) {

	x.HTML(http.StatusOK, "index.html", nil)
	log.Println(x.Request.Body)

}
func search(x *gin.Context) {
	x.Request.ParseForm()
	log.Println(x.Request.Form)
	data := x.Request.FormValue("City")
	getData()

	varToPass := gin.H{
		"City": data,
	}

	x.HTML(http.StatusOK, "results.html", varToPass)

}
func getData(x *gin.Context) (r *http.Request) {
	posturl := "https://jsonplaceholder.typicode.com/posts"
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	return r

}
