package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
	"shed/bookservice/services"
)

type postHandler struct {
	PostsService services.PostsService
}

func NewPostHandler() postHandler {
	return postHandler{PostsService: services.NewPostsService()}
}

func (p *postHandler) GetPosts(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body) /// Deserialize request
	var request api.GetPostsRequest
	json.Unmarshal(reqBody, &request) ///// deserialize and map it to object

	response, err := p.PostsService.GetPosts(request.UserId, request.ScreenName)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 400, Error: "No posts to show, why not follow someone!"})
		return
	}

	json.NewEncoder(w).Encode(response) ////write response to http writer
}
