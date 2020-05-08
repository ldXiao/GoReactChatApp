package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/ldXiao/GoReactChatApp/middleware"
	"github.com/ldXiao/GoReactChatApp/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func hah(w http.ResponseWriter, r *http.Request) {
	_, err := httputil.DumpRequest(r, true)

	if err != nil {
		fmt.Println("called0")
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Println("called1")
		fmt.Println(r.Cookie("w_auth"))
	}
}

type Test_struct struct {
}

func main() {
	r := gin.Default()
	r.GET("/api/users/auth", func(c *gin.Context) {
		var u models.User
		// fmt.Println(c.Request.Header)
		// auth request gives an empty body
		tok, err := c.Cookie("w_auth")
		var succ bool = false
		if err == nil {
			succ = u.LoadByToken(tok)
		}
		fmt.Printf(tok)

		c.JSON(200, gin.H{
			"isAuth": succ,
			"error":  !succ,
		})
	})
	r.POST("/api/users/login", func(c *gin.Context) {

		// fmt.Println(c.Request.Header)
		var logi models.Login
		fmt.Println(c.Cookie("w_auth"))
		err := c.Bind(&logi)
		if err != nil {
			log.Fatal(err)
		}
		var user models.User
		middleware.UsersCollection.FindOne(context.TODO(), bson.D{{"email", logi.Email}}).Decode(&user)
		if user.Email == "" {
			fmt.Println("email not found")
			c.JSON(404, gin.H{
				"loginSuccess": false,
				"message":      "Auth failed, email not found",
			})
			return
		}
		if !user.ComparePassword(logi.Password) {
			c.JSON(404, gin.H{
				"loginSuccess": false,
				"message":      "Wrong password",
			})
			return
		}
		// when login session expires, need to generate a new token
		user.GenerateToken()
		// and save it into database again
		err = user.UpdateToken()
		if err == nil {
			c.SetCookie("w_auth", user.Token, 0, "/", "localhost", false, true)
			c.SetCookie("w_authExp", user.TokenExp, 0, "/", "localhost", false, true)
			c.JSON(200, gin.H{
				"loginSuccess": true,
			})
		}
	})
	r.POST("/api/users/register", func(c *gin.Context) {

		var user models.User
		fmt.Println("register called")
		err := c.Bind(&user)
		fmt.Println("register called")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("register called")
		b, err := json.Marshal(user)
		fmt.Println("register called")
		if err != nil {
			fmt.Println(err)
			return
		}
		if user.Save() {
			c.JSON(200, gin.H{
				"success": true,
			})
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"err":     "Invalid registrtion, email exists",
			})
		}
		fmt.Println("register called")
		fmt.Println(string(b))
	})

	r.Run(":5000")
}
