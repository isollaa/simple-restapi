package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/isollaa/simple-restapi/config"
	"github.com/isollaa/simple-restapi/handler"
	"github.com/jinzhu/gorm"
)

const strToken = "mysecret"

func main() {
	// if onRequestToken() {
	// 	return
	// }

	session := config.DBInit()
	db := &handler.DB{DB: session}
	router := gin.Default()

	router.POST("/login", login)
	router.GET("/user/:id", auth, db.GetUser)
	router.GET("/users", auth, db.GetUsers)
	router.POST("/user", auth, db.CreateUser)
	router.PUT("/user/:id", auth, db.UpdateUser)
	router.DELETE("/user/:id", auth, db.DeleteUser)

	router.Run(":3000")
}

// func onRequestToken() bool {
// 	g := flag.Bool("token", false, "generate token")
// 	flag.Parse()
// 	if *g {
// 		err := generateToken(strToken)
// 		if err != nil {
// 			log.Print(err)
// 		}
// 	}
// 	return *g
// }

func generateToken(str string) (string, error) {
	fmt.Print()
	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	token, err := sign.SignedString([]byte(str))
	if err != nil {
		return "", err
	}
	return token, err
}

func login(c *gin.Context) {
	var db *gorm.DB = config.DBInit()
	client := handler.Client{
		Username: c.PostForm("username"),
		Password: c.PostForm("password")}
	result := gin.H{}

	if err := db.Table(handler.TABLENAME).Where("username = ? AND password = ?", client.Username, client.Password).First(&client).Error; err != nil {
		result = gin.H{"result": "username / password salah"}
	} else {
		token, err := generateToken(strToken)
		if err != nil {
			log.Print(err)
		}
		result = gin.H{"token": token}
	}
	c.JSON(http.StatusOK, result)
}

func auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	bearerToken := strings.Split(tokenString, " ")
	if len(bearerToken) == 2 {
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(strToken), nil
		})

		if token != nil && err == nil {
			fmt.Println("token verified")
		} else {
			result := gin.H{
				"message": "not authorized",
				"error":   err.Error(),
			}
			c.JSON(http.StatusUnauthorized, result)
			c.Abort()
		}
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   "An authorization header is required",
		}

		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}

}
