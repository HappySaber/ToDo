package midlleware

import (
	"ToDo/utils"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(401, gin.H{"error": "Couldn't get cookie 'token'"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)

		if err != nil {
			c.JSON(401, gin.H{"error": "Не удалось разобрать токен: " + err.Error()})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}
