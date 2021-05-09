package services

import (
	"errors"
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/common/constants"
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

func (p *PostsService) GetPosts(userId, screenName string, isSelf bool) ([]api.GetPostsResponse, error) {
	if screenName == constants.SCREEN_HOME {
		return p.getHomeScreen(userId)
	} else if screenName == constants.SCREEN_PROFILE {
		return p.getProfileScreen(userId, isSelf)
	} else {
		return p.getExploreScreen(userId)
	}
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

	if err != nil || response.Id == "" {
		return errors.New("post not created")
	}
	return nil
}

func (p *PostsService) LikePost(postId, userId string) error {
	client := p.PostRepo

	err := client.LikePost(postId, userId)

	if err != nil {
		return err
	}
	return nil
}

func (p *PostsService) UnlikePost(postId, userId string) error {
	client := p.PostRepo

	err := client.UnlikePost(postId, userId)

	if err != nil {
		return err
	}
	return nil
}

func (p *PostsService) getHomeScreen(userId string) ([]api.GetPostsResponse, error) {

	client := p.PostRepo
	user, err := client.GetUserFeedHomeScreen(userId)

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

func (p *PostsService) getProfileScreen(userId string, isSelf bool) ([]api.GetPostsResponse, error) {
	client := p.PostRepo
	user, err := client.GetUserHomeProfileScreen(userId)

	if err != nil {
		return []api.GetPostsResponse{}, err
	}

	var response []api.GetPostsResponse

	var showEditButton bool
	if isSelf {
		showEditButton = true
	}
	userFeedItem := api.GetPostsResponse{
		UserId:        userId,
		UserName:      user.Username,
		IsFollowed:    true,
		ShowFollowBtn: false,
		ShowEditBtn:   showEditButton,
		UserPhoto:     user.Username,
		IsLiked:       false,
	}

	for _, post := range user.Posts {

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

	////// sort according to date
	sort.Slice(response, func(i, j int) bool {
		return response[i].CreatedAt > response[j].CreatedAt
	})

	return response, nil
}

func (p *PostsService) getExploreScreen(userId string) ([]api.GetPostsResponse, error) {

	client := p.PostRepo
	posts, err := client.GetExporeScreen(userId)

	if err != nil {
		return []api.GetPostsResponse{}, err
	}

	var response []api.GetPostsResponse

	for _, post := range posts {

		postItem := api.GetPostsResponse{
			UserId:          post.Author.UserId,
			UserName:        post.Author.Username,
			PostColor:       post.Color,
			ShowEditBtn:     false,
			ShowFollowBtn:   true,
			ShowLikeSection: true,
			IsFollowed:      false,
			IsLiked:         false,
			UserPhoto:       post.Author.UserPhoto,
			Text:            post.Text,
			PostId:          post.Id,
			LikeCount:       fmt.Sprint(post.LikesAggregate.Count),
			CreatedAt:       fmt.Sprint(post.CreatedAt),
		}

		if len(post.Likes) > 0 {
			postItem.IsLiked = true
		}
		if len(post.Author.Followers) > 0 {
			postItem.IsFollowed = true
		}
		response = append(response, postItem)
	}

	////// sort according to likes
	sort.Slice(response, func(i, j int) bool {
		return response[i].LikeCount > response[j].LikeCount
	})

	return response, nil
}
