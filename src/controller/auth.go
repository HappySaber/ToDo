package controllers

import (
	"ToDo/database"
	"ToDo/models"
	"ToDo/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("sosal?da!")

func Login(c *gin.Context) {
	var user models.Auth

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Auth

	query := "SELECT email, password FROM users WHERE email = $1"

	err := database.DB.QueryRow(query, user.Email).Scan(&existingUser.Email, &existingUser.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
		return
	}

	//deleting spaces from passwords, if they some way managed to be
	user.Password = strings.TrimSpace(user.Password)
	existingUser.Password = strings.TrimSpace(existingUser.Password)

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &models.Claims{
		Email: existingUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   existingUser.Email,
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "could not create token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"succes": "user logged in"})
}

func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//var existingUser models.User

	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`
	rows, err := database.DB.Query(query, user.Email)

	if err != nil {
		log.Println("Error executing query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	if rows.Next() {
		var id int
		var username, email, password, createdAt string
		if err := rows.Scan(&id, &username, &email, &password, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user data"})
			return
		}
		if id != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
	}

	var errHash error

	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	query = "INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, NOW())"
	if _, err := database.DB.Exec(query, user.Username, user.Email, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "user created"})
}

func Home(c *gin.Context) {
	cookie, err := c.Cookie("token")

	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized 1"})
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized  2"})
		return
	}

	if claims.Email == "" {
		c.JSON(401, gin.H{"error": "unauthorized 3"})
		return
	}

	c.JSON(200, gin.H{"success": "home page", "user: ": claims.Email})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}
