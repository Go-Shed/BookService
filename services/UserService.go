package services

import "shed/bookservice/repos/dgraph/query"

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
