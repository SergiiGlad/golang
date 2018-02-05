package mysql

import (
  "database/sql"
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao"
)

type mysqlUserTokenDao struct {
  conn *sql.DB

  insert  *sql.Stmt
  update  *sql.Stmt
  delete  *sql.Stmt
  byid    *sql.Stmt
  byemail *sql.Stmt
  byphone *sql.Stmt
  friends *sql.Stmt
}

var _ interfaces.UserTokenDao = &mysqlUserTokenDao{}

func (tDao mysqlUserTokenDao) GetTokens() ([] dao.UserToken, error) {
  return nil, nil
}
func (tDao mysqlUserTokenDao) AddToken(token dao.UserToken) error {
  return nil
}
func (tDao mysqlUserTokenDao) DeleteToken(id int64) error {
  return nil
}
func (tDao mysqlUserTokenDao) UpdateToken(id int64, token *dao.UserToken) error {
  return nil
}

func (tDao mysqlUserTokenDao) FindTokenById(id int64) (*dao.UserToken, error) {
  return nil, nil
}
func (tDao mysqlUserTokenDao) FindTokenByEmail(email string) (*dao.UserToken, error) {
  return nil, nil
}
func (tDao mysqlUserTokenDao) FindTokensForUser(user_id int64) ([] dao.UserToken, error) {
  return nil, nil
}

func (tDao mysqlUserTokenDao) UpdateTokenForEmail(email string, token *dao.UserToken) error {
  return nil
}

func newMySqlTokenDao(db *sql.DB) (interfaces.UserTokenDao, error) {
  return nil, nil
}
