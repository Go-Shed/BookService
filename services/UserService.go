package services

import (
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/query"
)

type UserService struct {
	UserRepo query.UserRepo
}

func NewUserService() UserService {
	return UserService{UserRepo: query.NewUserRepo()}
}

func (u *UserService) FollowUser(userId, userIdToFollow string) error {
	err := u.UserRepo.FollowUser(userId, userIdToFollow)
	return err
}

func (u *UserService) UnfollowUser(userId, userIdToUnFollow string) error {
	err := u.UserRepo.UnFollowUser(userId, userIdToUnFollow)
	return err
}

func (u *UserService) SearchUser(userID, username string) api.SearchUserResponse {
	users := u.UserRepo.SearchUser(userID, username)

	var response []api.SearchResult

	for _, user := range users {

		if user.UserId == userID {
			continue
		}

		result := api.SearchResult{
			UserPhoto:        user.UserPhoto,
			UserId:           user.UserId,
			UserName:         user.Username,
			ShowFollowButton: false,
			IsFollowed:       false,
		}

		if len(user.Followers) > 0 {
			result.IsFollowed = true
			result.ShowFollowButton = true
		}

		response = append(response, result)
	}

	return api.SearchUserResponse{SearchResults: response}
}

func (u *UserService) FetchProfile(userId, profileUserId string, isSelf bool) (api.FetchProfileResponse, error) {

	if !isSelf && len(profileUserId) == 0 {
		return api.FetchProfileResponse{}, fmt.Errorf("profleUser id must exist")
	}

	if isSelf {
		profileUserId = userId
	}

	response, err := u.UserRepo.FetchProfile(profileUserId, userId)

	if err != nil {
		return api.FetchProfileResponse{}, fmt.Errorf("something went wrong")
	}

	isFollowing := false
	if len(response.Followers) > 0 && !isSelf {
		isFollowing = true
	}

	// fmt.Printf("%+v", response)

	return api.FetchProfileResponse{
		UserPhoto:   response.UserPhoto,
		Email:       response.Email,
		UserName:    response.Username,
		Followers:   response.FollowersAggregate.Count,
		Following:   response.FollowingAggregate.Count,
		IsFollowing: isFollowing,
	}, nil
}

func (u *UserService) GetFollowers(userId, profileUserId string, isSelf bool) (api.GetFollowersResponse, error) {

	if !isSelf && len(profileUserId) == 0 {
		return api.GetFollowersResponse{}, fmt.Errorf("profleUser id must exist")
	}

	if isSelf {
		profileUserId = userId
	}

	response, err := u.UserRepo.GetFollowers(profileUserId)

	if err != nil {
		return api.GetFollowersResponse{}, fmt.Errorf("something went wrong")
	}

	if len(response.Followers) == 0 {
		return api.GetFollowersResponse{}, nil
	}

	var follows []api.FollowItem

	for _, item := range response.Followers {

		followItem := api.FollowItem{
			UserPhoto: item.UserPhoto,
			UserName:  item.Username,
			UserId:    item.UserId,
		}

		follows = append(follows, followItem)

	}

	return api.GetFollowersResponse{Follows: follows, TotalFollowers: len(follows)}, nil
}

func (u *UserService) GetFollowing(userId, profileUserId string, isSelf bool) (api.GetFollowingResponse, error) {

	if !isSelf && len(profileUserId) == 0 {
		return api.GetFollowingResponse{}, fmt.Errorf("profleUser id must exist")
	}

	if isSelf {
		profileUserId = userId
	}

	response, err := u.UserRepo.GetFollowing(profileUserId)

	if err != nil {
		return api.GetFollowingResponse{}, fmt.Errorf("something went wrong")
	}

	if len(response.Following) == 0 {
		return api.GetFollowingResponse{}, nil
	}

	var follows []api.FollowItem

	for _, item := range response.Following {

		followItem := api.FollowItem{
			UserPhoto: item.UserPhoto,
			UserName:  item.Username,
			UserId:    item.UserId,
		}

		follows = append(follows, followItem)

	}

	return api.GetFollowingResponse{Follows: follows, TotalFollowing: len(follows)}, nil
}
