package mysql

import (
  "log"
  "go-team-room/models/dao/interfaces"
)

var (
  DB interfaces.Database
)

func init() {

  var err error

  DB, err = newMySQLDatabase()

  if err != nil {
    log.Fatal("Could not connect DB.")
  }
}
