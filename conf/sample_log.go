// File sample_log demonstares a sample logging by useing logrus package
// The simplest way to use Logrus is simply the package-level exported logger
// logrus.Info() and customize it all you want
// You also can create new loggger with specific parameters 
// For a full guide visit https://github.com/sirupsen/logrus



package conf

import (
	"github.com/sirupsen/logrus"

)

// print log message str with *logrus.Logger
func logout(logrus *logrus.Logger, str string) {

	println("===", str, "===")
	logrus.Debug("Debug level message from gohum ")        //"Useful debugging information."
	logrus.Info("Info level message from gohum")           //"Something noteworthy happened!"
	logrus.Warn("Warn level message from gohum")           //You should probably take a look at this."
	logrus.Error("Something failed but I'm not quitting.") //Something failed but I'm not quitting.

	return

	// after call program will exit
	logrus.Fatal("Bye.")         // Calls os.Exit(1) after logging
	logrus.Panic("I'm bailing.") // Calls panic() after logging
	println("===")
}

func Sample_log() {

	// With the default ASCII formatter logrus.SetFormatter(&logrus.TextFormatter{})
	logout(logrus.New(), "default settings")

	//Create new logger
	var log = logrus.New()

	//With json formatter, for easy parsing by logstash or Splunk
	log.Formatter = new(logrus.JSONFormatter) // default
	logout(log, "json formatter")

	// This works only on your instance of logger, obtained with `logrus.New()`.
	// PanicLevel
	// FatalLevel
	// ErrorLevel
	// WarnLevel
	// InfoLevel
	// DebugLevel
	// You can change DebugLevel to try output
	//log the warn level severity or above
	log.SetLevel(logrus.WarnLevel)
	logout(log, "warn level")

	// DisableTimestamp allows disabling automatic timestamps in output
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter.(*logrus.JSONFormatter).DisableTimestamp = true // remove timestamp from test output
	logout(log, "DisableTimestamp")

	// You should log the much more discoverable with added fields
	log.WithFields(logrus.Fields{
		"number": 1,
		"size":   10,
	}).Error("You can produce much more useful logging messages")

	// Often it's helpful to have fields always attached to log statements
	// You should create logger with fileds
	requestLogger := log.WithFields(logrus.Fields{"request_id": 1, "user_ip": "127.0.0.1"})
	// will log request_id and user_ip
	requestLogger.Warn("something not great happened")

}
