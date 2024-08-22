package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	redis_session "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	router "go-chat.io/src/Router"
	"go-chat.io/src/db"
	"go-chat.io/src/hubManager"
	"go-chat.io/src/redisManager"
)

func main() {

	// setup dotenv
	err := godotenv.Load(".config/.env.development")
	if err != nil {
		log.Panic("Error occur while loading dotenv file")
	}

	// setup redis manager
	redisManager := redisManager.NewRedisManager()

	// setup DB
	DB := db.DB_Init()
	log.Printf("%v", DB)

	// setup hubTable
	hubTable := hubManager.NewHubTable()

	// setup redis session store
	redis_url := os.Getenv("REDIS_URL")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_db := os.Getenv("REDIS_DB")
	session_secret := os.Getenv("SESSION_SECRET")
	redisStore, _ := redis_session.NewStoreWithDB(1000, "tcp", redis_url, redis_password, redis_db, []byte(session_secret))

	r := gin.Default()

	// register redis session store
	r.Use(sessions.Sessions("auth-session", redisStore))

	r.GET("/hc", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is Running!",
		})
	})

	// generate a hub
	r.POST("/hub/generate/:hubId", router.HubRegisterEndpoint(hubTable, redisManager))

	// connect to hub
	r.GET("/ws/:hubId", router.ChatEndpoint(hubTable))

	r.Run(":8080")
}
