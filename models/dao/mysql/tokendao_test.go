package mysql

import (
  "testing"
  "go-team-room/models/dao/entity"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "github.com/stretchr/testify/assert"
)

var tokenQueriesRegexes = []string{
  `INSERT INTO user_tokens \(token, email, is_active, user_id\) VALUES \((.+), (.+), (.+), (.+)\)`,
  `UPDATE user_tokens SET token = (.+), email = (.+), is_active = (.+), user_id = (.+) WHERE token_id = (.+)`,
  `DELETE FROM user_tokens WHERE token_id = (.+)`,
  `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE token_id = (.+)`,
  `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE token = (.+)`,
  `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens`,
  `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE user_id = (.+)`,
}
var preps map[string]*sqlmock.ExpectedPrepare = make(map[string]*sqlmock.ExpectedPrepare)

func TestGetToken(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()
  rows := sqlmock.NewRows(
    []string{"token_id", "token", "email", "is_active", "user_id"}).
    AddRow(1, "token1", "arsenzhd@gmail.com", false, 1).
    AddRow(2, "token2", "arsenzhd@gmail.com", true, 1)

  query := `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range tokenQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  tokenRepo, err := newMySqlTokenDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  tokens, err := tokenRepo.GetTokens()
  assert.NoError(t, err)
  assert.NotNil(t, tokens)
  assert.Equal(t, 2, len(tokens))

}

func TestGetTokenByToken(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()
  rows := sqlmock.NewRows(
    []string{"token_id", "email", "token", "is_active", "user_id"}).
    AddRow(1, "arsenzhd@gmail.com", "token1", false, 1)

  query := `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE token = (.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range tokenQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  tokenRepo, err := newMySqlTokenDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  token, err := tokenRepo.FindTokenByToken("token1");
  assert.NoError(t, err)
  assert.NotNil(t, token)
  assert.Equal(t, int64(1), token.ID)

}

func TestGetTokenById(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()
  rows := sqlmock.NewRows(
    []string{"token_id", "email", "token", "is_active", "user_id"}).
    AddRow(1, "arsenzhd@gmail.com", "token1", false, 1)

  query := `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE token_id = (.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range tokenQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  tokenRepo, err := newMySqlTokenDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  token, err := tokenRepo.FindTokenById(1);
  assert.NoError(t, err)
  assert.NotNil(t, token)
  assert.Equal(t, int64(1), token.ID)

}

func TestGetTokenUserId(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()
  rows := sqlmock.NewRows(
    []string{"token_id", "email", "token", "is_active", "user_id"}).
    AddRow(1, "arsenzhd@gmail.com", "token1", false, 1).
    AddRow(2, "arsenzhd@gmail.com", "token1", false, 1)

  query := `SELECT (.+), (.+), (.+), (.+), (.+) FROM user_tokens WHERE user_id = (.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range tokenQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  tokenRepo, err := newMySqlTokenDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  tokens, err := tokenRepo.FindTokensForUser(1);
  assert.NoError(t, err)
  assert.NotNil(t, tokens)
  assert.Equal(t, 2, len(tokens))
  assert.Equal(t, int64(1), tokens[0].ID)
  assert.Equal(t, int64(2), tokens[1].ID)
}

func TestAddToken(t *testing.T) {
  token := entity.UserToken{
    0,
    "email@gmail.com",
    "token1",
    true,
    1,
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `INSERT INTO user_tokens \(token, email, is_active, user_id\) VALUES \((.+), (.+), (.+), (.+)\)`
  var prep *sqlmock.ExpectedPrepare
  for _, query := range tokenQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectExec().WithArgs(token.Token, token.Email, token.IsActive, token.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

  tokenRepository, err := newMySqlTokenDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  insertedToken, err := tokenRepository.AddToken(token)
  assert.NoError(t, err)
  assert.Equal(t, int64(1), insertedToken)
}
