package mysql

import (
  "testing"
  "go-team-room/models/dao"
  sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "github.com/stretchr/testify/assert"
)

var userQueriesRegexes []string = []string{
  //insert user
  `INSERT INTO users_data \(email, first_name, last_name, phone, role_in_network, account_status, avatar_ref\) VALUES \((.+),(.+),(.+),(.+),(.+),(.+),(.+)\)`,
  //update user
  `UPDATE users_data SET email =(.+), first_name =(.+), last_name =(.+), phone =(.+), role_in_network =(.+), account_status =(.+), avatar_ref =(.+) WHERE user_id =(.+)`,
  //delete statement
  `UPDATE users_data SET account_status = 'deleted' WHERE user_id =(.+)`,
  //forese delete statement
  `DELETE FROM users_data WHERE user_id =(.+)`,
  //find by id
  `SELECT (.+) FROM users_data WHERE user_id =(.+)`,
  //find by email
  `SELECT (.+) FROM users_data WHERE email =(.+)`,
  //find by phone
  `SELECT (.+) FROM users_data WHERE phone =(.+)`,
  //user friends
  `SELECT friend_id FROM friend_list WHERE user_id =(.+)`,
  //count row by user-role
  `SELECT COUNT\(\*\) FROM users_data WHERE role_in_network =(.+)`,
}

var passQueriesRegexes []string = []string {
  //insert password
  `INSERT INTO users_passwords \(password, password_created, user_id\) VALUES \((.+),(.+),(.+)\)`,
  //current password
  `SELECT (.+) FROM users_passwords WHERE user_id =(.+) ORDER BY password_created DESC LIMIT 1`,
  //all user passwords
  `SELECT (.+) FROM users_passwords WHERE user_id =(.+)`,
}

var preps map [string] *sqlmock.ExpectedPrepare = make(map[string] *sqlmock.ExpectedPrepare)

func TestAddUser(t *testing.T) {
  user := dao.User {
    0,
    "email@gmail.com",
    "Name",
    "Surname",
    "+3805436857",
    dao.UserRole,
    dao.Active,
    "",
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `INSERT INTO users_data \(email, first_name, last_name, phone, role_in_network, account_status, avatar_ref\) VALUES \((.+),(.+),(.+),(.+),(.+),(.+),(.+)\)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectExec().WithArgs(user.Email, user.FirstName,
    user.LastName, user.Phone, user.Role, user.AccStatus, user.AvatarRef).WillReturnResult(sqlmock.NewResult(1, 1))

  userRepository, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  insertedUser, err := userRepository.AddUser(&user)
  assert.NoError(t, err)
  assert.Equal(t, int64(1), insertedUser.ID)
}

func TestFindUserById(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"user_id", "email", "first_name", "last_name", "phone", "role_in_network", "account_status", "avatar_ref"}).
      AddRow(1, "email@gmail.com", "Name", "Surname", "phone", "user", "active", "")

  query := `SELECT (.+) FROM users_data WHERE user_id =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  userRepo, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  user, err := userRepo.FindUserById(1)
  assert.NoError(t, err)
  assert.NotNil(t, user)
}

func TestFindUserByEmail(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"user_id", "email", "first_name", "last_name", "phone", "role_in_network", "account_status", "avatar_ref"}).
    AddRow(1, "email@gmail.com", "Name", "Surname", "phone", "user", "active", "")

  query := `SELECT (.+) FROM users_data WHERE email =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  userRepo, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  user, err := userRepo.FindUserByEmail("email@gmail.com")
  assert.NoError(t, err)
  assert.NotNil(t, user)
}

func TestFindUserByPhone(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"user_id", "email", "first_name", "last_name", "phone", "role_in_network", "account_status", "avatar_ref"}).
    AddRow(1, "email@gmail.com", "Name", "Surname", "phone", "user", "active", "")

  query := `SELECT (.+) FROM users_data WHERE email =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  userRepo, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  user, err := userRepo.FindUserByEmail("email@gmail.com")
  assert.NoError(t, err)
  assert.NotNil(t, user)
}

func TestFindFriendByUserId(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"friend_id"}).
    AddRow(1).
      AddRow(3).
        AddRow(4).
          AddRow(5).
            AddRow(6)

  query := `SELECT friend_id FROM friend_list WHERE user_id =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  userRepo, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  friends, err := userRepo.FriendsByUserID(2)
  assert.NoError(t, err)
  assert.NotNil(t, friends)
}

func TestInsertPass(t *testing.T) {
  password := dao.Password{
    0,
    "123456",
    "1999-01-01 10:10:01",
    1,
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `INSERT INTO users_passwords \(password, password_created, user_id\) VALUES \((.+),(.+),(.+)\)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range passQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectExec().WithArgs(password.Password, password.CreatedAt, password.UserId).
    WillReturnResult(sqlmock.NewResult(1, 1))

  passwordRepository, err := newMySqlPassDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  lastId, err := passwordRepository.InsertPass(&password)
  assert.NoError(t, err)
  assert.Equal(t, int64(1), lastId)
}

func TestLatestUserPassword(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"password_id", "password", "password_created", "user_id"}).
    AddRow(1, "password1", "1999-01-01 10:10:01", 1)

  query := `SELECT (.+) FROM users_passwords WHERE user_id =(.+) ORDER BY password_created DESC LIMIT 1`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range passQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  passRepo, err := newMySqlPassDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  password, err := passRepo.LastPassByUserId(1)
  assert.NoError(t, err)
  assert.NotNil(t, password)
}

func TestUserPasswords(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"password_id", "password", "password_created", "user_id"}).
    AddRow(1, "password1", "1999-01-01 10:10:01", 1).
      AddRow(2, "password2", "2000-01-01 10:10:01", 1)

  query := `SELECT (.+) FROM users_passwords WHERE user_id =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range passQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  passRepo, err := newMySqlPassDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  password, err := passRepo.LastPassByUserId(1)
  assert.NoError(t, err)
  assert.NotNil(t, password)
}
