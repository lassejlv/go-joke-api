package main

import (
	"encoding/json"
	"os"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var API_URL string = "https://api.chucknorris.io/jokes/random"

type Joke struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	URL       string `json:"url"`
	IconURL   string `json:"icon_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type API_KEY struct {
	Key string
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.GET("/joke", func(c *gin.Context) {
		joke := GetJoke()

		c.JSON(http.StatusOK, gin.H{
			"joke": joke, // Use a string key here
		})
	})

	// Listen and serve on 0.0.0.0:8080
	port := ":" + os.Getenv("PORT")
	r.Run(port)
}

func GetJoke() Joke {
	res, err := http.Get(API_URL)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	var joke Joke
	err = json.NewDecoder(res.Body).Decode(&joke)

	if err != nil {
		log.Fatalln(err)
	}

	return joke
}
