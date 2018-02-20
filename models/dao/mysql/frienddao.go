package mysql

import (
  "database/sql"
  "go-team-room/models/dao/interfaces"
  "fmt"
  "go-team-room/models/dao/entity"
)

type friendDaoImpl struct {
  conn *sql.DB

  findConnection   *sql.Stmt
  insertConnection *sql.Stmt
  delete           *sql.Stmt
  updateStatus     *sql.Stmt
  friendsById      *sql.Stmt
  requestsToId     *sql.Stmt
}

var _ interfaces.FriendDao = &friendDaoImpl{}

func newMySqlFriendshipDao(conn *sql.DB) (interfaces.FriendDao, error) {

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &friendDaoImpl{
    conn: conn,
  }

  var err error
  if db.findConnection, err = conn.Prepare(findConnectionStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.insertConnection, err = conn.Prepare(insertConnectionStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.updateStatus, err = conn.Prepare(updateConnectionStatus); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.delete, err = conn.Prepare(deleteFriendship); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.friendsById, err = conn.Prepare(friendsByUserIdStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }
  if db.requestsToId, err = conn.Prepare(usersWithRequestsToIdStatement); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }

  return db, nil
}

const findConnectionStatement = `SELECT * FROM friend_list WHERE friend_user_id = ? AND user_id = ?`

func (f *friendDaoImpl) FindConnection(connection *entity.Connection) (entity.Connection, error) {
  foundConnection, err := scanConnection(f.findConnection.QueryRow(connection.FriendUserId, connection.UserId))
  if err != nil {
    return foundConnection, err
  }

  return foundConnection, nil
}

const friendsByUserIdStatement =
  `SELECT users_data.* FROM users_data JOIN friend_list
   ON (users_data.user_id = friend_list.user_id)
   WHERE friend_list.friend_user_id = ?
--     AND users_data.account_status <> 'deleted'
    AND friend_list.connection_status='approved'`

func (f friendDaoImpl) FriendsByUserID(id int64) ([]entity.User, error) {
  return getUserConnections(id, f.friendsById)
}

const usersWithRequestsToIdStatement =
  `SELECT users_data.* FROM users_data JOIN friend_list
   ON (users_data.user_id = friend_list.user_id)
   WHERE friend_list.friend_user_id = ?
--     AND users_data.account_status <> 'deleted'
    AND friend_list.connection_status='waiting'`


func (f *friendDaoImpl) UsersWithRequestsTo(id int64) ([]entity.User, error) {
  return getUserConnections(id, f.requestsToId)
}

func getUserConnections(id int64, stmt *sql.Stmt) ([]entity.User, error) {
  rows, err := stmt.Query(id)

  if err != nil {
    return nil, err
  }
  defer rows.Close()

  friends := []entity.User{}
  var user entity.User
  for rows.Next() {
    user, err = scanUser(rows)
    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }

    friends = append(friends, user)
  }

  return friends, nil
}

const insertConnectionStatement =
  `INSERT INTO friend_list (friend_user_id, user_id, connection_status) VALUES (?, ?, ?)`

func (f *friendDaoImpl) InsertConnection(connection *entity.Connection) error {
  _, err := execAffectingOneRow(f.insertConnection, connection.FriendUserId, connection.UserId, connection.Status)

  return err
}

const updateConnectionStatus = `UPDATE friend_list SET connection_status = ?
                                WHERE friend_user_id = ? AND user_id = ?`

func (f *friendDaoImpl) UpdateStatus(connection *entity.Connection) error {
  _, err := execAffectingOneRow(f.updateStatus, connection.Status, connection.FriendUserId, connection.UserId)

  return err
}

const deleteFriendship = `DELETE FROM friend_list WHERE friend_user_id IN (?, ?) AND user_id IN (?, ?)`

func (f *friendDaoImpl) Delete(connection *entity.Connection) error {
  _, err := f.delete.Exec(
    connection.FriendUserId,
    connection.UserId,
    connection.FriendUserId,
    connection.UserId)

  return err
}

var (
  friend_user_id  int64
  u_id            int64
  status          sql.NullString
  constraint      sql.NullString
)

func scanConnection(s rowScanner) (entity.Connection, error) {

  connection := entity.Connection{}

  if err := s.Scan(&friend_user_id, &u_id, &status, &constraint); err != nil {
    return connection, err
  }

  connection.FriendUserId = friend_user_id
  connection.UserId = u_id
  connection.Status = entity.ConnectionStatus(status.String)
  return connection, nil
}
