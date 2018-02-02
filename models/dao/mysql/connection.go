package mysql

import (
  "log"
  "go-team-room/models/dao/interfaces"
)

//this file defines global mysql db connection that can be used and accessed in any project part

var (
  DB interfaces.MySqlDal
)

func init() {

  var err error

  DB, err = newMySQLDatabase()

  if err != nil {
    log.Fatalf("Could not connect DB: %s", err)
  }
}
