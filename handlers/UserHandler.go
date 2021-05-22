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

func (u *userHandler) SearchUser(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.SearchUserRequest
	json.Unmarshal(reqBody, &request)

	user := u.UserService.SearchUser(request.Search)

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200, Data: user})
}

func (u *userHandler) FetchProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FetchProfileRequest
	json.Unmarshal(reqBody, &request)

	response, err := u.UserService.FetchProfile(uid, request.ProfileUserId, request.IsSelf)

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 500, Message: "Something went wrong. Please try again!"})
	}

	json.NewEncoder(w).Encode(response)
}
