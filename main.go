package main

import (
	"net/http"
	"shed/bookservice/handlers"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/services"

	"github.com/gorilla/mux"
)

/*

Mux mathches routes in the order they were defined
*/
func handleRequests() *mux.Router {

	postsHandler := handlers.NewPostHandler()
	userHandler := handlers.NewUserHandler()
	bookHandler := handlers.NewBookHandler()

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/signup", auth.AuthorizeUsers(http.HandlerFunc(auth.Signup))).Methods("POST")
	router.HandleFunc("/getPosts", postsHandler.GetPosts).Methods("POST")
	router.Handle("/addPost", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.AddPost))).Methods("POST")
	router.Handle("/updatePost", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.AddPost))).Methods("POST")
	router.Handle("/follow", auth.AuthorizeUsers(http.HandlerFunc(userHandler.FollowUser))).Methods("POST")
	router.Handle("/unfollow", auth.AuthorizeUsers(http.HandlerFunc(userHandler.UnfollowUser))).Methods("POST")
	router.Handle("/like", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.LikePost))).Methods("POST")
	router.Handle("/unlike", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.UnlikePost))).Methods("POST")
	router.HandleFunc("/searchUser", userHandler.SearchUser).Methods("POST")
	router.HandleFunc("/getBooks", bookHandler.GetBooks).Methods("POST")
	router.HandleFunc("/deletePost", postsHandler.DeletePost).Methods("POST")

	return router
}

func main() {

	// router := handleRequests()

	// fmt.Println("Starting server")
	// srv := &http.Server{
	// 	Handler:      router,
	// 	Addr:         "localhost:8000",
	// 	WriteTimeout: 5 * time.Second,
	// 	ReadTimeout:  5 * time.Second,
	// }

	// log.Fatal(srv.ListenAndServe())

	srv := services.NewPostsService()
	srv.AddPost("Yes I am working", "efgh", "", "Six of crows")
}
