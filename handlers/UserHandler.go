package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200})
}

func (u *userHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.UpdateUser
	json.Unmarshal(reqBody, &request)

	showMessage, err := u.UserService.UpdateUserName(uid, request.UserName)

	if err != nil {
		log.Print(err)
		message := "Some error occurred, please try again"
		if showMessage {
			message = err.Error()
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{HTTPCode: 500, Message: message})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200, Message: "Cool name!😎"})
}

func (u *userHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FollowUserRequest
	json.Unmarshal(reqBody, &request)

	err := u.UserService.UnfollowUser(uid, request.UserId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200})
}

func (u *userHandler) SearchUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.SearchUserRequest
	json.Unmarshal(reqBody, &request)

	results := u.UserService.SearchUser(uid, request.Search)

	if len(results.SearchResults) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{HTTPCode: 400, Message: "No users found"})
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (u *userHandler) FetchProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FetchProfileRequest
	json.Unmarshal(reqBody, &request)

	response, err := u.UserService.FetchProfile(uid, request.ProfileUserId, request.IsSelf)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Something went wrong. Please try again!"})
		return
	}

	if len(response.UserName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "You are lost, this profile does not exist!"})
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (u *userHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FetchProfileRequest
	json.Unmarshal(reqBody, &request)

	response, err := u.UserService.GetFollowers(uid, request.ProfileUserId, request.IsSelf)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Something went wrong. Please try again!"})
		return
	}

	if len(response.Follows) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "No followers here!"})
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (u *userHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.FetchProfileRequest
	json.Unmarshal(reqBody, &request)

	response, err := u.UserService.GetFollowing(uid, request.ProfileUserId, request.IsSelf)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "Something went wrong. Please try again!"})
		return
	}

	if len(response.Follows) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 500, Message: "No one here!"})
		return
	}

	json.NewEncoder(w).Encode(response)
}
