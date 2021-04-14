package main

import (
	"fmt"
	"log"
	"net/http"
	auth "shed/bookservice/handlers/Auth"
	"time"

	"github.com/gorilla/mux"
)

/*

Mux mathches routes in the order they were defined
*/
func handleRequests() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/signup", auth.AuthorizeUsers(http.HandlerFunc(auth.Signup))).Methods("POST")
	return router
}

func main() {

	router := handleRequests()

	fmt.Println("Starting server")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
