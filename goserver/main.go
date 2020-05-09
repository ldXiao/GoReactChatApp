package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ldXiao/GoReactChatApp/router"
	"github.com/ldXiao/GoReactChatApp/websocket"
	"gopkg.in/olahol/melody.v1"
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

func setupWSRoutes(r *gin.Engine) {
	// pool := websocket.NewPool()
	// go pool.Start()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

}
func main() {
	r := router.Router()
	setupWSRoutes(r)
	r.Run(":5000")
}
