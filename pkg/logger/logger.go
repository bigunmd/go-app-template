package logger

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoggerConfigSetDefault() {
	_ = viper.BindEnv("logger.level", "LOGGER_LEVEL", "LOG_LEVEL", "LOGGING")
	viper.SetDefault("logger.level", "info")

	_ = viper.BindEnv("logger.writeToFile", "LOGGER_WRITE_TO_FILE", "LOG_WRITE_TO_FILE")
	viper.SetDefault("logger.writeToFile", false)

	_ = viper.BindEnv("logger.file.path", "LOGGER_FILE_PATH", "LOG_FILE_PATH")
	viper.SetDefault("logger.file.path", "~/")

	_ = viper.BindEnv("logger.file.name", "LOGGER_FILE_NAME", "LOG_FILE_NAME")
	viper.SetDefault("logger.file.name", "app.log")

	_ = viper.BindEnv("logger.file.maxAge", "LOGGER_FILE_MAX_AGE", "LOG_FILE_MAX_AGE")
	viper.SetDefault("logger.file.maxAge", 24*time.Hour)

	_ = viper.BindEnv("logger.file.rotationTime", "LOGGER_FILE_ROTATION_TIME", "LOG_FILE_ROTATION_TIME")
	viper.SetDefault("logger.file.rotationTime", 24*7*time.Hour)

}

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	AddFields(fields map[string]interface{}) Entry
	AddField(key string, value interface{}) Entry
	AddError(err error) Entry

	GetLogLevel() string
	SetLogLevel(level string) error
}

type Entry interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
}

type LogrusLogger struct {
	*logrus.Logger
}

func (l *LogrusLogger) AddFields(fields map[string]interface{}) Entry {
	return l.WithFields(fields)
}

func (l *LogrusLogger) AddField(key string, value interface{}) Entry {
	return l.WithField(key, value)
}

func (l *LogrusLogger) GetLogLevel() string {
	return l.GetLevel().String()
}

func (l *LogrusLogger) SetLogLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	l.SetLevel(lvl)
	return nil
}

// AddError implements Logger
func (l *LogrusLogger) AddError(err error) Entry {
	return l.WithError(err)
}

func NewLogger() Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lvl)
	logger.SetReportCaller(true)
	jsonFmt := &logrus.JSONFormatter{PrettyPrint: true}
	if viper.GetBool("logger.writeToFile") {
		jsonFmt.PrettyPrint = false
		writer, err := rotatelogs.New(
			viper.GetString("logger.file.path")+viper.GetString("logger.file.name")+".%Y%m%dT%H%M",
			rotatelogs.WithLinkName(viper.GetString("logger.file.path")+viper.GetString("logger.file.name")),
			rotatelogs.WithMaxAge(viper.GetDuration("logger.file.maxAge")),
			rotatelogs.WithRotationTime(viper.GetDuration("logger.file.rotationTime")),
		)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, writer)
		logger.SetOutput(mw)
	}
	logger.SetFormatter(jsonFmt)
	return &LogrusLogger{logger}
}
