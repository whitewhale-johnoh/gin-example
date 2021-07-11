package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var user1 = User{
	ID:       "jhoh",
	Username: "johnoh",
	Password: "password",
}

func main() {
	ConnectToDB()
	defer db.Close()
	r := gin.Default()
	r.Use(JSONMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/health", health)
		v1.POST("/signup", signup)
		v1.POST("/login", login)
	}
	r.Run(":8000")
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func health(c *gin.Context) {
	tokenAuth, err := ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := FetchAccessAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": userId + " accessed",
	})
}

func signup(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	added, err := AddUrMyUser(u.ID, u.Password, u.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid signup details")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": added,
	})
}

func login(c *gin.Context) {

	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	auth := GetUrMyUser(u.ID, u.Password)
	if !auth {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := CreateToken(user1.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAccessAuth(user1.ID, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	tokens := map[string]string{
		"access_token": token.AccessToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
		"token":   tokens,
	})
}
