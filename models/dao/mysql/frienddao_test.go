package mysql

import (
  "gopkg.in/DATA-DOG/go-sqlmock.v1"
  "testing"
  "github.com/stretchr/testify/assert"
  "go-team-room/models/dao/entity"
)

var friendQueriesRegexes []string = []string {
  `SELECT (.+) FROM friend_list WHERE friend_user_id =(.+) AND user_id =(.+)`,

  `INSERT INTO friend_list \(friend_user_id, user_id, connection_status\) VALUES \((.+),(.+),(.+)\)`,

  `UPDATE friend_list SET connection_status =(.+) WHERE friend_user_id =(.+) AND user_id =(.+)`,

  `DELETE FROM friend_list WHERE friend_user_id IN \((.+),(.+)\) AND user_id IN \((.+),(.+)\)`,

  `SELECT users_data.* FROM users_data JOIN friend_list ON \(users_data.user_id = friend_list.user_id\)
   WHERE friend_list.friend_user_id =(.+) AND friend_list.connection_status='approved'`,

   `SELECT users_data.* FROM users_data JOIN friend_list ON \(users_data.user_id = friend_list.user_id\)
   WHERE friend_list.friend_user_id =(.+) AND friend_list.connection_status='waiting'`,
}

var friendPreps = make(map[string] *sqlmock.ExpectedPrepare)

func TestFindConnection(t *testing.T) {
  connection := entity.Connection{
    1,
    2,
    "approved",
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  rows := sqlmock.NewRows([]string{
    "friend_user_id", "user_id", "connection_status", "user_id_equals_friend_id"}).
      AddRow(1, 2, "approved", "")

  query := `SELECT (.+) FROM friend_list WHERE friend_user_id =(.+) AND user_id =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectQuery().WithArgs(connection.FriendUserId, connection.UserId).WillReturnRows(rows)

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  foundConnection, err := friendRepo.FindConnection(&connection)
  assert.NoError(t, err)
  assert.NotNil(t, foundConnection)
}

func TestInsertConnection(t *testing.T) {
  connection := entity.Connection{
    1,
    2,
    "approved",
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `INSERT INTO friend_list \(friend_user_id, user_id, connection_status\) VALUES \((.+),(.+),(.+)\)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectExec().WithArgs(connection.FriendUserId, connection.UserId, connection.Status).
    WillReturnResult(sqlmock.NewResult(1, 1))

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  err = friendRepo.InsertConnection(&connection)
  assert.NoError(t, err)
}

func TestUpdateConnection(t *testing.T) {
  connection := entity.Connection{
    1,
    2,
    "approved",
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `UPDATE friend_list SET connection_status =(.+) WHERE friend_user_id =(.+) AND user_id =(.+)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectExec().WithArgs(connection.Status, connection.FriendUserId, connection.UserId).
    WillReturnResult(sqlmock.NewResult(1, 1))

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  err = friendRepo.UpdateStatus(&connection)
  assert.NoError(t, err)
}

func TestDeleteConnection(t *testing.T) {
  connection := entity.Connection{
    1,
    2,
    "approved",
  }

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `DELETE FROM friend_list WHERE friend_user_id IN \((.+),(.+)\) AND user_id IN \((.+),(.+)\)`

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectExec().WithArgs(connection.FriendUserId, connection.UserId,
    connection.FriendUserId, connection.UserId).
  WillReturnResult(sqlmock.NewResult(1, 1))

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  err = friendRepo.Delete(&connection)
  assert.NoError(t, err)
}

func TestFindFriendsByUserId(t *testing.T) {

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `SELECT users_data.* FROM users_data JOIN friend_list ON \(users_data.user_id = friend_list.user_id\)
   WHERE friend_list.friend_user_id =(.+) AND friend_list.connection_status='approved'`

  rows := sqlmock.NewRows([]string{
    "user_id", "email", "first_name", "last_name", "phone", "role_in_network", "account_status", "avatar_ref"}).
    AddRow(2, "email@gmail.com", "Name", "Surname", "phone", "user", "active", "")

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectQuery().WithArgs(1).WillReturnRows(rows)

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  _, err = friendRepo.FriendsByUserID(1)
  assert.NoError(t, err)
}

func TestUsersWithRequestsTo(t *testing.T) {

  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  query := `SELECT users_data.* FROM users_data JOIN friend_list ON \(users_data.user_id = friend_list.user_id\)
   WHERE friend_list.friend_user_id =(.+) AND friend_list.connection_status='waiting'`

  rows := sqlmock.NewRows([]string{
    "user_id", "email", "first_name", "last_name", "phone", "role_in_network", "account_status", "avatar_ref"}).
    AddRow(2, "email@gmail.com", "Name", "Surname", "phone", "user", "active", "")

  var prep *sqlmock.ExpectedPrepare
  for _, query := range friendQueriesRegexes {
    prep = mock.ExpectPrepare(query)
    friendPreps[query] = prep
  }

  friendPreps[query].ExpectQuery().WithArgs(1).WillReturnRows(rows)

  friendRepo, err := newMySqlFriendshipDao(db)
  if err != nil {
    t.Fatalf("an error:\n'%s'\nwas not expected when opening a stub database connection", err)
  }

  _, err = friendRepo.UsersWithRequestsTo(1)
  assert.NoError(t, err)
}

