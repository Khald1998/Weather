package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func main() {
	router := gin.Default()
	const ContentTypeHTML = "./site/*"
	router.LoadHTMLGlob(ContentTypeHTML)
	router.GET("/", home)
	router.POST("/search", search)

	router.Run("localhost:8080")
}

func home(x *gin.Context) { x.HTML(http.StatusOK, "index.html", nil) }

func search(x *gin.Context) {
	x.Request.ParseForm()
	log.Println(x.Request.Form)
	data := x.Request.FormValue("City")
	fitchData := getData(data)
	temperature := gjson.Get(fitchData, "main.temp")
	feels_like := gjson.Get(fitchData, "main.feels_like")
	Max := gjson.Get(fitchData, "main.temp_max")
	Min := gjson.Get(fitchData, "main.temp_min")
	weather := gjson.Get(fitchData, "weather.0.main")
	weather_description := gjson.Get(fitchData, "weather.0.description")
	weather_icon := gjson.Get(fitchData, "weather.0.icon")

	templ, err := template.ParseFiles("./site/results.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templ.Execute(x.Writer, weather_icon.String())
	if err != nil {
		log.Println("Error = ", err)
	}
	log.Println("weather is ", weather.String())
	varToPass := gin.H{
		"City":                data,
		"temperature":         temperature.String(),
		"feels_like":          feels_like.String(),
		"Max":                 Max.String(),
		"Min":                 Min.String(),
		"weather":             weather.String(),
		"weather_description": weather_description.String(),
	}

	x.HTML(http.StatusOK, "results.html", varToPass)

}
func getData(c string) string {
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
	return string(resBody)

}
