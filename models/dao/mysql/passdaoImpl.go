package mysql

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao"
  "fmt"
)

type mysqlPassDaoImpl struct {
  conn *sql.DB

  insert  *sql.Stmt
  //update  *sql.Stmt
  //delete  *sql.Stmt
  //byuid   *sql.Stmt
}

var _ interfaces.PasswordDao = &mysqlPassDaoImpl{}

func newMySqlPassDao(conn *sql.DB) (interfaces.PasswordDao, error) {

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  dao := &mysqlPassDaoImpl{
    conn: conn,
  }

  var err error

  if dao.insert, err = conn.Prepare(insertPassStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }

  return dao, nil
}

const insertPassStatement = `INSERT INTO
  users_passwords (password, password_created, user_id)
  VALUES (?, ?, ?)`

func (d *mysqlPassDaoImpl) InsertPass(pass *dao.Password) (int64, error) {
  r, err := execAffectingOneRow(d.insert, pass.Password, pass.CreatedAt, pass.UserId)

  if err != nil {
    return 0, err
  }

  lastInsertID, err := r.LastInsertId()

  if err != nil {
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }

  pass.ID = lastInsertID

  return lastInsertID, nil
}

func (db *mysqlPassDaoImpl) Close() {
  db.conn.Close()
}

