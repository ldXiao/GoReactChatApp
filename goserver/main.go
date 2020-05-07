package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ldXiao/GoReactChatApp/models"
)

func main() {
	fmt.Println("GoReactChatApp server started")
	u := new(models.User)
	u.Name = "haha"
	u.Password = "dark secret"
	fmt.Println(u.Save())
	u.ComparePassword("dark secret")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
