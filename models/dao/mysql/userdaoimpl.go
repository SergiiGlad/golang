package mysql

import (
  "go-team-room/models/dao"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "go-team-room/models/dao/interfaces"
)

type mysqlUserDao struct {
  conn *sql.DB

  insert  *sql.Stmt
  update  *sql.Stmt
  delete  *sql.Stmt
  byid    *sql.Stmt
  byemail *sql.Stmt
  byphone *sql.Stmt
  friends *sql.Stmt
}

var _ interfaces.UserDao = &mysqlUserDao{}

func newMySqlUserDao(conn *sql.DB) (interfaces.UserDao, error) {

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &mysqlUserDao{
    conn: conn,
  }

  var err error

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
  if db.friends, err = conn.Prepare(findFriendsByUserId); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  return db, nil
}

// Close closes the database, freeing up any resources.
func (d *mysqlUserDao) Close() {
  d.conn.Close()
}

const insertStatement = `INSERT INTO
  users_data (email, first_name, second_name, phone, current_password, role_in_network, account_status, avatar_ref)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

func (d *mysqlUserDao) AddUser(user *dao.User) (int64, error) {
  r, err := execAffectingOneRow(d.insert, user.Email, user.FirstName, user.SecondName, user.Phone, user.CurrentPass, user.Role,
    user.AccStatus, user.AvatarRef)

  if err != nil {
    return 0, err
  }

  lastInsertID, err := r.LastInsertId()

  if err != nil {
    return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
  }

  user.ID = lastInsertID

  return lastInsertID, nil
}


const updateStatement = `UPDATE users_data SET
  email = ?, first_name = ?, second_name = ?, phone = ?, current_password = ?, role_in_network = ?, account_status = ?, avatar_ref = ?
  WHERE user_id = ?`

func (d *mysqlUserDao) UpdateUser(id int64, user *dao.User) error {
  _, err := execAffectingOneRow(d.update, user.Email, user.FirstName, user.SecondName, user.Phone, user.CurrentPass, user.Role,
    user.AccStatus, user.AvatarRef, id)

  return err
}

//It just changes account_status without actual deleting table row
const deleteStatement = `UPDATE users_data SET account_status = 'deleted' WHERE user_id = ?`

func (d *mysqlUserDao) DeleteUser(id int64) error {
  _, err := execAffectingOneRow(d.delete, id)

  return err
}


const findByIdStatement = `SELECT * FROM users_data WHERE user_id = ?`

func (d *mysqlUserDao) FindUserById(id int64) (*dao.User, error) {
  user, err := scanUser(d.byid.QueryRow(id))

  if err != nil {
    return nil, err
  }

  return user, nil
}


const findByEmailStatement = `SELECT * FROM users_data WHERE email = ?`

func (d *mysqlUserDao) FindUserByEmail(email string) (*dao.User, error) {
  user, err := scanUser(d.byemail.QueryRow(email))

  if err != nil {
    return nil, err
  }

  return user, nil
}


const findByPhoneStatement = `SELECT * FROM users_data WHERE phone = ?`

func (d *mysqlUserDao) FindUserByPhone(phone string) (*dao.User, error) {
  user, err := scanUser(d.byphone.QueryRow(phone))

  if err != nil {
    return nil, err
  }

  return user, nil
}

const findFriendsByUserId = `SELECT friend_id FROM friend_list WHERE user_id = ?;`

func (d *mysqlUserDao) FriendsByUserID(id int64) ([]int64, error) {
  rows, err := d.friends.Query()

  if err != nil {
    return nil, err
  }
  defer rows.Close()

  friendIds := []int64{}
  var friendId int64
  for rows.Next() {
    err = rows.Scan(&friendId)
    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }

    friendIds = append(friendIds, friendId)
  }

  return friendIds, nil
}

//scanUser reads a user from a sql.Row or sql.Rows

var (
  user_id    int64
  email      sql.NullString
  firstName  sql.NullString
  secondName sql.NullString
  phone      sql.NullString
  pass       sql.NullString
  role       sql.NullString
  accStat    sql.NullString
  avRef      sql.NullString
)

func scanUser(s rowScanner) (*dao.User, error) {

  if err := s.Scan(&user_id, &email, &firstName, &secondName, &phone,
    &pass, &role, &accStat, &avRef); err != nil {
      return nil, err
  }

  user := dao.User{
    user_id,
    email.String,
    firstName.String,
    secondName.String,
    phone.String,
    pass.String,
    dao.Role(role.String),
    dao.AccountStatus(accStat.String),
    avRef.String,
  }

  return &user, nil
}
