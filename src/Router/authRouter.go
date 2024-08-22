package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// check if user is logged in...
func AuthMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userID")

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	c.Next()
}

func LoginEndpoint(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// check if username or password is empty
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "username or password cannot be empty",
			})
		return
	}

	// check the database
}
