package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
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

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.GetPostsRequest
	json.Unmarshal(reqBody, &request)

	response, err := p.PostsService.GetPosts(request.UserId, request.ScreenName, true)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 400, Error: "No posts to show, why not follow someone!"})
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

	err := p.PostsService.AddPost(request.Post, request.PostColor, uid)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200})
}

func (p *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.UpdatePostRequest
	json.Unmarshal(reqBody, &request)
	err := p.PostsService.UpdatePost(request.PostId, request.Text, uid)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200})
}

func (p *postHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.LikePost(request.PostId, uid)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200})
}

func (p *postHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.UnlikePost(request.PostId, uid)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200}) ////write response to http writer
}

func (p *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.PostIdRequest
	json.Unmarshal(reqBody, &request)

	err := p.PostsService.DeletePost(request.PostId)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200}) ////write response to http writer
}
