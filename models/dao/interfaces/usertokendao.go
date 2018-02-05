package interfaces

import (
  "go-team-room/models/dao"
)

type UserTokenDao interface {
  GetTokens() ([] dao.UserToken, error)
  AddToken(token dao.UserToken) error
  DeleteToken(id int64) error
  UpdateToken(id int64, token *dao.UserToken) error

  FindTokenById(id int64) (*dao.UserToken, error)
  FindTokenByEmail(email string) (*dao.UserToken, error)
  FindTokensForUser(user_id int64) ([] dao.UserToken, error)

  UpdateTokenForEmail(email string, token *dao.UserToken) error
}
