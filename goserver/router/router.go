package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllTask(w http.ResponseWriter, r *http.Request) {
	bs := string(r.Body)
	fmt.Println(bs)
}

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/login").Methods("GET", "OPTIONS")
	return router
}
