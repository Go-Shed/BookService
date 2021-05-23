package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/services"
)

type bookHandler struct {
	bookService services.BookService
}

func NewBookHandler() bookHandler {
	return bookHandler{
		bookService: services.NewBookService(),
	}
}

func (handler *bookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	if len(uid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ApiResponse{ResponseCode: 400, Message: "Sign in to explore world around books!"})
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.GetBooksRequest
	json.Unmarshal(reqBody, &request)

	response, _ := handler.bookService.GetBooks(uid, request.ProfileUserId, request.IsSelf)

	json.NewEncoder(w).Encode(response)
}
