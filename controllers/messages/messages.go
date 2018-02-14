package messages

import (
	"go-team-room/conf"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var logRus = conf.GetLog()

func init() {
	//logRus = logRus.WithField(logRus.Fields{"packet": "messages"})

	logRus.Debug("Messages packet initialized Ok")
}
