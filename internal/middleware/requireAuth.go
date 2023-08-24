package middleware

import (
	"log"
	"net/http"
	jwthelper "senpainikolay/go-internship-smartdata/utils/jwt"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {

	tokenStr, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
	}

	log.Println(tokenStr)

	usrInfo, err := jwthelper.ValidateToken(tokenStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
	}

	c.Set("user", usrInfo)

	c.Next()

}
