package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
	"shed/bookservice/common/constants"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/services"
)

type postHandler struct {
	PostsService services.PostsService
}

func NewPostHandler() postHandler {
	return postHandler{PostsService: services.NewPostsService()}
}

func (p *postHandler) GetPosts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.GetPostsRequest
	json.Unmarshal(reqBody, &request)

	response, err := p.PostsService.GetPosts(uid, request.ScreenName, request.ForUserID, request.IsSelf)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Looks like there is some issue here, please try again after some time"})
		return
	}

	if len(response.Posts) == 0 {
		msg := getMessageToShownOnScreen(request.ScreenName)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 400, Message: msg})
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (p *postHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.AddPostRequest
	json.Unmarshal(reqBody, &request)

	if len(request.Post) == 0 || len(uid) == 0 || (len(request.BookId) == 0 && len(request.BookTitle) == 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 400, Message: "post or book empty"})
		return
	}

	err := p.PostsService.AddPost(request.Post, uid, request.BookId, request.BookTitle)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200, Message: "Post created successfully"})
}

func (p *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.UpdatePostRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.UpdatePost(request.PostId, request.Post, request.BookTitle, request.BookId, uid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200, Message: "Post updated Successfully"})
}

func (p *postHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.LikePost(request.PostId, uid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200})
}

func (p *postHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.UnlikePost(request.PostId, uid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200}) ////write response to http writer
}

func (p *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.DeletePost(request.PostId, uid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200, Message: "Post deleted successfully"}) ////write response to http writer
}

func getMessageToShownOnScreen(screen string) string {
	var msg string

	if screen == constants.SCREEN_HOME {
		msg = "No Posts here, why not follow someone!"
	} else if screen == constants.SCREEN_EXPLORE {
		msg = "Couldn't find posts for you, try again later!"
	} else {
		msg = "Ah, looks like someone's been unsocial lately. Post about your favourite books and find others like you ! :) "
	}
	return msg
}
