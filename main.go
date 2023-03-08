package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// Define constants for various values used throughout the code
const (
	contentTypeHTML = "./site/*"                                        // Path to HTML templates directory
	openWeatherAPI  = "https://api.openweathermap.org/data/2.5/weather" // OpenWeather API URL
	apiKey          = "51c242b3396833f6078589fc5411066e"                // API key for OpenWeather API
	addr            = "localhost:8080"                                  // Server address
)

// Setup the Gin router with necessary routes and middleware
func setupRouter() *gin.Engine {
	router := gin.Default() // Create a new Gin router instance with default middleware

	router.LoadHTMLGlob(contentTypeHTML) // Load HTML templates from the specified directory

	// Serve static files from the specified directory under the "/static" route
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// Add handlers for the root path, the "Search" path and the "Results" path
	router.GET("/", home)
	router.GET("/Search", Search)
	router.POST("/Results", Results)

	return router // Return the configured router instance
}

// Start the server
func main() {
	router := setupRouter() // Set up the Gin router

	err := router.Run(addr) // Start the server
	if err != nil {
		log.Fatalf("Failed to start server: %v", err) // If the server fails to start, log the error and exit
	}
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "Main.html", nil)
}
func Search(c *gin.Context) {
	c.HTML(http.StatusOK, "Search.html", nil)
}

func Results(c *gin.Context) {
	c.Request.ParseForm()
	City_name := c.Request.FormValue("City")
	if len(City_name) == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": "City name is missing",
		})
		return
	}

	var temperature, feelsLike, maxTemp, minTemp, weather, weatherDesc, weatherIcon string
	data := getData(City_name)

	if gjson.Get(data, "cod").String() != "200" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": gjson.Get(data, "message").String(),
		})
		return
	}
	temperature = gjson.Get(data, "main.temp").String()
	feelsLike = gjson.Get(data, "main.feels_like").String()
	maxTemp = gjson.Get(data, "main.temp_max").String()
	minTemp = gjson.Get(data, "main.temp_min").String()
	weather = gjson.Get(data, "weather.0.main").String()
	weatherDesc = gjson.Get(data, "weather.0.description").String()
	weatherIcon = gjson.Get(data, "weather.0.icon").String()

	varToPass := gin.H{
		"City":                City_name,
		"Temperature":         temperature,
		"Feels_like":          feelsLike,
		"Max":                 maxTemp,
		"Min":                 minTemp,
		"Weather":             weather,
		"Weather_description": weatherDesc,
		"Icon":                weatherIcon,
	}

	c.HTML(http.StatusOK, "Results.html", varToPass)

}

func getData(city string) string {
	// Define the base URL for the OpenWeatherMap API
	// Construct the URL for the API request
	url := openWeatherAPI + "?q=" + city + "&appid=" + apiKey

	// Create a new GET request with the constructed URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP client to make the request
	client := &http.Client{}

	// Send the request and store the response in 'resp'
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing the response body until the function returns
	defer resp.Body.Close()

	// Read the response body into a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice to a string and return it
	return string(body)
}
