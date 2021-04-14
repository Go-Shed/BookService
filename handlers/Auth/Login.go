package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
)

type AddUserRequest struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
}

type AddUserResponse struct {
	Username string `json:"user_name"`
}

func Signup(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(RequestContextKey{Key: "uid"}).(string)
	reqBody, _ := ioutil.ReadAll(r.Body)
	var request AddUserRequest
	json.Unmarshal(reqBody, &request)

	userRepo := query.NewUserRepo()
	_, err := userRepo.CreateUser(model.User{Username: request.Username, UserId: uid, Email: request.Email})

	if err != nil {
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 400, Error: "username already exists"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 200, Error: "", Data: request.Username})
}
