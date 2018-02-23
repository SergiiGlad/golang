package mysql

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao/entity"
)

//mysqlUserDao implements UserDao interface
type mysqlUserDao struct {
  conn *sql.DB

  insert      *sql.Stmt
  update      *sql.Stmt
  delete      *sql.Stmt
  forceDelete *sql.Stmt
  countByRole *sql.Stmt
  byid        *sql.Stmt
  byemail     *sql.Stmt
  byphone     *sql.Stmt
}

var _ interfaces.UserDao = &mysqlUserDao{}

//newMySqlUserDao creates new mysqlUserDao object by instantiating every statement field. Any statement
// field can then be used without repeating Prepare() performing before next db query.
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
  if db.forceDelete, err = conn.Prepare(forceDeleteStatement); err != nil {
    return nil, fmt.Errorf("mysql:prepare list: %v", err)
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
  if db.countByRole, err = conn.Prepare(coundByRoleStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  return db, nil
}

// Close closes the database, freeing up any resources.
func (d *mysqlUserDao) Close() {
  d.conn.Close()
}

const insertStatement = `INSERT INTO
  users_data (email, first_name, last_name, phone, role_in_network, account_status, avatar_ref)
  VALUES (?, ?, ?, ?, ?, ?, ?)`

func (d *mysqlUserDao)  AddUser(user *entity.User) (entity.User, error) {

  r, err := execAffectingOneRow(d.insert, user.Email, user.FirstName, user.LastName, user.Phone, user.Role,
    user.AccStatus, user.AvatarRef)

  if err != nil {
    return *user, err
  }

  lastInsertID, err := r.LastInsertId()

  if err != nil {
    return *user, fmt.Errorf("mysql: could not get last insertConnection ID: %v", err)
  }

  user.ID = lastInsertID

  return *user, nil
}


const updateStatement = `UPDATE users_data SET
  email = ?, first_name = ?, last_name = ?, phone = ?, role_in_network = ?, account_status = ?, avatar_ref = ?
  WHERE user_id = ?`

func (d *mysqlUserDao) UpdateUser(id int64, user *entity.User) (entity.User, error) {
  _, err := execAffectingOneRow(d.update, user.Email, user.FirstName, user.LastName, user.Phone, user.Role,
    user.AccStatus, user.AvatarRef, id)

  if err != nil {
    return *user, err
  }

  return *user, nil
}

//It just changes account_status without actual deleting table row
const deleteStatement = `UPDATE users_data SET account_status = 'deleted' WHERE user_id = ?`

func (d *mysqlUserDao) DeleteUser(id int64) error {
  _, err := execAffectingOneRow(d.delete, id)

  return err
}

const forceDeleteStatement = `DELETE FROM users_data WHERE user_id = ?`

func (d *mysqlUserDao) ForceDeleteUser(id int64) error {
  _, err := execAffectingOneRow(d.forceDelete, id)

  return err
}

const coundByRoleStatement = `SELECT COUNT(*) FROM users_data WHERE role_in_network = ?`

func (d *mysqlUserDao) CountByRole(role entity.Role) (int64, error) {
  var count int64
  err := d.countByRole.QueryRow(role).Scan(&count)
  if err != nil {
    return 0, err
  }

  return count, nil
}

const findByIdStatement = `SELECT * FROM users_data WHERE user_id = ?`

func (d *mysqlUserDao) FindUserById(id int64) (entity.User, error) {
  user, err := scanUser(d.byid.QueryRow(id))

  if err != nil {
    return user, err
  }

  return user, nil
}


const findByEmailStatement = `SELECT user_id, first_name, last_name, email, phone, role_in_network, account_status, avatar_ref FROM users_data WHERE email = ?`

func (d *mysqlUserDao) FindUserByEmail(email string) (entity.User, error) {
  user, err := scanUser(d.byemail.QueryRow(email))

  if err != nil {
    return user, err
  }

  return user, nil
}


const findByPhoneStatement = `SELECT * FROM users_data WHERE phone = ?`

func (d *mysqlUserDao) FindUserByPhone(phone string) (entity.User, error) {
  user, err := scanUser(d.byphone.QueryRow(phone))

  if err != nil {
    return user, err
  }

  return user, nil
}

//scanUser reads a user from a sql.Row or sql.Rows

var (
  user_id   int64
  firstName sql.NullString
  lastName  sql.NullString
  email     sql.NullString
  phone     sql.NullString
  role      sql.NullString
  accStat   sql.NullString
  avRef     sql.NullString
)

func scanUser(s rowScanner) (entity.User, error) {

  user := entity.User{}

  if err := s.Scan(&user_id, &firstName, &lastName, &email, &phone, &role, &accStat, &avRef); err != nil {
      return user, err
  }

  user = entity.User{
    user_id,
    email.String,
    firstName.String,
    lastName.String,
    phone.String,
    entity.Role(role.String),
    entity.AccountStatus(accStat.String),
    avRef.String,
  }

  return user, nil
}
