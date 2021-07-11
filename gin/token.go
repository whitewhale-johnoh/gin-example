package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"gopkg.in/redis.v5"
)

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

type AccessTokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

var client *redis.Client

var mySigningKey = []byte("mysupersecretphrase")

func CreateToken(userid string) (*AccessTokenDetails, error) {
	td := &AccessTokenDetails{}
	td.AtExpires = time.Now().Add(time.Second * 240).Unix()
	u, err := uuid.NewV4()
	if err != nil {
		// TODO: Handle error.
		fmt.Println(err)
	}
	td.AccessUuid = u.String()

	atclaims := jwt.MapClaims{}
	atclaims["authorized"] = true
	atclaims["access_uuid"] = td.AccessUuid
	atclaims["user_id"] = userid
	atclaims["exp"] = td.AtExpires
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, atclaims)
	td.AccessToken, err = atoken.SignedString(mySigningKey)

	if err != nil {
		return nil, err
	}
	return td, nil

}

func ExtractAccessToken(r *http.Request) string {
	bearToken := r.Header.Get("AuthorizationAccess")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
}

func VerifyAccessToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractAccessToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractAccessTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyAccessToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		//userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		userId := claims["user_id"].(string)
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func FetchAccessAuth(authD *AccessDetails) (string, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", nil
	}
	//userID, _ := strconv.ParseUint(userid, 10, 64)

	return userid, nil
}

func CreateAccessAuth(userid string, td *AccessTokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}

func ResetAccessAuth(authD *AccessDetails) error {
	//converting Unix to UTC(to Time object)
	rt := time.Duration(time.Minute * 15)
	//.AtExpires = time.Now().Add(time.Minute * 15).Unix()

	errRefresh := client.Expire(authD.AccessUuid, rt).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func DeleteAccessAuth(givenAccessUuid string) (int64, error) {
	accessdeleted, err := client.Del(givenAccessUuid).Result()
	if err != nil {
		return 0, err
	}
	return accessdeleted, nil
}

func init() {
	//Initializing redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.160:6379",
		Password: "qwer1234", // no password set
		DB:       0,          // use default DB
	})
	ping(rdb)

}

func ping(client *redis.Client) error {
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return err

	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}
