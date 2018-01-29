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

  insert        *sql.Stmt
  lastPassword  *sql.Stmt
  passwords     *sql.Stmt
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
  if dao.lastPassword, err = conn.Prepare(lastPassStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if dao.passwords, err = conn.Prepare(passwordsStatement); err != nil {
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

const lastPassStatement = `SELECT * FROM users_passwords WHERE user_id = ?
                           ORDER BY password_created DESC LIMIT 1`

func (d *mysqlPassDaoImpl) LastPassByUserId(id int64) (dao.Password, error) {
  pass, err := scanPass(d.lastPassword.QueryRow(id))

  if err != nil {
    return *pass, err
  }

  return *pass, nil
}

const passwordsStatement = `SELECT * FROM users_passwords WHERE user_id = ?`

func (d *mysqlPassDaoImpl) PasswdsByUserId(id int64) ([]dao.Password, error) {
  rows, err := d.passwords.Query(id)

  if err != nil {
    return nil, err
  }
  rows.Close()

  passwords := []dao.Password{}
  var pass dao.Password

  for rows.Next() {
    err = rows.Scan(&id, &password, created_at, usr_id)

    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }

    pass.ID = id
    pass.Password = password.String
    pass.CreatedAt = password.String
    pass.UserId = usr_id

    passwords = append(passwords, pass)
  }

  return passwords, nil
}

var (
  id        int64
  password  sql.NullString
  created_at sql.NullString
  usr_id    int64
)

func scanPass(s rowScanner) (*dao.Password, error) {

  if err := s.Scan(&id, &password, &created_at, &usr_id); err != nil {
    return nil, err
  }

  return &dao.Password{
    id,
    password.String,
    created_at.String,
    usr_id,
  }, nil
}
