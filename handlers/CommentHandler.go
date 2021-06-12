package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"shed/bookservice/api"
	auth "shed/bookservice/handlers/Auth"
	"shed/bookservice/services"
)

type commentHandler struct {
	CommentService services.CommentService
}

func NewCommentHandler() commentHandler {
	return commentHandler{CommentService: services.NewCommentService()}
}

func (p *commentHandler) AddComment(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	uid := ctx.Value(auth.RequestContextKey{Key: "uid"}).(string)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var request api.AddCommentRequest
	json.Unmarshal(reqBody, &request)

	err := p.CommentService.AddComment(request.Text, uid, request.PostId)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{HTTPCode: 500, Message: "Some error occurred, please try again"})
		return
	}

	json.NewEncoder(w).Encode(api.ApiResponse{HTTPCode: 200, Message: "Done😃"})
}
