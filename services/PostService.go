package services

import (
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/query"
	"sort"
)

type PostsService struct {
	PostRepo query.PostRepo
}

func NewPostsService() PostsService {
	return PostsService{PostRepo: query.NewPostRepo()}
}

func (p *PostsService) GetPosts(userId, screenName string) ([]api.GetPostsResponse, error) {
	client := query.NewPostRepo()

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
