package services

import (
	"fmt"
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/model"
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

func (u *UserService) SearchUser(username string) model.User {
	user := u.UserRepo.SearchUser(username)
	return user
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
