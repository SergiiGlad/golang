package interfaces

import "go-team-room/models/dao/entity"

type FriendshipDao interface {
  FriendsByUserID(id int64) ([]entity.Friendship, error)
}
