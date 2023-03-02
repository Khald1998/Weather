package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Pallinder/go-randomdata"
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

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func search(c *gin.Context) {
	c.Request.ParseForm()
	City_name := c.Request.FormValue("City")
	if City_name == "" {
		City_name = randomdata.City()
	}

	Data := getData(City_name)
	if gjson.Get(Data, "cod").String() == "200" {
		temperature := gjson.Get(Data, "main.temp")
		feels_like := gjson.Get(Data, "main.feels_like")
		Max := gjson.Get(Data, "main.temp_max")
		Min := gjson.Get(Data, "main.temp_min")
		weather := gjson.Get(Data, "weather.0.main")
		weather_description := gjson.Get(Data, "weather.0.description")
		weather_icon := gjson.Get(Data, "weather.0.icon")
	} else {
		temperature := gjson.Get(Data, "main.temp")
		feels_like := gjson.Get(Data, "main.feels_like")
		Max := gjson.Get(Data, "main.temp_max")
		Min := gjson.Get(Data, "main.temp_min")
		weather := gjson.Get(Data, "weather.0.main")
		weather_description := gjson.Get(Data, "weather.0.description")
		weather_icon := gjson.Get(Data, "weather.0.icon")
	}

	varToPass := gin.H{
		"City":                City_name,
		"Temperature":         temperature.String(),
		"Feels_like":          feels_like.String(),
		"Max":                 Max.String(),
		"Min":                 Min.String(),
		"Weather":             weather.String(),
		"Weather_description": weather_description.String(),
		"Icon":                weather_icon.String(),
	}

	c.HTML(http.StatusOK, "results.html", varToPass)

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
