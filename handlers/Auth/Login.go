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
	UserId            string `json:"user_id"`
	AlreadyRegistered bool   `json:"alreadyRegistered"`
}

func Signup(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(RequestContextKey{Key: "uid"}).(string)
	reqBody, _ := ioutil.ReadAll(r.Body)
	var request AddUserRequest
	json.Unmarshal(reqBody, &request)

	if len(request.Email) == 0 || len(request.Username) == 0 || len(uid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{HTTPCode: 400, Message: "Username, email or userId can not be empty"})
		return
	}

	userRepo := query.NewUserRepo()
	_, err := userRepo.CreateUser(model.User{Username: request.Username, UserId: uid, Email: request.Email})

	if err != nil {
		json.NewEncoder(w).Encode(AddUserResponse{
			UserId:            uid,
			AlreadyRegistered: true,
		})
		return
	}

	json.NewEncoder(w).Encode(AddUserResponse{
		UserId:            uid,
		AlreadyRegistered: false,
	})
}
