package interfaces

import (
  "go-team-room/models/dao"
)

type UserTokenDao interface {
  GetTokens() ([] dao.UserToken, error)
  AddToken(token dao.UserToken) (int64, error)
  DeleteToken(id int64) error
  UpdateToken(id int64, token *dao.UserToken) error

  FindTokenById(id int64) (*dao.UserToken, error)
  FindTokenByToken(token string) (*dao.UserToken, error)
  FindTokensForUser(user_id int64) ([] dao.UserToken, error)
}
