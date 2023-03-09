package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Pallinder/go-randomdata"
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
	router.GET("/", homePage)
	router.GET("/Search", SearchPage)
	router.POST("/Results", ResultsPage)

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

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "Main.html", nil)

}
func SearchPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Search.html", nil)

}

func ResultsPage(c *gin.Context) {

	c.Request.ParseForm()
	City_name := c.Request.FormValue("City")
	if c.Request.FormValue("button") == "Search" {
		Search(City_name, c)
	} else if c.Request.FormValue("button") == "Random" {
		Random(c)
	}
}
func Random(c *gin.Context) {
	City_name := randomdata.City()
	data := getData(City_name)
	if gjson.Get(data, "message").String() != "" {
		// loop that runs infinitely
		for {
			City_name = randomdata.City()
			data = getData(City_name)
			// condition to terminate the loop
			if gjson.Get(data, "message").String() == "" {
				break
			}
		}
		Search(City_name, c)
	} else {
		Search(City_name, c)
	}

}
func Search(City_name string, c *gin.Context) {
	if len(City_name) == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": "City name is missing",
		})
		return
	}

	data := getData(City_name)

	if gjson.Get(data, "cod").String() != "200" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": gjson.Get(data, "message").String(),
		})
		return
	}

	varToPass := gin.H{
		"City":                City_name,
		"Temperature":         gjson.Get(data, "main.temp").String(),
		"Feels_like":          gjson.Get(data, "main.feels_like").String(),
		"Max":                 gjson.Get(data, "main.temp_max").String(),
		"Min":                 gjson.Get(data, "main.temp_min").String(),
		"Weather":             gjson.Get(data, "weather.0.main").String(),
		"Weather_description": gjson.Get(data, "weather.0.description").String(),
		"Icon":                gjson.Get(data, "weather.0.icon").String(),
	}

	c.HTML(http.StatusOK, "Results.html", varToPass)

}
func getData(city string) string {
	// Define the base URL for the OpenWeatherMap API
	// Construct the URL for the API request
	url := openWeatherAPI + "?q=" + city + "&appid=" + apiKey + "&units=metric"

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
