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
	readSetupLogger()

	// end logger setup configuration
	log.Info("Logger configuration has finished")

	// test messages
	log.Debug("Debug level message from gohum ")        //"Useful debugging information."
	log.Info("Info level message from gohum")           //"Something noteworthy happened!"
	log.Warn("Warn level message from gohum")           //You should probably take a look at this."
	log.Error("Something failed but I'm not quitting.") //Something failed but I'm not quitting.

}

func readSetupLogger() {

	// new hook, you just need a registration
	// beacause logrus_mate doesn't have a lsfhook in package
	// func newHook func(logrus_mate.Options) (hook logrus.Hook, err error)
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

	logHook := log.WithFields(logrus.Fields{
		"Setup Hook": "lfshoook",
		"Useing":     "default values"})

	// names of fields which will be parsed from viper file conf.json
	// section "hooks"

	formatter := "formatter"
	logfile := "logfile"
	rotatehours := "rotatehours"
	maxdays := "maxdays"

	//default values for settings
	filenameDeafult := "/logs/current.log"
	rotationTimeDefault := "24" // 24 hours
	maxDaysDefault := "7"
	formatterDefault := "text"

	filename, err := options.String(logfile)
	if err != nil {
		logHook.Error(err)
		filename = filenameDeafult
	}

	//Interval between file rotation.
	//By default logs are rotated every 86400 seconds.
	//Note: Remember to use time.Duration values.
	rotationTime, err := options.String(rotatehours)
	if err != nil {
		logHook.Error(err)
		rotationTime = rotationTimeDefault
	}
	rotation, err := strconv.ParseInt(rotationTime, 10, 64)

	maxAge, err := options.String(maxdays)
	if err != nil {
		logHook.Error(err)
		maxAge = maxDaysDefault
	}
	age, err := strconv.ParseInt(maxAge, 10, 64)

	writer, err := rotatelogs.New(
		viper.GetString("work_dir")+filename+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(time.Duration(86400*age)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(3600*rotation)*time.Second),
		rotatelogs.WithClock(rotatelogs.UTC),
	)

	pathMap := lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}

	formatterHook, err := options.String(formatter)
	if err != nil {
		logHook.Error(err)
		formatterHook = formatterDefault
	}

	if err != nil {
		logHook.Warn("Useing defaul values!")
	}

	switch formatterHook {
	case "json":
		hook = lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		)

	default:
		hook = lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{},
		)
	}

	if err = options.ToObject(&hook); err != nil {
		return
	}

	return
}
