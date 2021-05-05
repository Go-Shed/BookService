package services

import (
	"errors"
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
	"sort"
	"time"
)

type PostsService struct {
	PostRepo query.PostRepo
}

func NewPostsService() PostsService {
	return PostsService{PostRepo: query.NewPostRepo()}
}

func (p *PostsService) GetPosts(userId, screenName string) ([]api.GetPostsResponse, error) {
	client := p.PostRepo

	user, err := client.GetUserFeedHomeScreen(userId)

	if err != nil {
		panic(err)
	}

	if err != nil {
		return []api.GetPostsResponse{}, err
	}

	var response []api.GetPostsResponse

	for _, following := range user.Following {

		userFeedItem := api.GetPostsResponse{
			UserId:        following.UserId,
			UserName:      following.Username,
			IsFollowed:    true,
			ShowFollowBtn: true,
			ShowEditBtn:   false,
			UserPhoto:     following.UserPhoto,
			IsLiked:       false,
		}

		for _, post := range following.Posts {
			item := userFeedItem
			item.Text = post.Text
			item.PostColor = post.Color
			item.LikeCount = fmt.Sprint(post.LikesAggregate.Count)
			item.PostId = fmt.Sprint(post.Id)
			item.CreatedAt = fmt.Sprint(post.CreatedAt)

			if len(post.Likes) > 0 {
				item.IsLiked = true
			}
			response = append(response, item)
		}
	}

	////// sort according to date
	sort.Slice(response, func(i, j int) bool {
		return response[i].CreatedAt > response[j].CreatedAt
	})
	return response, nil
}

func (p *PostsService) AddPost(text, color, userId string) error {
	client := p.PostRepo
	timeNow := time.Now().Local().String()
	post := model.Post{Text: text, Color: color, Author: model.User{UserId: userId},
		CreatedAt: timeNow, UpdatedAt: timeNow}

	response, err := client.CreatePost(post)

	fmt.Println(err)
	if err != nil || response.Id == "" {
		return errors.New("post not created")
	}
	return nil
}

func (p *PostsService) UpdatePost(text, userId, postId string) error {
	client := p.PostRepo
	timeNow := time.Now().Local().String()
	post := model.Post{Id: postId, Text: text, UpdatedAt: timeNow}

	response, err := client.UpdatePost(post)

	fmt.Println(err)
	if err != nil || response.Id == "" {
		return errors.New("post not created")
	}
	return nil
}
