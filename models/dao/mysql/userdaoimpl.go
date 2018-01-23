package mysql

import (
  "go-team-room/models/dao"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "go-team-room/models/dao/interfaces"
)

type mysqlUserDB struct {
  conn *sql.DB

  insert  *sql.Stmt
  update  *sql.Stmt
  delete  *sql.Stmt
  byid    *sql.Stmt
  byemail *sql.Stmt
  byphone *sql.Stmt
}

var _ interfaces.UserDatabase = &mysqlUserDB{}

func newMySqlUserDB(conn *sql.DB) (interfaces.UserDatabase, error) {

  // Check database and table exists. If not, create it.
  err := ensureTableExists()

  if err != nil {
    return nil, err
  }

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &mysqlUserDB{
    conn: conn,
  }

  if db.insert, err = conn.Prepare(insertStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.update, err = conn.Prepare(updateStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.delete, err = conn.Prepare(deleteStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byid, err = conn.Prepare(findByIdStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byemail, err = conn.Prepare(findByEmailStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.byphone, err = conn.Prepare(findByPhoneStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }

  return db, nil
}

// Close closes the database, freeing up any resources.
func (db *mysqlUserDB) Close() {
  db.conn.Close()
}

const insertStatement = `INSERT INTO
  users_data (email, first_name, second_name, phone, current_password, role_in_network, account_status, avatar_ref)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

func (db *mysqlUserDB) AddUser(user *dao.User) (int64, error) {
  r, err := execAffectingOneRow(db.insert, user.Email, user.FirstName, user.SecondName, user.CurrentPass, user.Role,
    user.AccStatus, user.AvatarRef)

  if err != nil {
    return 0, err
  }

  lastInsertID, err := r.LastInsertId()

  if err != nil {
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }
  return lastInsertID, nil
}


const updateStatement = `UPDATE users_data SET
  email = ?, first_name = ?, second_name = ?, phone = ?, current_password = ?, role_in_network = ?, account_status = ?, avatar_ref = ?
  WHERE user_id = ?`

func (db *mysqlUserDB) UpdateUser(id int, user *dao.User) error {
  _, err := execAffectingOneRow(db.update, user.Email, user.FirstName, user.SecondName, user.CurrentPass, user.Role,
    user.AccStatus, user.AvatarRef, id)

  return err
}

//It just changes account_status without actual deleting table row
const deleteStatement = `UPDATE users_data SET account_status = 'deleted' WHERE user_id = ?`

func (db *mysqlUserDB) DeleteUser(id int) error {
  _, err := execAffectingOneRow(db.delete, id)

  return err
}


const findByIdStatement = `SELECT * FROM users_data WHERE id = ?`

func (db *mysqlUserDB) FindUserById(id int) (*dao.User, error) {
  user, err := scanUser(db.byid.QueryRow(id))

  if err != nil {
    return nil, err
  }

  return user, nil
}


const findByEmailStatement = `SELECT * FROM users_data WHERE email = ?`

func (db *mysqlUserDB) FindUserByEmail(email string) (*dao.User, error) {
  user, err := scanUser(db.byemail.QueryRow(email))

  if err != nil {
    return nil, err
  }

  return user, nil
}


const findByPhoneStatement = `SELECT * FROM users_data WHERE phone = ?`

func (db *mysqlUserDB) FindUserByPhone(phone string) (*dao.User, error) {
  user, err := scanUser(db.byphone.QueryRow(phone))

  if err != nil {
    return nil, err
  }

  return user, nil
}

//scanUser reads a user from a sql.Row or sql.Rows
func scanUser(s rowScanner) (*dao.User, error) {

  user := dao.User{}

  if err := s.Scan(&user.ID, &user.Email, &user.FirstName, &user.SecondName, &user.Phone,
    &user.CurrentPass, &user.Role, &user.AccStatus, &user.AvatarRef); err != nil {
      return nil, err
  }

  return &user, nil
}
