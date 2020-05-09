package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ldXiao/GoReactChatApp/router"
	"github.com/ldXiao/GoReactChatApp/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes(r *gin.Engine) {
	pool := websocket.NewPool()
	go pool.Start()

	r.GET("/", func(c *gin.Context) {
		serveWs(pool, c.Writer, c.Request)
	})
}
func main() {
	r := router.Router()
	setupRoutes(r)
	r.Run(":5000")
}
