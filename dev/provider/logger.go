package provider

import (
	"excel-receiver/config"
	"excel-receiver/provider/dailylogger"
	"excel-receiver/util"
	"fmt"
	"io"
	"path"

	"github.com/sirupsen/logrus"
)

type LogType int

const (
	AppLog LogType = 0
	AmqLog LogType = 1
	DBLog  LogType = 2
)

type ILogger interface {
	Infof(logType LogType, format string, args ...interface{})
	Errorf(logType LogType, format string, args ...interface{})
	Debugf(logType LogType, format string, args ...interface{})
	WithFields(logType LogType, fields logrus.Fields) *logrus.Entry
	InfoWithFields(logType LogType, fields map[string]interface{}, format string, args ...interface{})
	ErrorWithFields(logType LogType, fields map[string]interface{}, format string, args ...interface{})
}

type logrusLogger struct {
	loggers map[LogType]*logrus.Logger
}

func NewLogger() *logrusLogger {
	return &logrusLogger{
		loggers: map[LogType]*logrus.Logger{
			AppLog: createLogger("info", "error"),
			DBLog:  createLogger("db", "dberror"),
		},
	}
}

func createLogger(infoDir, errorDir string) *logrus.Logger {
	infoLogFile := path.Join(config.Configuration.Logger.Dir, infoDir, fmt.Sprintf("%s.info.log", config.Configuration.Logger.FileName))
	errorLogFile := path.Join(config.Configuration.Logger.Dir, errorDir, fmt.Sprintf("%s.error.log", config.Configuration.Logger.FileName))

	logger := logrus.New()

	maxAge := config.Configuration.Logger.MaxAge
	maxBackups := config.Configuration.Logger.MaxBackups
	maxSize := config.Configuration.Logger.MaxSize
	compress := config.Configuration.Logger.Compress
	localTime := config.Configuration.Logger.LocalTime

	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@time", // to move time in log to left
		},
	}

	logger.SetFormatter(formatter)

	logger.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(infoLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	logger.AddHook(&WriterHook{
		Writer: dailylogger.NewDailyRotateLogger(errorLogFile, maxSize, maxBackups, maxAge, localTime, compress),
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	return logger
}

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func (l *logrusLogger) Infof(logType LogType, format string, args ...interface{}) {
	l.getLog(logType).Infof(format, args...)
}

func (l *logrusLogger) Errorf(logType LogType, format string, args ...interface{}) {
	l.getLog(logType).Errorf(format, args...)
}

func (l *logrusLogger) Debugf(logType LogType, format string, args ...interface{}) {
	l.getLog(logType).Debugf(format, args...)
}

func (l *logrusLogger) WithFields(logType LogType, fields logrus.Fields) *logrus.Entry {
	logger := l.getLog(logType)
	return logger.WithFields(fields)
}

func (l *logrusLogger) InfoWithFields(logType LogType, fields map[string]interface{}, format string, args ...interface{}) {
	l.getLog(logType).WithFields(fields).Infof(format, args...)
}

func (l *logrusLogger) ErrorWithFields(logType LogType, fields map[string]interface{}, format string, args ...interface{}) {
	l.getLog(logType).WithFields(fields).Errorf(format, args...)
}

func (l logrusLogger) getLog(logType LogType) *logrus.Logger {
	if v, ok := l.loggers[logType]; ok {
		return v
	}

	return l.loggers[AppLog]
}

func InitLogDir() {
	workingDirectory := config.Configuration.Logger.Dir
	logRootDirectory := path.Join(workingDirectory)

	logDirectories := []string{
		logRootDirectory,
		path.Join(logRootDirectory, "info"),
		path.Join(logRootDirectory, "error"),
		path.Join(logRootDirectory, "db"),
		path.Join(logRootDirectory, "dberror"),
	}

	if err := util.CreateDirectory(logDirectories...); err != nil {
		panic(err)
	}
}
