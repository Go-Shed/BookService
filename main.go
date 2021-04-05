package main

import (
	"log"
	"net/http"
	"shed/bookservice/handlers"
	"time"

	"github.com/gorilla/mux"
)

/*

Mux mathches routes in the order they were defined
*/
func handleRequests() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/facts", handlers.GetFacts).Methods("POST")
	return router
}

/**
1. Install go, SHAME on you though 😛, if you have not done it already
2. Run service using [go run main.go] or create executable using [go build main.go]
4. CURL : curl -X POST --data "{\"user_id\": \"afds\"}" localhost:8000/facts
4. Add your routes in handleRequests func and add corresponding handler to handlers folder.
5. Handler should be a function: f func(http.ResponseWriter,
	*http.Request)
6. No need to use servies, only use services when there is too much business logic, otherwise just make call to repos.
7. Db and external calls strictly in repos folder.
8. Feel free to create folder if required.
9. PS: How go packages work: every folder is a package. Aceess method and object using these packages
10. What is go.mod here? Used to manage depndency.
**/
func main() {

	// db := dgraph.Dgraph{
	// 	URL:    "https://billowing-night.ap-south-1.aws.cloud.dgraph.io/graphql",
	// 	Secret: "ZTE4YjRhNGEwYTAxNWRiYjMwMGI4YmEyMjc1ODU5ZmI="}

	// fmt.Printf("%+v", db.GetUsers("lol"))

	router := handleRequests()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
