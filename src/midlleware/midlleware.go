package midlleware

import (
	"ToDo/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized 1"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)

		if err != nil {
			log.Println("Ошибка при разборе токена:", err)
			c.JSON(401, gin.H{"error": "unauthorized 2"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}
