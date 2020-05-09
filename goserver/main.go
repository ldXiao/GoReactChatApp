package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ldXiao/GoReactChatApp/middleware"
	"github.com/ldXiao/GoReactChatApp/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

type authmodel struct {
	_id      string
	isAdmin  bool
	isAuth   bool
	email    string
	name     string
	lastname string
	role     int
	image    string
}

func main() {
	r := gin.Default()
	r.GET("/api/users/auth", func(c *gin.Context) {
		var u models.User
		// fmt.Println(c.Request.Header)
		// auth request gives an empty body
		// requestDump, err := httputil.DumpRequest(c.Request, true)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("body", string(requestDump))
		tok, err := c.Cookie("w_auth")
		var succ bool = false
		if err == nil {
			succ = u.LoadByToken(tok) // load every filed in to user and return succees
		}
		fmt.Printf(tok)
		if succ {
			c.JSON(200, gin.H{
				"isAuth":   succ,
				"error":    !succ,
				"_id":      u.ID.Hex(),
				"isAdmin":  u.Role == 1,
				"email":    u.Email,
				"name":     u.Name,
				"lastname": u.LastName,
				"role":     u.Role,
				"image":    u.Image,
			})
		} else {
			c.JSON(200, gin.H{
				"isAuth": succ,
				"error":  !succ,
			})
		}
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
		err := c.Bind(&user)
		if err != nil {
			log.Fatal(err)
		}
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
	})

	r.GET("/api/users/logout", func(c *gin.Context) {
		var u models.User
		tok, err := c.Cookie("w_auth")
		var succ bool = false
		if err == nil {
			succ = u.LoadByToken(tok) // load every filed in to user and return succees
		}
		u.Token = ""
		u.TokenExp = ""

		u.UpdateToken()
		c.JSON(200, gin.H{
			"success": succ,
		})
	})

	r.Run(":5000")
}
