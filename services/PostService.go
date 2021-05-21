package services

import (
	"errors"
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/common/constants"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
	"sort"
	"strings"
	"time"
)

type PostsService struct {
	PostRepo query.PostRepo
	BookRepo query.BookRepo
}

func NewPostsService() PostsService {
	return PostsService{PostRepo: query.NewPostRepo(), BookRepo: query.NewBookRepo()}
}

func (p *PostsService) GetPosts(userId, screenName, forUserId string, isSelf bool) (api.GetPostsResponse, error) {
	if screenName == constants.SCREEN_HOME {
		return p.getHomeScreen(userId)
	} else if screenName == constants.SCREEN_PROFILE {
		return p.getProfileScreen(userId, forUserId, isSelf)
	} else {
		return p.getExploreScreen(userId)
	}
}

func (p *PostsService) AddPost(text, userId, bookId, bookTitle string) error {

	client := p.PostRepo
	timeNow := time.Now().Local().String()

	bookTitle = strings.TrimSpace(strings.ToLower(bookTitle))
	bookTitle = strings.Join(strings.Fields(bookTitle), " ")

	newBook, err := p.BookRepo.CreateOrGetBook(bookId, bookTitle)

	if err != nil {
		return err
	}

	book := model.Book{Id: newBook}

	post := model.Post{Text: text, CreatedAt: timeNow, UpdatedAt: timeNow, Book: book}
	user := model.User{UserId: userId, Books: []model.Book{book}, Posts: []model.Post{post}}

	err = client.CreatePost(user)

	if err != nil {
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

func (p *PostsService) DeletePost(postId string) error {
	client := p.PostRepo

	err := client.DeletePost(postId)

	if err != nil {
		return err
	}
	return nil
}

func (p *PostsService) getHomeScreen(userId string) (api.GetPostsResponse, error) {

	client := p.PostRepo
	user, err := client.GetUserFeedHomeScreen(userId)

	if err != nil {
		return api.GetPostsResponse{}, err
	}

	var response []api.GetPostResponse

	for _, following := range user.Following {

		userFeedItem := api.GetPostResponse{
			UserId:        following.UserId,
			UserName:      following.Username,
			IsFollowed:    true,
			ShowFollowBtn: false,
			ShowEditBtn:   false,
			UserPhoto:     following.UserPhoto,
			IsLiked:       false,
		}

		for _, post := range following.Posts {
			item := userFeedItem
			item.Text = post.Text
			item.LikeCount = fmt.Sprint(post.LikesAggregate.Count)
			item.PostId = fmt.Sprint(post.Id)
			item.CreatedAt = fmt.Sprint(post.CreatedAt)
			item.Book = api.GetBooksResponse{BookId: post.Book.Id, BookName: post.Book.Name}

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
	return api.GetPostsResponse{Posts: response}, nil
}

func (p *PostsService) getProfileScreen(userId, forUserId string, isSelf bool) (api.GetPostsResponse, error) {

	if (!isSelf && len(forUserId) == 0) || len(userId) == 0 {
		return api.GetPostsResponse{}, fmt.Errorf("UserId mising")
	}

	client := p.PostRepo

	if isSelf {
		forUserId = userId
	} else {
		userId = forUserId
		forUserId = userId
	}
	user, err := client.GetUserHomeProfileScreen(userId, forUserId)

	if err != nil {
		return api.GetPostsResponse{}, err
	}

	var response []api.GetPostResponse

	showEditButton, showFollowBtn := false, true
	if isSelf {
		showEditButton = true
		showFollowBtn = false
	}
	userFeedItem := api.GetPostResponse{
		UserId:        userId,
		UserName:      user.Username,
		IsFollowed:    len(user.Followers) > 0,
		ShowFollowBtn: showFollowBtn,
		ShowEditBtn:   showEditButton,
		UserPhoto:     user.Username,
		IsLiked:       false,
	}

	for _, post := range user.Posts {

		item := userFeedItem
		item.Text = post.Text
		item.LikeCount = fmt.Sprint(post.LikesAggregate.Count)
		item.PostId = fmt.Sprint(post.Id)
		item.CreatedAt = fmt.Sprint(post.CreatedAt)
		item.Book = api.GetBooksResponse{BookId: post.Book.Id, BookName: post.Book.Name}

		if len(post.Likes) > 0 {
			item.IsLiked = true
		}
		response = append(response, item)
	}

	////// sort according to date
	sort.Slice(response, func(i, j int) bool {
		return response[i].CreatedAt > response[j].CreatedAt
	})

	return api.GetPostsResponse{Posts: response}, nil
}

func (p *PostsService) getExploreScreen(userId string) (api.GetPostsResponse, error) {

	client := p.PostRepo
	posts, err := client.GetExporeScreen(userId)

	if err != nil {
		return api.GetPostsResponse{}, err
	}

	var response []api.GetPostResponse

	for _, post := range posts {

		postItem := api.GetPostResponse{
			UserId:          post.Author.UserId,
			UserName:        post.Author.Username,
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
			Book:            api.GetBooksResponse{BookId: post.Book.Id, BookName: post.Book.Name},
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

	return api.GetPostsResponse{Posts: response}, nil
}
