package router

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-chat.io/src/hubManager"
	"go-chat.io/src/redisManager"
)

func HubRegisterEndpoint(hubTable *hubManager.HubTable, redisManager *redisManager.RedisManager) func(*gin.Context) {
	return func(c *gin.Context) {
		hubid_query, ok := c.Params.Get("hubId")
		if !ok {
			log.Printf("Incorrect HubID %v", hubid_query)
			return
		}

		hubID, err := strconv.ParseUint(hubid_query, 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Incorrect hubId...",
			})
			return
		}

		err = hubTable.Add(hubID, redisManager)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Generated Hub with hubID %v", hubID),
		})
	}
}
