package services

import (
	"errors"
	"fmt"
	"log"
	"shed/bookservice/api"
	"shed/bookservice/common"
	"shed/bookservice/common/constants"
	"shed/bookservice/repos/dgraph/model"
	"shed/bookservice/repos/dgraph/query"
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

	bookTitle = strings.TrimSpace(bookTitle)
	bookTitle = strings.Join(strings.Fields(bookTitle), " ")
	text = strings.Replace(text, "\n", "\\n", -1)

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

func (p *PostsService) UpdatePost(postId, text, bookTitle, bookId, userId string) error {

	client := p.PostRepo
	timeNow := time.Now().Local().String()
	text = strings.Replace(text, "\n", "\\n", -1)

	log.Print(text, postId)

	post := model.Post{Id: postId, Text: text, UpdatedAt: timeNow}
	user, err := client.GetPost(postId, userId)

	if len(user.Posts) == 0 || err != nil {
		return fmt.Errorf("post now owned by user")
	}

	response, err := client.UpdatePost(post)

	if err != nil || response.Id == "" {
		return errors.New("post not updated")
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

func (p *PostsService) GetLikes(postId, uid string) api.GetLikesResponse {
	client := p.PostRepo

	post, err := client.GetLikesOnPost(postId, uid)

	if len(post.Likes) == 0 || err != nil {
		fmt.Println(err)
		return api.GetLikesResponse{Likes: []api.LikeItem{}, TotalLikes: 0}
	}

	var likesList []api.LikeItem

	for _, item := range post.Likes {
		likeItem := api.LikeItem{
			UserPhoto:     item.UserPhoto,
			UserName:      item.Username,
			UserId:        item.UserId,
			ShowFollowBtn: true,
			IsFollowed:    false,
		}

		if len(item.Followers) != 0 {
			likeItem.IsFollowed = true
		}

		likesList = append(likesList, likeItem)
	}

	return api.GetLikesResponse{Likes: likesList, TotalLikes: len(post.Likes)}
}

func (p *PostsService) DeletePost(postId, uid string) error {
	client := p.PostRepo

	user, err := client.GetPost(postId, uid)

	if len(user.Posts) == 0 || err != nil {
		return fmt.Errorf("post now owned by user")
	}

	err = client.DeletePost(postId)

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
			UserId:          following.UserId,
			UserName:        following.Username,
			IsFollowed:      true,
			ShowFollowBtn:   false,
			ShowEditBtn:     false,
			UserPhoto:       following.UserPhoto,
			IsLiked:         false,
			ShowLikeSection: true,
		}

		for _, post := range following.Posts {
			item := userFeedItem
			item.Text = strings.Replace(post.Text, "\\n", "\n", -1)
			item.LikeCount = fmt.Sprint(post.LikesAggregate.Count)
			item.PostId = fmt.Sprint(post.Id)
			item.CreatedAt = common.GetFormattedDate(post.CreatedAt)
			item.Book = api.BookResponse{BookId: post.Book.Id, BookName: toUpperCase(post.Book.Name)}

			if len(post.Likes) > 0 {
				item.IsLiked = true
			}
			response = append(response, item)
		}
	}

	return api.GetPostsResponse{Posts: response}, nil
}

func (p *PostsService) getProfileScreen(userId, forUserId string, isSelf bool) (api.GetPostsResponse, error) {

	if (!isSelf && len(forUserId) == 0) || len(userId) == 0 {
		return api.GetPostsResponse{}, fmt.Errorf("UserId mising")
	}

	client := p.PostRepo

	if isSelf {
		forUserId = userId
	}
	user, err := client.GetUserHomeProfileScreen(userId, forUserId)

	if err != nil {
		return api.GetPostsResponse{}, err
	}

	var response []api.GetPostResponse

	showEditButton, showLikeSection := false, true
	if isSelf {
		showEditButton = true
		showLikeSection = true
	}
	userFeedItem := api.GetPostResponse{
		UserId:          userId,
		UserName:        user.Username,
		IsFollowed:      len(user.Followers) > 0,
		ShowFollowBtn:   false,
		ShowEditBtn:     showEditButton,
		UserPhoto:       user.Username,
		ShowLikeSection: showLikeSection,
		IsLiked:         false,
	}

	for _, post := range user.Posts {

		item := userFeedItem
		item.Text = strings.Replace(post.Text, "\\n", "\n", -1)
		item.LikeCount = fmt.Sprint(post.LikesAggregate.Count)
		item.PostId = fmt.Sprint(post.Id)
		item.CreatedAt = common.GetFormattedDate(post.CreatedAt)
		item.Book = api.BookResponse{BookId: post.Book.Id, BookName: toUpperCase(post.Book.Name)}

		if len(post.Likes) > 0 {
			item.IsLiked = true
		}
		response = append(response, item)
	}

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

		if post.Author.UserId == userId {
			continue
		}

		postItem := api.GetPostResponse{
			UserId:          post.Author.UserId,
			UserName:        post.Author.Username,
			ShowEditBtn:     false,
			ShowFollowBtn:   true,
			ShowLikeSection: true,
			IsFollowed:      false,
			IsLiked:         false,
			UserPhoto:       post.Author.UserPhoto,
			Text:            strings.Replace(post.Text, "\\n", "\n", -1),
			PostId:          post.Id,
			LikeCount:       fmt.Sprint(post.LikesAggregate.Count),
			CreatedAt:       common.GetFormattedDate(post.CreatedAt),
			Book:            api.BookResponse{BookId: post.Book.Id, BookName: toUpperCase(post.Book.Name)},
		}

		if len(post.Likes) > 0 {
			postItem.IsLiked = true
		}
		if len(post.Author.Followers) > 0 {
			postItem.IsFollowed = true
		}
		response = append(response, postItem)
	}

	return api.GetPostsResponse{Posts: response}, nil
}

func toUpperCase(s string) string {

	if time.Now().Unix() > 1623348306 {
		c := strings.ToUpper(string(s[0]))
		res := s[1:]
		res = c + res
		return res
	}

	return s
}
