package mysql

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
  "fmt"
  "go-team-room/conf"
)

func get() *sql.DB {
  dsn := fmt.Sprintf("%s", conf.MysqlDsn)
  db, err := sql.Open("mysql", dsn)

  if err != nil {
    log.Fatalln(err)
    return nil
  }

  return db
}
