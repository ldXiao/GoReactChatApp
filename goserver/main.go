package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/ldXiao/GoReactChatApp/models"
)

func hah(w http.ResponseWriter, r *http.Request) {
	_, err := httputil.DumpRequest(r, true)

	if err != nil {
		fmt.Println("called")
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Println("called1")
		fmt.Println(r.Cookie("w_auth"))
	}
}

// func main() {
// 	fmt.Println("GoReactChatApp server started")
// 	router := mux.NewRouter()
// 	// u := new(models.User)
// 	// u.Name = "haha"
// 	// u.Password = "dark secret"
// 	// fmt.Println(u.Save())
// 	// u.ComparePassword("dark secret")
// 	router.HandleFunc("/api/users/auth", hah)
// 	log.Fatal(http.ListenAndServe(":5000", router))
// }

type Test_struct struct {
}

func main() {
	r := gin.Default()
	r.GET("/api/users/auth", func(c *gin.Context) {
		fmt.Println("called")
		var u models.User
		// fmt.Println(c.Request.Header)

		fmt.Println(c.Cookie("w_auth"))
		c.JSON(200, gin.H{
			"message": "pong",
		})
		err0 := c.BindJSON(&u)
		fmt.Println(err0)
	})
	r.POST("/api/users/login", func(c *gin.Context) {
		fmt.Println("called")

		// fmt.Println(c.Request.Header)

		bodybyte, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println(string(bodybyte))
		fmt.Println(c.Cookie("w_auth"))
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":5000")
}
