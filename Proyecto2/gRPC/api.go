package main

import(
	"net/http"
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
)

type message struct {
	Message string `json:"message"`
	Country string `json:"country"`
}

type result struct {
	Proccess string `json:"proccess"`
}

func main() {
	router := gin.Default()
	router.GET("/health", getHealth)
	router.POST("/incert", postMatch)
	router.Run(":3000")
}

func getHealth(c *gin.Context) {
	var nerResult result
	newResult.Proccess = "done"
	c.IndentedJSON(http.StatusOK, newResult)
}

func postMatch(c *gin.Context) {
	var newMessage message
	if err := c.BindJSON(&newMessage); err != nil {
		return
	}

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	})
	err := rdb.HIncrBy(ctx, "countries", newMessage.Country, 1).Err()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newMessage)
}
