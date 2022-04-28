package middleware

import (
	"net/http"

	token "github.com/akhil/ecommerce-yt/tokens"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		// Why do you add so many claims to your JWT token?
		// Is just the email or uid (which I assume are both unique) and fetch the user if you need to get more data from the user?
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
