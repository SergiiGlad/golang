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
  log.Debugf("Start creating new token dao for connection %s", conn)
  if err := conn.Ping(); err != nil {
    conn.Close()
    log.Errorf("mysql: could not establish a good connection: %s, error: %s", conn, err)
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &mysqlUserTokenDao{
    conn: conn,
  }

  var err error

  if db.insert, err = conn.Prepare(insertTokenStatement); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", insertTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.update, err = conn.Prepare(updateTokenStatement); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", updateTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.delete, err = conn.Prepare(deleteTokenStatement); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", deleteTokenStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byid, err = conn.Prepare(findTokenByIdStatement); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", findTokenByIdStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.bytoken, err = conn.Prepare(findTokenByToken); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", findTokenByToken, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.getAll, err = conn.Prepare(getAllTokensStatement); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", getAllTokensStatement, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.forUser, err = conn.Prepare(findTokenForUser); err != nil {
    log.Errorf("mysql: cant prepare %s, error: %s", findTokenForUser, err)
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  log.Debugf("Finish creating new token dao: %s", db)
  return db, nil
}

const getAllTokensStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens`

//Return all tokens.
func (tDao mysqlUserTokenDao) GetTokens() ([] entity.UserToken, error) {
  log.Debugf("Get all tokens statement: %s", getAllTokensStatement)
  rows, err := tDao.getAll.Query()
  if err != nil {
    log.Errorf("Failed to get all tokens, error %s", err)
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  log.Debugf("Get all tokens statement return rows %s", rows)
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      log.Errorf("Get all token statement fails, mysql: could not read row: %s", err)
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    log.Debugf("Successfully read row for get all tokens statement token: %s", token)
    tokens = append(tokens, *token)
  }
  log.Infof("Get all tokens returned: %s ", tokens)
  return tokens, nil
}

const insertTokenStatement = `INSERT INTO user_tokens (token, email, is_active, user_id) VALUES (?, ?, ?, ?)`

//Insert new token.
func (tDao mysqlUserTokenDao) AddToken(token entity.UserToken) (int64, error) {
  log.Debugf("Insert token statement: %s", insertTokenStatement)
  log.Debugf("Inserting new token: %s in database ", token)
  r, err := execAffectingOneRow(tDao.insert, token.Token, token.Email, token.IsActive, token.UserId)
  if err != nil {
    log.Errorf("Fail to insert token in database error: %s", err)
    return 0, err
  }
  log.Debugf("Insert statement return row: %s, for token: %s ", r, token)
  lastInsertID, err := r.LastInsertId()
  if err != nil {
    log.Errorf("Fail to insert token in database error: msql: could not get last insert ID: %s", err)
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }
  token.ID = lastInsertID
  log.Infof("Token %s inserted id is: %s", token, lastInsertID)
  return lastInsertID, nil
}

const deleteTokenStatement = `DELETE FROM user_tokens WHERE token_id = ?`

// Delete token for given id.
func (tDao mysqlUserTokenDao) DeleteToken(id int64) error {
  log.Debugf("Delete token statement: %s", deleteTokenStatement)
  log.Debugf("Deleting token for id: %s", id)
  _, err := execAffectingOneRow(tDao.delete, id)
  if err != nil {
    log.Errorf("Error deleting token for id %s , error: %s", id, err)
    return err
  }
  log.Infof("Token for id:%s deleted.", id)
  return nil
}

const updateTokenStatement = `UPDATE user_tokens SET
  token = ?, email = ?, is_active = ?, user_id = ? WHERE token_id = ?`

// Update token.
func (tDao mysqlUserTokenDao) UpdateToken(id int64, token *entity.UserToken) error {
  log.Debugf("Update token statement: %s", updateTokenStatement)
  log.Debugf("Updating token with id: %s to token: %s", id, token)
  _, err := execAffectingOneRow(tDao.update, token.Token, token.Email, token.IsActive, token.UserId, id)
  if err != nil {
    log.Errorf("Failed to update token with id: %s, to token: %s, error: %s", id, token, err)
    return err
  }
  log.Infof("Successfully update token with id: %s to token: %s", id, token)
  return nil
}

const findTokenByIdStatement = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token_id = ?`

// Find token with given id.
func (tDao mysqlUserTokenDao) FindTokenById(id int64) (*entity.UserToken, error) {
  log.Debugf("Find token by id statement: %s", findTokenByIdStatement)
  log.Debugf("Start looking for token with id: %s", id)
  token, err := scanToken(tDao.byid.QueryRow(id))
  if err != nil {
    log.Errorf("Failed to find token with id: %s, error: %S", id, err)
    return nil, err
  }
  log.Infof("Found token with id: %s, token: %s", id, token)
  return token, nil
}

const findTokenByToken = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE token = ?`

//Find token with given token value.
func (tDao mysqlUserTokenDao) FindTokenByToken(token string) (*entity.UserToken, error) {
  log.Debugf("Find token by token statement: %s", findTokenByToken)
  log.Debugf("Start looking for token with token value: %s", token)
  t, err := scanToken(tDao.bytoken.QueryRow(token))
  if err != nil {
    log.Errorf("Failed to find token with token value: %s, error: %s", token, err)
    return nil, err
  }
  log.Infof("Found token with token value: %s, token: %s", token, t)
  return t, nil
}

const findTokenForUser = `SELECT token_id, token, email, is_active, user_id FROM user_tokens 
WHERE user_id = ?`

// Find token for user with given id.
func (tDao mysqlUserTokenDao) FindTokensForUser(user_id int64) ([] entity.UserToken, error) {
  log.Debugf("Find token for user statement: %s", findTokenForUser)
  log.Debugf("Start looking for toking with user_id: %s", user_id)
  rows, err := tDao.forUser.Query(user_id)
  if err != nil {
    log.Errorf("Failed to find tokens with user_id: %s, error: %s", user_id, err)
    return nil, err
  }
  defer rows.Close()

  tokens := [] entity.UserToken{}
  for rows.Next() {
    token, err := scanToken(rows)
    if err != nil {
      log.Errorf("Failed to find tokens with user_id: %s error: mysql: could not read row: %s", user_id, err)
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }
    log.Debugf("Found new token with user_id: %s, token: %s", user_id, token)
    tokens = append(tokens, *token)
  }
  log.Infof("Found tokens with user_id: %s, tokens: %s", user_id, tokens)
  return tokens, nil
}

// Close closes the database, freeing up any resources.
func (d *mysqlUserTokenDao) Close() {
  log.Debugf("Closing connection for token dao: %s", d)
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
  log.Debugf("Scanning row for new token with scanner: %s", s)
  if err := s.Scan(&token_id, &token,  &user_email, &is_active, &user_id_token); err != nil {
    log.Errorf("Failed to scan row for token, error: %s", err)
    return nil, err
  }
  token := entity.UserToken{
    token_id,
    user_email.String,
    token.String,
    is_active,
    user_id_token,
  }
  log.Debugf("Successfully scanned row for token: %s", token)
  return &token, nil
}
