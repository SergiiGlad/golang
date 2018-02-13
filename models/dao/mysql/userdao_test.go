package mysql

import (
  "go-team-room/models/dao/entity"
  "testing"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "github.com/stretchr/testify/assert"
)

var userQueriesRegexes []string = []string{
  //insert user
  `INSERT INTO users_data \(email, first_name, last_name, phone, role_in_network, account_status, avatar_ref\) VALUES \((.+),(.+),(.+),(.+),(.+),(.+),(.+)\)`,
  //updateStatus user
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
  //count row by user-role
  `SELECT COUNT\(\*\) FROM users_data WHERE role_in_network =(.+)`,
}

var tokenPreps map[string]*sqlmock.ExpectedPrepare = make(map[string]*sqlmock.ExpectedPrepare)


func TestAddUser(t *testing.T) {
  user := entity.User{
    0,
    "email@gmail.com",
    "Name",
    "Surname",
    "+3805436857",
    entity.UserRole,
    entity.Active,
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
    tokenPreps[query] = prep
  }

  tokenPreps[query].ExpectExec().WithArgs(user.Email, user.FirstName,
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
    tokenPreps[query] = prep
  }

  tokenPreps[query].ExpectQuery().WillReturnRows(rows)

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
    tokenPreps[query] = prep
  }

  tokenPreps[query].ExpectQuery().WillReturnRows(rows)

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
    tokenPreps[query] = prep
  }

  tokenPreps[query].ExpectQuery().WillReturnRows(rows)

  userRepo, err := newMySqlUserDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  user, err := userRepo.FindUserByEmail("email@gmail.com")
  assert.NoError(t, err)
  assert.NotNil(t, user)
}
