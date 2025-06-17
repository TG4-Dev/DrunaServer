package services

import (
	"BlobbyServer/pkg/repositories"
	"errors"
)

type friendService struct{}

var FriendService = friendService{}

func (f *friendService) FriendRequest(user_id int, username string) error {
	friend_id, err := repositories.UsersRepo.ExistsByUsername(username)
	if err != nil {
		return errors.New("user doesn't exists")
	}

	err = repositories.UsersRepo.CreateFriendRequest(user_id, friend_id)
	if err != nil {
		return errors.New("something gone wrong")
	}

	return err
}
