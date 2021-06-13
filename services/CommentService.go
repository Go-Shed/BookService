package services

import (
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
	"strings"
	"time"
)

type CommentService struct {
	CommentRepo query.CommentRepo
}

func NewCommentService() CommentService {
	return CommentService{CommentRepo: query.NewCommentRepo()}
}

func (p *CommentService) AddComment(text, userId, postId string) error {

	text = strings.TrimSpace(text)

	if len(text) == 0 || len(userId) == 0 || len(postId) == 0 {
		return fmt.Errorf("comment can not be empty")
	}

	timeNow := time.Now().Local().String()
	comment := model.Comment{
		Text:      text,
		User:      model.User{UserId: userId},
		Post:      model.Post{Id: postId},
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	err := p.CommentRepo.AddComment(comment)
	return err
}

func (p *CommentService) GetComments(postId, userId string) (api.GetCommentsResponse, error) {

	postId = strings.TrimSpace(postId)

	if len(postId) == 0 {
		return api.GetCommentsResponse{}, fmt.Errorf("post id can not be empty")
	}

	comments, err := p.CommentRepo.GetComments(postId)

	if err != nil {
		return api.GetCommentsResponse{}, err
	}

	var otherComments []api.CommentItem
	var myComments []api.CommentItem

	for _, comment := range comments {

		item := api.CommentItem{
			Text:      comment.Text,
			UserName:  comment.User.Username,
			UserId:    comment.User.UserId,
			UserPhoto: comment.User.UserPhoto,
			CreatedAt: comment.CreatedAt,
			CommentId: comment.Id,
		}

		if item.UserId == userId {
			myComments = append(myComments, item)
			continue
		}

		otherComments = append(otherComments, item)
	}

	response := append(myComments, otherComments...)

	return api.GetCommentsResponse{Comments: response}, nil
}
