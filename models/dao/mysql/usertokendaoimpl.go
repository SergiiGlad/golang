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

//Return new instance of mysql token dao.
func newMySqlTokenDao(conn *sql.DB) (interfaces.UserTokenDao, error) {
  log.Debug("Start creating new token dao for connection {}", conn)
  if err := conn.Ping(); err != nil {
    conn.Close()
    log.Error("mysql: could not establish a good connection: {}, error: {}", conn, err)
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &mysqlUserTokenDao{
    conn: conn,
  }

  var err error

  if db.insert, err = conn.Prepare(insertTokenStatement); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", insertTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.update, err = conn.Prepare(updateTokenStatement); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", updateTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.delete, err = conn.Prepare(deleteTokenStatement); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", deleteTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byid, err = conn.Prepare(findTokenByIdStatement); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", findTokenByIdStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.bytoken, err = conn.Prepare(findTokenByToken); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", findTokenByToken, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.getAll, err = conn.Prepare(getAllTokensStatement); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", getAllTokensStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.forUser, err = conn.Prepare(findTokenForUser); err != nil {
    log.Error("mysql: cant prepare {}, error: {}", findTokenForUser, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  log.Debug("Finish creating new token dao: {}", db)
  return db, nil
}

const getAllTokensStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens`

//Return all tokens.
func (tDao mysqlUserTokenDao) GetTokens() ([] entity.UserToken, error) {
  log.Debug("Get all tokens statement: {}", getAllTokensStatement)
  rows, err := tDao.getAll.Query()
  if err != nil {
    log.Error("Failed to get all tokens, err {}", err)
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  log.Debug("Get all tokens statement return rows {}", rows)
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      log.Error("Get all token statement fails, mysql: could not read row: {}", err)
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    log.Debug("Successfully read row for get all tokens statement token: {}", token)
    tokens = append(tokens, *token)
  }
  log.Info("Get all tokens returned: {} ", tokens)
  return tokens, nil
}

const insertTokenStatement = `INSERT INTO user_tokens (token, email, is_active, user_id) VALUES (?, ?, ?, ?)`

//Insert new token.
func (tDao mysqlUserTokenDao) AddToken(token entity.UserToken) (int64, error) {
  log.Debug("Insert token statement: {}", insertTokenStatement)
  log.Debug("Inserting new token: {} in database ", token)
  r, err := execAffectingOneRow(tDao.insert, token.Token, token.Email, token.IsActive, token.UserId)
  if err != nil {
    log.Error("Fail to insert token in database error: {} ", err)
    return 0, err
  }
  log.Debug("Insert statement return row: {}, for token: {} ", r, token)
  lastInsertID, err := r.LastInsertId()
  if err != nil {
    log.Error("Fail to insert token in database error: msql: could not get last insert ID: {}", err)
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }
  token.ID = lastInsertID
  log.Info("Token {} inserted id is: {}", token, lastInsertID)
  return lastInsertID, nil
}

const deleteTokenStatement = `DELETE FROM user_tokens WHERE token_id = ?`

// Delete token for given id.
func (tDao mysqlUserTokenDao) DeleteToken(id int64) error {
  log.Debug("Delete token statement: {}", deleteTokenStatement)
  log.Debug("Deleting token for id: {}", id)
  _, err := execAffectingOneRow(tDao.delete, id)
  if err != nil {
    log.Error("Error deleting token for id {} , error: {}", id, err)
    return err
  }
  log.Info("Token for id:{} deleted.", id)
  return nil
}

const updateTokenStatement = `UPDATE user_tokens SET
  token = ?, email = ?, is_active = ?, user_id = ? WHERE token_id = ?`

// Update token.
func (tDao mysqlUserTokenDao) UpdateToken(id int64, token *entity.UserToken) error {
  log.Debug("Update token statement: {}", updateTokenStatement)
  log.Debug("Updating token with id: {} to token: {}", id, token)
  _, err := execAffectingOneRow(tDao.update, token.Token, token.Email, token.IsActive, token.UserId, id)
  if err != nil {
    log.Error("Failed to update token with id: {}, to token: {}, error: {}", id, token, err)
    return err
  }
  log.Info("Successfully update token with id: {} to token: {}", id, token)
  return nil
}

const findTokenByIdStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token_id = ?`

// Find token with given id.
func (tDao mysqlUserTokenDao) FindTokenById(id int64) (*entity.UserToken, error) {
  log.Debug("Find token by id statement: {}", findTokenByIdStatement)
  log.Debug("Start looking for token with id: {}", id)
  token, err := scanToken(tDao.byid.QueryRow(id))
  if err != nil {
    log.Error("Failed to find token with id: {}, error: {}", id, err)
    return nil, err
  }
  log.Info("Found token with id: {}, token: {}", id, token)
  return token, nil
}

const findTokenByToken = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token = ?`

//Find token with given token value.
func (tDao mysqlUserTokenDao) FindTokenByToken(token string) (*entity.UserToken, error) {
  log.Debug("Find token by token statement: {}", findTokenByToken)
  log.Debug("Start looking for token with token value: {}", token)
  t, err := scanToken(tDao.bytoken.QueryRow(token))
  if err != nil {
    log.Error("Failed to find token with token value: {}, error: {}", token, err)
    return nil, err
  }
  log.Info("Found token with token value: {}, token: {}", token, t)
  return t, nil
}

const findTokenForUser = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE user_id = ?`

// Find token for user with given id.
func (tDao mysqlUserTokenDao) FindTokensForUser(user_id int64) ([] entity.UserToken, error) {
  log.Debug("Find token for user statement: ", findTokenForUser)
  log.Debug("Start looking for toking with user_id: {}", user_id)
  rows, err := tDao.forUser.Query(user_id)
  if err != nil {
    log.Error("Failed to find tokens with user_id: {}, error: {}", user_id, err)
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      log.Error("Failed to find tokens with user_id: {} error: mysql: could not read row: {}", user_id, err)
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    log.Debug("Found new token with user_id: {}, token: {}", user_id, token)
    tokens = append(tokens, *token)
  }
  log.Info("Found tokens with user_id: {}, tokens: {}", user_id, tokens)
  return tokens, nil
}

// Close closes the database, freeing up any resources.
func (d *mysqlUserTokenDao) Close() {
  log.Debug("Closing connection for token dao: {}", d)
  d.conn.Close()
}

var (
  token_id      int64
  user_email    sql.NullString
  token         sql.NullString
  is_active     bool
  user_id_token int64
)
//Read sql row to token.
func scanToken(s rowScanner) (*entity.UserToken, error) {
  log.Debug("Scanning row for new token with scanner: {}", s)
  if err := s.Scan(&token_id, &user_email, &token, &is_active, &user_id_token); err != nil {
    log.Error("Failed to scan row for token, error: {}", err)
    return nil, err
  }
  token := entity.UserToken{
    token_id,
    user_email.String,
    token.String,
    is_active,
    user_id_token,
  }
  log.Debug("Successfully scanned row for token: {}", token)
  return &token, nil
}
