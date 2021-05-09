package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shed/bookservice/api"
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

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.GetBooksRequest
	json.Unmarshal(reqBody, &request)

	response, _ := handler.bookService.GetBooks(request.UserId)

	json.NewEncoder(w).Encode(response) ////write response to http writer
}
