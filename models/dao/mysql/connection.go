package mysql

import (
  "log"
  "go-team-room/models/dao/interfaces"
)

var (
  DB interfaces.Dal
)

func init() {

  var err error

  DB, err = newMySQLDatabase()

  if err != nil {
    log.Fatalf("Could not connect DB: %s", err)
  }
}
