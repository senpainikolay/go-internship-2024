package middleware

import (
	"net/http"
	"senpainikolay/go-internship-smartdata/internal/auth"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {

	accessTk := c.GetHeader("Authorization")

	if accessTk == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "no authorization token provided",
		})
		return
	}

	usrInfo, err := auth.ValidateAccessToken(accessTk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.Set("user", usrInfo)

	c.Next()

}
