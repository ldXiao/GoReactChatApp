package router

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ldXiao/GoReactChatApp/middleware"
	"github.com/ldXiao/GoReactChatApp/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

// Router is exported and used in main.go
func Router() *gin.Engine {

	r := gin.Default()
	r.GET("/api/users/auth", func(c *gin.Context) {
		var u models.User
		tok, err := c.Cookie("w_auth")
		var succ bool = false
		if err == nil {
			succ = u.LoadByToken(tok) // load every filed in to user and return succees
		}
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

		var logi models.Login
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

	r.GET("/api/chat/getChats", func(c *gin.Context) {

		cur, _ := middleware.ChatsCollection.Find(context.TODO(), bson.D{})
		// }
		var doc []models.Chat
		var sample models.Chat
		sample.ID, _ = primitive.ObjectIDFromHex("5eb6af835236b8ffbc985f62")
		sample.CreatedAt = primitive.Timestamp{1, 1}
		sample.UpdatedAt = primitive.Timestamp{1, 1}
		sample.Message = "test"
		sample.Sender, _ = primitive.ObjectIDFromHex("5eb6af835236b8ffbc985f61")
		sample.Type = "Text"
		err := cur.All(context.TODO(), &doc)
		doc = append(doc, sample)
		if err == nil {
			c.JSON(200, doc)
		}
		log.Println("Fetched ", len(doc), "chat history")
	})

	// r.GET("/", func(c *gin.Context) {
	// 	c.String(200, "We got Gin")
	// })

	// r.GET("/ws", func(c *gin.Context) {
	// 	fmt.Println("calledws")
	// 	wshandler(c.Writer, c.Request)
	// })
	return r
}
