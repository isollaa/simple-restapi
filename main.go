package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/isollaa/simple-restfull-api/config"
	"github.com/isollaa/simple-restfull-api/handler"
)

const strToken = "mysecret"

func main() {
	if onRequestToken() {
		return
	}

	session := config.DBInit()
	db := &handler.DB{DB: session}

	router := gin.Default()

	// router.POST("/login", loginHandler)
	router.GET("/client/:id", auth, db.GetClient)
	router.GET("/clients", auth, db.GetClients)
	router.POST("/client", auth, db.CreateClient)
	router.PUT("/client", auth, db.UpdateClient)
	router.DELETE("/client/:id", auth, db.DeleteClient)
	router.Run(":3000")
}

func onRequestToken() bool {
	g := flag.Bool("token", false, "generate token")
	flag.Parse()
	if *g {
		err := generateToken(strToken)
		if err != nil {
			log.Print(err)
		}
	}
	return *g
}

func generateToken(str string) error {
	fmt.Print()
	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	token, err := sign.SignedString([]byte(str))
	if err != nil {
		return err
	}
	log.Printf("token bearer : %s", token)
	return nil
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

// type Credential struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// func loginHandler(c *gin.Context) {
// 	var user Credential
// 	err := c.Bind(&user)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  http.StatusBadRequest,
// 			"message": "can't bind struct",
// 		})
// 	}
// 	if user.Username != "myname" {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"status":  http.StatusUnauthorized,
// 			"message": "wrong username or password",
// 		})
// 	} else {
// 		if user.Password != "myname123" {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"status":  http.StatusUnauthorized,
// 				"message": "wrong username or password",
// 			})
// 		}
// 	}
// 	sign := jwt.New(jwt.GetSigningMethod("HS256"))
// 	token, err := sign.SignedString([]byte("secret"))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 		c.Abort()
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"token": token,
// 	})
// }

// func auth(c *gin.Context) {
// 	tokenString := c.Request.Header.Get("Authorization")
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if jwt.GetSigningMethod("HS256") != token.Method {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte("secret"), nil
// 	})

// 	// if token.Valid && err == nil {
// 	if token != nil && err == nil {
// 		fmt.Println("token verified")
// 	} else {
// 		result := gin.H{
// 			"message": "not authorized",
// 			"error":   err.Error(),
// 		}
// 		c.JSON(http.StatusUnauthorized, result)
// 		c.Abort()
// 	}
// }
