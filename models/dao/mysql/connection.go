package mysql

import (
	"go-team-room/conf"
	"go-team-room/models/dao/interfaces"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

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
