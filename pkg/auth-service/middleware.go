package authservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth(ctrl *UserController) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessTk := c.GetHeader("Authorization")

		if accessTk == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "no authorization token provided",
			})
			return
		}

		usrInfo, err := ctrl.TokenValSvcClient.ValidateUsrAccessToken(accessTk)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		c.Set("user", *usrInfo)

		c.Next()
	}
}
