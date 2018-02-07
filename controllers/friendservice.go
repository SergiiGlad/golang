package controllers

import (
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao/entity"
)

type FriendService struct {
  DB interfaces.FriendshipDao
}

func (fs *FriendService) GetUserFriends(id int64) ([]entity.User, error) {

  friendships, err := fs.DB.FriendsByUserID(id)
  if err != nil {
    return nil, err
  }

  friends, err := fs.

  return fs.DB.FriendsByUserID(id)
}

func (fs *FriendService) GetUserFriedsIds(id int64)([]int64, error) {
  var fids []int64
  friends, err := fs.GetUserFriends(id)
  if err != nil {
    return fids, err
  }

  for _, friend := range friends {
    fids = append(fids, friend.FriendUserId)
  }

  return fids, nil
}
