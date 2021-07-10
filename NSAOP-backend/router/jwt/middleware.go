package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"nsaop/router/resp"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.ERROR(c, http.StatusUnauthorized, "authHeader not found")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp.ERROR(c, http.StatusUnauthorized, "authHeader format error")
			c.Abort()
			return
		}
		accessToken := parts[1]
		j := NewJWT("access")
		claims, err := j.ParseToken(accessToken)
		if err != nil {
			resp.ERROR(c, http.StatusUnauthorized, err.Error())
			c.Abort()
		} else {
			c.Set("userId", claims.UserId)
			c.Set("userRole", claims.UserRole)
			c.Next()
		}
	}
}
