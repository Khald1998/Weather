package main

import (
	"encoding/json"
	"io/ioutil"
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
	fitchData := getData(data)
	log.Println("-------------------------------------")
	foomapone := fitchData["main"]
	temp := foomapone.(map[string]interface{})
	foomaptwo := fitchData["weather"]
	state := foomaptwo.([]map[string]interface{})

	log.Println(state[0])
	log.Println("-------------------------------------")
	varToPass := gin.H{
		"City":        data,
		"temperature": temp["temp"],
		"feels_like":  temp["feels_like"],
		"Max":         temp["temp_max"],
		"Min":         temp["temp_min"],
	}

	x.HTML(http.StatusOK, "results.html", varToPass)

}
func getData(c string) map[string]interface{} {
	posturl := "https://api.openweathermap.org/data/2.5/weather"
	city := c
	API := "51c242b3396833f6078589fc5411066e"
	url := posturl + "?q=" + city + "&appid=" + API

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
	}
	var jsonRes map[string]interface{}    // declaring a map for key names as string and values as interface
	_ = json.Unmarshal(resBody, &jsonRes) // Unmarshalling

	log.Println(jsonRes)
	return jsonRes

}
