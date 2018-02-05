package interfaces

import (
  "go-team-room/models/dao/entity"
)

type UserTokenDao interface {
  GetTokens() ([] entity.UserToken, error)
  AddToken(token entity.UserToken) (int64, error)
  DeleteToken(id int64) error
  UpdateToken(id int64, token *entity.UserToken) error

  FindTokenById(id int64) (*entity.UserToken, error)
  FindTokenByToken(token string) (*entity.UserToken, error)
  FindTokensForUser(user_id int64) ([] entity.UserToken, error)
}
