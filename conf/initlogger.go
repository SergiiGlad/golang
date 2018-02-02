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

// create global logger's instance
var log *logrus.Logger

// A function to fetch current logger
func GetLog() *logrus.Logger {
	return log
}

// Init logger
func init() {

	// Creates a new logger. Configuration should be set by changing `Formatter`,
	// `Out` and `Hooks` directly on the default logger instance. You can also just
	// instantiate your own:
	//
	//    var log = &Logger{
	//      Out: os.Stderr,
	//      Formatter: new(JSONFormatter),
	//      Hooks: make(LevelHooks),
	//      Level: logrus.DebugLevel,
	//    }
	//
	// It's recommended to make this a global instance called `log`.
	log = logrus.New()

	// read configuration from conf/conf.json and setup logger
	readSetupLogger(log)

	// end logger setup configuration
	log.Info("Logger configuration has finished")

}

func readSetupLogger(log *logrus.Logger) {

	// new hook, you just need a registration
	// beacause logrus_mate doesn't have a lsfhook in package
	// func newHook()
	logrus_mate.RegisterHook("lfshook", newHook)

	// Read and unmarshal configuration from viper
	mate_conf := logrusHelper.UnmarshalConfiguration(viper.GetViper())

	// apply the configuration to logger
	if err := logrusHelper.SetConfig(log, mate_conf); err != nil {
		// Handle errors reading the mate config
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

// If you want to use your own hook, you just need todo as follow
func newHook(options logrus_mate.Options) (hook logrus.Hook, err error) {

	// name of field viper file conf.json
	logfile := "logfile"
	rotatehours := "rotatehours"
	maxdays := "maxdays"

	filename, err := options.String(logfile)
	if err != nil {
		filename = viper.GetString("work_dir") + "/logs/current.log"
		logrus.Info("Useing default name for logs file current.log")
	} else {
		filename = viper.GetString("work_dir") + filename
	}

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
