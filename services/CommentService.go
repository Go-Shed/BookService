package services

import (
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/common"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
	"shed/bookservice/repos/notification"
	"strings"
	"time"
)

type CommentService struct {
	CommentRepo      query.CommentRepo
	UserRepo         query.UserRepo
	NotificationRepo notification.NotificationRepo
}

func NewCommentService() CommentService {
	return CommentService{CommentRepo: query.NewCommentRepo(), UserRepo: query.NewUserRepo(), NotificationRepo: notification.NewNotificationRepo()}
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

	notification, err := p.CommentRepo.AddComment(comment)

	if err != nil {
		return err
	}

	if notification.UserToSend != notification.UserBy {
		p.NotificationRepo.AddNotificationTODB(notification)
	}
	return nil
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

	var response []api.CommentItem

	for _, comment := range comments {

		item := api.CommentItem{
			Text:      comment.Text,
			UserName:  comment.User.Username,
			UserId:    comment.User.UserId,
			UserPhoto: comment.User.UserPhoto,
			CreatedAt: common.GetFormattedDate(comment.CreatedAt),
			CommentId: comment.Id,
		}

		response = append(response, item)
	}

	return api.GetCommentsResponse{Comments: response}, nil
}
