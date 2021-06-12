package services

import (
	"fmt"
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
