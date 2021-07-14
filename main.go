package main

import (
	"net/http"
	"shed/bookservice/handlers"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/scripts"

	"github.com/gorilla/mux"
)

/*
Mux mathches routes in the order they were defined
*/
func handleRequests() *mux.Router {

	postsHandler := handlers.NewPostHandler()
	userHandler := handlers.NewUserHandler()
	bookHandler := handlers.NewBookHandler()
	commentHandler := handlers.NewCommentHandler()

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/signup", auth.AuthorizeUsers(http.HandlerFunc(auth.Signup))).Methods("POST")
	router.Handle("/getPosts", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.GetPosts))).Methods("POST")
	router.Handle("/addPost", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.AddPost))).Methods("POST")
	router.Handle("/updatePost", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.UpdatePost))).Methods("POST")
	router.Handle("/updateUser", auth.AuthorizeUsers(http.HandlerFunc(userHandler.UpdateUserName))).Methods("POST")
	router.Handle("/addComment", auth.AuthorizeUsers(http.HandlerFunc(commentHandler.AddComment))).Methods("POST")

	router.Handle("/follow", auth.AuthorizeUsers(http.HandlerFunc(userHandler.FollowUser))).Methods("POST")
	router.Handle("/unfollow", auth.AuthorizeUsers(http.HandlerFunc(userHandler.UnfollowUser))).Methods("POST")
	router.Handle("/like", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.LikePost))).Methods("POST")
	router.Handle("/unlike", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.UnlikePost))).Methods("POST")
	router.Handle("/getLikes", auth.AuthorizeUsers(http.HandlerFunc(postsHandler.GetLikes))).Methods("POST")
	router.Handle("/getComments", auth.AuthorizeUsers(http.HandlerFunc(commentHandler.GetComments))).Methods("POST")

	router.Handle("/searchUser", auth.AuthorizeUsers(http.HandlerFunc(userHandler.SearchUser))).Methods("POST")
	router.Handle("/getBooks", auth.AuthorizeUsers(http.HandlerFunc(bookHandler.GetBooks))).Methods("POST")
	router.Handle("/fetchProfile", auth.AuthorizeUsers(http.HandlerFunc(userHandler.FetchProfile))).Methods("POST")
	router.HandleFunc("/deletePost", postsHandler.DeletePost).Methods("POST")

	router.Handle("/getFollowers", auth.AuthorizeUsers(http.HandlerFunc(userHandler.GetFollowers))).Methods("POST")
	router.Handle("/getFollowing", auth.AuthorizeUsers(http.HandlerFunc(userHandler.GetFollowing))).Methods("POST")

	return router
}

func main() {

	// router := handleRequests()

	// fmt.Println("Starting server")
	// srv := &http.Server{
	// 	Handler:      router,
	// 	Addr:         ":8080",
	// 	WriteTimeout: 5 * time.Second,
	// 	ReadTimeout:  5 * time.Second,
	// }

	// log.Fatal(srv.ListenAndServe())

	scripts.CleanData("/Users/dhairya/Desktop/ol_dump_works_2021-05-13.txt")
}
