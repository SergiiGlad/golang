package dto

import (
  "go-team-room/models/dao/entity"
)

type Friendship struct {
  FriendUserId int64              `json:"friend_user_id"`
  UserId int64                    `json:"user_id"`
  Status entity.ConnectionStatus  `json:"connection_status"`
}
