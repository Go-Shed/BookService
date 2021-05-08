package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/services"
)

type userHandler struct {
	UserService services.UserService
}

func NewUserHandler() userHandler {
	return userHandler{UserService: services.NewUserService()}
}

func (u *userHandler) FollowUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FollowUserRequest
	json.Unmarshal(reqBody, &request)

	err := u.UserService.FollowUser(uid, request.UserId)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200})
}

func (u *userHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FollowUserRequest
	json.Unmarshal(reqBody, &request)

	err := u.UserService.UnfollowUser(uid, request.UserId)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Error: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200})
}