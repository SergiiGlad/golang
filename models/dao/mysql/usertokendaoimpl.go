package mysql

import (
  "database/sql"
  "go-team-room/models/dao/interfaces"
  "fmt"
  "go-team-room/models/dao/entity"
)

type mysqlUserTokenDao struct {
  conn *sql.DB

  getAll  *sql.Stmt
  insert  *sql.Stmt
  update  *sql.Stmt
  delete  *sql.Stmt
  byid    *sql.Stmt
  bytoken *sql.Stmt
  forUser *sql.Stmt
}

var _ interfaces.UserTokenDao = &mysqlUserTokenDao{}

func newMySqlTokenDao(conn *sql.DB) (interfaces.UserTokenDao, error) {

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &mysqlUserTokenDao{
    conn: conn,
  }

  var err error

  if db.insert, err = conn.Prepare(insertTokenStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.update, err = conn.Prepare(updateTokenStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.delete, err = conn.Prepare(deleteTokenStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byid, err = conn.Prepare(findTokenByIdStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.bytoken, err = conn.Prepare(findTokenByToken); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.getAll, err = conn.Prepare(getAllTokensStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.forUser, err = conn.Prepare(findTokenForUser); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  return db, nil
}

const getAllTokensStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens`

func (tDao mysqlUserTokenDao) GetTokens() ([] entity.UserToken, error) {
  rows, err := tDao.getAll.Query()
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    tokens = append(tokens, *token)
  }
  return tokens, nil
}

const insertTokenStatement = `INSERT INTO user_tokens (token, email, is_active, user_id) VALUES (?, ?, ?, ?)`

func (tDao mysqlUserTokenDao) AddToken(token entity.UserToken) (int64, error) {
  r, err := execAffectingOneRow(tDao.insert, token.Token, token.Email, token.IsActive, token.UserId)
  if err != nil {
    return 0, err
  }
  lastInsertID, err := r.LastInsertId()
  if err != nil {
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }
  token.ID = lastInsertID
  return lastInsertID, nil
}

const deleteTokenStatement = `DELETE FROM user_tokens WHERE token_id = ?`

func (tDao mysqlUserTokenDao) DeleteToken(id int64) error {
  _, err := execAffectingOneRow(tDao.delete, id)
  if err != nil {
    return err
  }
  return nil
}

const updateTokenStatement = `UPDATE user_tokens SET
  token = ?, email = ?, is_active = ?, user_id = ? WHERE token_id = ?`

func (tDao mysqlUserTokenDao) UpdateToken(id int64, token *entity.UserToken) error {
  _, err := execAffectingOneRow(tDao.update, token.Token, token.Email, token.IsActive, token.UserId, id)
  if err != nil {
    return err
  }
  return nil
}

const findTokenByIdStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token_id = ?`

func (tDao mysqlUserTokenDao) FindTokenById(id int64) (*entity.UserToken, error) {
  token, err := scanToken(tDao.byid.QueryRow(id))
  if err != nil {
    return nil, err
  }
  return token, nil
}

const findTokenByToken = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token = ?`

func (tDao mysqlUserTokenDao) FindTokenByToken(token string) (*entity.UserToken, error) {
  t, err := scanToken(tDao.bytoken.QueryRow(token))
  if err != nil {
    return nil, err
  }
  return t, nil
}

const findTokenForUser = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE user_id = ?`

func (tDao mysqlUserTokenDao) FindTokensForUser(user_id int64) ([] entity.UserToken, error) {
  rows, err := tDao.forUser.Query(user_id)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    tokens = append(tokens, *token)
  }
  return tokens, nil
}

// Close closes the database, freeing up any resources.
func (d *mysqlUserTokenDao) Close() {
  d.conn.Close()
}

var (
  token_id      int64
  user_email    sql.NullString
  token         sql.NullString
  is_active     bool
  user_id_token int64
)

func scanToken(s rowScanner) (*entity.UserToken, error) {
  if err := s.Scan(&token_id, &user_email, &token, &is_active, &user_id_token); err != nil {
    return nil, err
  }
  token := entity.UserToken{
    token_id,
    user_email.String,
    token.String,
    is_active,
    user_id_token,
  }
  return &token, nil
}
