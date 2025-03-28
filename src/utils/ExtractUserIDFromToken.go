package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ExtractUserEmailFromToken(c *gin.Context) (string, error) {
	userEmail, exists := c.Get("email")

	if !exists {
		return "", errors.New("user email not found in context")
	}

	return userEmail.(string), nil
}
