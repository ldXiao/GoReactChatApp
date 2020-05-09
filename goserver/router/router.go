package router

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ldXiao/GoReactChatApp/middleware"
	"github.com/ldXiao/GoReactChatApp/models"
	"go.mongodb.org/mongo-driver/bson"
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
		doc := []models.Chat_info{}
		var err error = nil
		for cur.Next(context.Background()) {
			var chat models.Chat
			err = cur.Decode(&chat)
			doc = append(doc, chat.GetChatInfo())
		}
		// doc = append(doc, sample)

		if err == nil {
			if len(doc) == 0 {
				c.JSON(200, doc)
			} else {
				c.JSON(200, doc)
			}
		}
		log.Println("Fetched ", len(doc), "chat history")
	})

	r.POST("/api/chat/uploadfiles", func(c *gin.Context) {
		mf, fh, err := c.Request.FormFile("file")
		if err != nil {
			fmt.Println(err)
		}
		defer mf.Close()
		// create sha for file name
		splits := strings.Split(fh.Filename, ".")
		ext := splits[len(splits)-1]
		h := sha1.New()
		io.Copy(h, mf)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		// create new file
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "../uploads", fname)
		pathrel := filepath.Join("uploads", fname)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()
		// copy
		mf.Seek(0, 0)
		io.Copy(nf, mf)

		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"err":     true,
			})
		} else {
			c.JSON(200, gin.H{
				"success": true,
				"url":     pathrel,
			})
		}
	})
	r.Static("/uploads", "../uploads")

	return r
}
