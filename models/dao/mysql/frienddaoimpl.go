package mysql

import (
  "database/sql"
  "go-team-room/models/dao/interfaces"
  "fmt"
  "go-team-room/models/dao/entity"
)

type friendshipDaoImpl struct {
  conn *sql.DB

  friendsById     *sql.Stmt

}

var _ interfaces.FriendshipDao = &friendshipDaoImpl{}

func newMySqlFriendshipDao(conn *sql.DB) (interfaces.FriendshipDao, error) {

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  db := &friendshipDaoImpl{
    conn: conn,
  }

  var err error

  if db.friendsById, err = conn.Prepare(findFriendsByUserId); err != nil {
    return nil, fmt.Errorf("mysql: prepare list: %v", err)
  }

  return db, nil
}

const findFriendsByUserId = `SELECT * FROM friend_list WHERE user_id = ? AND connection_status='approved'`

func (f *friendshipDaoImpl) FriendsByUserID(id int64) ([]entity.Friendship, error) {
  rows, err := f.friendsById.Query()

  if err != nil {
    return nil, err
  }
  defer rows.Close()

  friendships := []entity.Friendship{}
  var friendship *entity.Friendship
  for rows.Next() {
    friendship, err = scanFriends(rows)
    if err != nil {
      return nil, fmt.Errorf("mysql: could not read row: %v", err)
    }

    friendships = append(friendships, *friendship)
  }

  return friendships, nil
}

var (
  friend_user_id  int64
  u_id            int64
  status          sql.NullString
)

func scanFriends(s rowScanner) (*entity.Friendship, error) {

  if err := s.Scan(&friend_user_id, &u_id, &status); err != nil {
    return nil, err
  }

  return &entity.Friendship{
    friend_user_id,
    u_id,
    entity.ConnectionStatus(status.String),
  }, nil
}
