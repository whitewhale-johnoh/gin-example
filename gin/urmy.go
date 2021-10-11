package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type SignInUser struct {
	LoginID       string `json:"loginId"`
	Password      string `json:"password"`
	Token         string `json:"notification_token"`
	UserDataExist bool   `json:"userdataexist"`
}

type SignUpUser struct {
	LoginID               string `json:"loginId"`
	Password              string `json:"password"`
	Name                  string `json:"name"`
	PhoneNo               string `json:"phoneNo"`
	Gender                bool   `json:"genderState"`
	Birthdate             string `json:"birthdate"`
	Country               string `json:"country"`
	Hometown              string `json:"hometown"`
	Phonecode             string `json:"phonecode"`
	Isoverage             bool   `json:"isoverage"`
	Urmyaccount           bool   `json:"urmyaccount"`
	Urmyoverallservice    bool   `json:"urmyoverallservice"`
	Urmynotiad            bool   `json:"urmynotiad"`
	Urmypersonaldataacc   bool   `json:"urmypersonaldataacc"`
	Urmylocation          bool   `json:"urmylocation"`
	Urmyprofileadditional bool   `json:"urmyprofileadditional"`
}

func init() {
	awsconfig = ConfigAWS()
}

func main() {
	ConnectToDB()
	ConnectToRedis()

	defer db.Close()
	defer rdb.Close()
	r := gin.Default()
	r.Use(JSONMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/signup", signup)
		v1.POST("/signin", signin)
	}

	v2 := r.Group("/v2")
	v2.Use(TokenAuthMiddleware)
	{
		v2.POST("/service", Service)
		v2.POST("/signout", signout)
	}
	r.Run(":3010")
}

func Service(c *gin.Context) {
	fmt.Println("ping")
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0")
		c.Writer.Header().Set("Last-Modified", time.Now().String())
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "-1")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func TokenAuthMiddleware(c *gin.Context) {
	tokenAuth, err := ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}
	userId, err := FetchAccessAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized,
				gin.H{
					"status": http.StatusUnauthorized,
					"error":  "token is expired",
				})
			c.Abort()
			return
		}
		c.JSON(http.StatusForbidden,
			gin.H{
				"status": http.StatusForbidden,
				"error":  "Authentication failed",
			})
		c.Abort()
		return

	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": userId + " accessed",
		})
		c.Next()
	}

}

func signup(c *gin.Context) {
	var u *SignUpUser
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		fmt.Println("Invalid json provided")
		return
	}
	uuid, err := AddUrMyUser(u)
	if err != nil {
		fmt.Println("invalid add user")
		fmt.Println(err.Error())
		c.JSON(http.StatusUnauthorized, "Please provide valid signup details")
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"message":        "created",
			"uuid":           uuid,
			"profilePicPath": uuid + "/profile/",
		})
	}
}

func signin(c *gin.Context) {
	var u SignInUser
	var urmyuuid string
	var autherr error
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if strings.Contains(u.LoginID, "@") {
		urmyuuid, autherr = GetUrMyUserEmail(&u)
		if autherr != nil {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
	} else {
		urmyuuid, autherr = GetUrMyUserPhone(&u)
		if autherr != nil {
			c.JSON(http.StatusUnauthorized, "Please provide valid login details")
			return
		}
	}

	token, err := CreateToken(urmyuuid)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := CreateAccessAuth(urmyuuid, token)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}

	userdata, userdataerr := GetUrMyUserData(urmyuuid)
	if userdataerr != nil {
		c.JSON(http.StatusUnauthorized, "User does not have proper information")
		return
	} else {
		userdata.AccessToken = token.AccessToken
		c.JSON(http.StatusOK, gin.H{
			"message":  "signin",
			"Userdata": userdata,
		})
	}

}

func signout(c *gin.Context) {

	tokenAuth, err := ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		fmt.Println("Etracting token invalid")
		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	deleted, deletederr := DeleteAccessAuth(tokenAuth.AccessUuid)
	if deletederr != nil || deleted == 0 {
		fmt.Println("token is not deleted")
		fmt.Println(deleted)
		fmt.Println(deletederr.Error())
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully logged out",
			"deleted": deleted,
		})
	}
}
