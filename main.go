package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Key string `bson:"key"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("DATABASE_URL")))

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("⛺️ Database has connected!")
	}
	defer client.Disconnect(context.Background())

	// Database insert schema
	// Accessing a database and collection
	collection := client.Database("go-joke-api").Collection("api_keys")

	r := gin.Default()

	r.GET("/joke", func(c *gin.Context) {

		apikey := c.Query("api-key")

		if apikey == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Missing api_key as a query",
			})
			return
		}

		// Define a filter to search for the API key
		filter := bson.M{"key": apikey}

		// Find key in the database
		var result API_KEY
		err = collection.FindOne(context.Background(), filter).Decode(&result)

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid api key",
			})
			return
		}

		joke := GetJoke()

		c.JSON(http.StatusOK, gin.H{
			"joke": joke, // Use a string key here
		})
	})

	r.POST("/api-key", func(c *gin.Context) {
		apiKey := API_KEY{Key: CreateCode()}
		_, err = collection.InsertOne(context.Background(), apiKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not create an api key :/",
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"data": apiKey,
			})
		}
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

func CreateCode() string {
	// Define the length of the API key
	keyLength := 32 // You can adjust the length as needed

	// Create a byte slice to store random bytes
	randomBytes := make([]byte, keyLength)

	// Fill the slice with random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Handle error, if any
		panic(err)
	}

	// Encode the random bytes to a hexadecimal string
	key := hex.EncodeToString(randomBytes)

	// Return the API key
	return key
}
