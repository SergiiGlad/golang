// Package logs provide a general log settings
// It defines a type logrus Logger, with methods for formatting output.
// Let see ./conf/sample_log.go for example ho uses logrus Logger

package conf

import (
	"fmt"
	"github.com/heirko/go-contrib/logrusHelper"
	"github.com/heralight/logrus_mate"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

// The variable log is a name for our logger
var log *logrus.Logger

// A function to fetch current logger
func GetLog() *logrus.Logger {
	return log
}

// Init logger
func init() {

	// new hook, you just need a registration
	// beacause logrus_mate doesn't have a lsfhook in package
	logrus_mate.RegisterHook("lfshook", NewHook)

	// Read and unmarshal configuration from viper

	var mate_conf logrus_mate.LoggerConfig
	mate_conf = logrusHelper.UnmarshalConfiguration(viper.GetViper())

	// create logger's instance
	log = logrus.New()

	// apply the configuration to logger
	if logrusHelper.SetConfig(log, mate_conf); err != nil {
		// Handle errors reading the mate config
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// end logger setup configuration

	log.Info("Logger configuration has finished")

}

// If you want to use your own hook, you just need todo as follow
func NewHook(options logrus_mate.Options) (hook logrus.Hook, err error) {

	// name of field viper file conf.json
	logfile := "logfile"
	rotatehours := "rotatehours"
	maxdays := "maxdays"

	filename, err := options.String(logfile)

	//Interval between file rotation.
	//By default logs are rotated every 86400 seconds.
	//Note: Remember to use time.Duration values.
	rotationTime, err := options.String(rotatehours)
	maxAge, err := options.String(maxdays)

	rotation, err := strconv.ParseInt(rotationTime, 10, 64)
	age, err := strconv.ParseInt(maxAge, 10, 64)

	writer, err := rotatelogs.New(
		filename+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Duration(86400*age)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(3600*rotation)*time.Second),
		rotatelogs.WithClock(rotatelogs.UTC),
	)

	pathMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.ErrorLevel: writer,
	}

	hook = lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	)

	if err = options.ToObject(&hook); err != nil {
		return
	}

	return
}
