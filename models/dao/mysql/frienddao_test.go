package mysql

import (
  "testing"
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "github.com/stretchr/testify/assert"
)

var friendQueriesRegexes []string = []string {
  //select friends by user id
  `SELECT (.+) FROM friend_list WHERE user_id =(.+) AND connection_status='approved'`,
}

var friendPreps map [string] *sqlmock.ExpectedPrepare = make(map[string] *sqlmock.ExpectedPrepare)

func TestFriendsByUserID(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows(
    []string{"friend_user_id", "user_id", "connection_status"}).
    AddRow(1, 2, "approved").
    AddRow(3, 2, "approved").
    AddRow(4, 2, "approved").
    AddRow(5, 2, "approved").
    AddRow(6, 2, "approved")

  query := `SELECT (.+) FROM friend_list WHERE user_id =(.+) AND connection_status='approved'`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range userQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    preps[query] = prep
  }

  preps[query].ExpectQuery().WillReturnRows(rows)

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  friends, err := friendRepo.FriendsByUserID(2)
  assert.NoError(t, err)
  assert.NotNil(t, friends)
}
