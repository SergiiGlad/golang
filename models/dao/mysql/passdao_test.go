package mysql

import (
  "testing"
  "go-team-room/models/dao/entity"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "github.com/stretchr/testify/assert"
)

var passQueriesRegexes []string = []string {
  //insert password
  `INSERT INTO users_passwords \(password, password_created, user_id\) VALUES \((.+),(.+),(.+)\)`,
  //current password
  `SELECT (.+) FROM users_passwords WHERE user_id =(.+) ORDER BY password_created DESC LIMIT 1`,
  //all user passwords
  `SELECT (.+) FROM users_passwords WHERE user_id =(.+)`,
}

var passPreps = make(map[string] *sqlmock.ExpectedPrepare)

func TestInsertPass(t *testing.T) {
  password := entity.Password{
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
    passPreps[query] = prep
  }

  passPreps[query].ExpectExec().WithArgs(password.Password, password.CreatedAt, password.UserId).
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
    passPreps[query] = prep
  }

  passPreps[query].ExpectQuery().WillReturnRows(rows)

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
    passPreps[query] = prep
  }

  passPreps[query].ExpectQuery().WillReturnRows(rows)

  passRepo, err := newMySqlPassDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  password, err := passRepo.LastPassByUserId(1)
  assert.NoError(t, err)
  assert.NotNil(t, password)
}
