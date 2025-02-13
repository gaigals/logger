package logger

import (
	"errors"
	"fmt"
	"log/syslog"
	"os"
	"path/filepath"
	"time"
)

// ErrSyslogConnFailed gets triggered if syslog setup failed.
var ErrSyslogConnFailed = errors.New("new syslog setup error")

var logger Logger

type Logger struct {
	syslogWriter        *syslog.Writer
	fileWriter          *os.File
	formatter           func(level syslog.Priority, msg string, args ...any) string
	useStdBackupWritter bool // Use std log package to write logs.
}

// Printf without applied formatters.
func (logger *Logger) Printf(format string, args ...any) {
	logger.printLogf(false, syslog.LOG_INFO, os.Stdout, format, args...)
}

// Println withotu applied formatters.
func (logger *Logger) Println(s ...any) {
	logger.printLog(false, syslog.LOG_INFO, os.Stdout, s...)
}

// Info ...
func (logger *Logger) Info(s ...any) {
	logger.printLog(true, syslog.LOG_INFO, os.Stdout, s...)
}

// Infof
func (logger *Logger) Infof(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_INFO, os.Stdout, format, args...)
}

// Debug ...
func (logger *Logger) Debug(s ...any) {
	logger.printLog(true, syslog.LOG_DEBUG, os.Stdout, s...)
}

// Debugf ...
func (logger *Logger) Debugf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_DEBUG, os.Stdout, format, args...)
}

// Warn ...
func (logger *Logger) Warn(s ...any) {
	logger.printLog(true, syslog.LOG_WARNING, os.Stdout, s...)
}

// Warnf ...
func (logger *Logger) Warnf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_WARNING, os.Stdout, format, args...)
}

// Alert ...
func (logger *Logger) Alert(s ...any) {
	logger.printLog(true, syslog.LOG_ALERT, os.Stdout, s...)
}

// Alertf ...
func (logger *Logger) Alertf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_ALERT, os.Stdout, format, args...)
}

// Notice ...
func (logger *Logger) Notice(s ...any) {
	logger.printLog(true, syslog.LOG_NOTICE, os.Stdout, s...)
}

// Noticef ...
func (logger *Logger) Noticef(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_NOTICE, os.Stdout, format, args...)
}

// Error ...
func (logger *Logger) Error(s ...any) {
	logger.printLog(true, syslog.LOG_ERR, os.Stderr, s...)
}

// Errorf ...
func (logger *Logger) Errorf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_ERR, os.Stderr, format, args...)
}

// Critical ...
func (logger *Logger) Critical(s ...any) {
	logger.printLog(true, syslog.LOG_CRIT, os.Stderr, s...)
}

// Criticalf ...
func (logger *Logger) Criticalf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_CRIT, os.Stderr, format, args...)
}

// Emergency ...
func (logger *Logger) Emergency(s ...any) {
	logger.printLog(true, syslog.LOG_EMERG, os.Stderr, s...)
}

// Emergencyf ...
func (logger *Logger) Emergancyf(format string, args ...any) {
	logger.printLogf(true, syslog.LOG_EMERG, os.Stderr, format, args...)
}

func (logger *Logger) openLogFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	err := checkFileDir(filePath)
	if err != nil {
		return err
	}

	logger.fileWriter, err = os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("log file=%s open error: %w", filePath, err)
	}

	return nil
}

func (logger *Logger) printLog(
	applyFormating bool,
	level syslog.Priority,
	osFile *os.File,
	s ...any,
) {
	msg := fmt.Sprint(s...)

	if applyFormating {
		msg = logger.formatter(level, msg)
	}

	if logger.fileWriter != nil {
		fmt.Fprintln(logger.fileWriter, msg)
	}

	fmt.Fprintln(osFile, msg)
	logger.writeToSyslog(level, msg)
}

func (logger *Logger) printLogf(
	applyFormating bool,
	level syslog.Priority,
	osFile *os.File,
	msg string,
	args ...any,
) {
	if applyFormating {
		msg = logger.formatter(level, msg, args...)
	} else {
		msg = fmt.Sprintf(msg, args...)
	}

	if logger.fileWriter != nil {
		fmt.Fprintln(logger.fileWriter, msg)
	}

	fmt.Fprintln(osFile, msg)
	logger.writeToSyslog(level, msg)
}

func (logger *Logger) writeToSyslog(level syslog.Priority, msg string) {
	if logger.syslogWriter == nil {
		return
	}

	switch level {
	case syslog.LOG_DEBUG:
		logger.syslogWriter.Debug(msg)
	case syslog.LOG_INFO:
		logger.syslogWriter.Info(msg)
	case syslog.LOG_WARNING:
		logger.syslogWriter.Warning(msg)
	case syslog.LOG_NOTICE:
		logger.syslogWriter.Notice(msg)
	case syslog.LOG_ALERT:
		logger.syslogWriter.Alert(msg)
	case syslog.LOG_ERR:
		logger.syslogWriter.Err(msg)
	case syslog.LOG_CRIT:
		logger.syslogWriter.Crit(msg)
	case syslog.LOG_EMERG:
		logger.syslogWriter.Emerg(msg)
	default:
		logger.syslogWriter.Info(msg)
	}
}

func checkFileDir(filePath string) error {
	parentDir := filepath.Dir(filePath)

	info, err := os.Stat(parentDir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("the path=%s exists but is not directory",
				parentDir)
		}

		return nil
	}

	if os.IsNotExist(err) {
		return fmt.Errorf("log file dir=%s does not exist", parentDir)
	}

	return fmt.Errorf("log file dir=%s stat error: %w", parentDir, err)
}

// NewLoggerOrFatal creates new logger or exits with status code 1 if logger
// creation failed.
func NewLoggerOrFatal(
	appName, logFilePath string,
	enableSyslog bool,
	flags syslog.Priority,
) *Logger {
	l, err := NewLogger(appName, logFilePath, enableSyslog, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return l
}

// NewGlobalLoggerOrFatal creates new global logger or exits with status code
// 1 if logger creation failed.
func NewGlobalLoggerOrFatal(
	appName, logFilePath string,
	enableSyslog bool,
	flags syslog.Priority,
) {
	l, err := newLogger(appName, logFilePath, enableSyslog, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger = *l
}

// NewLogger ...
func NewLogger(
	appName, logFilePath string,
	enableSyslog bool,
	flags syslog.Priority,
) (*Logger, error) {
	l, err := newLogger(appName, logFilePath, enableSyslog, flags)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// NewGlobalLogger ...
func NewGlobalLogger(appName, logFilePath string, enableSyslog bool, flags syslog.Priority) error {
	l, err := newLogger(appName, logFilePath, enableSyslog, flags)
	if err != nil {
		return err
	}

	logger = *l
	return nil
}

func newLogger(
	appName, logFilePath string,
	enableSyslog bool,
	flags syslog.Priority,
) (*Logger, error) {
	var err error
	var syslogWriter *syslog.Writer

	if enableSyslog {
		syslogWriter, err = newSysLogger(appName, flags)
		if err != nil {
			return nil, err
		}
	}

	l := Logger{
		syslogWriter: syslogWriter,
		formatter:    formatLog,
	}

	err = l.openLogFile(logFilePath)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func newSysLogger(appName string, flags syslog.Priority) (*syslog.Writer, error) {
	slog, err := newSyslog(appName, flags)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSyslogConnFailed, err)
	}

	return slog, nil
}

// Printf without applied formatters.
func Printf(format string, args ...any) {
	logger.printLogf(false, syslog.LOG_INFO, os.Stdout, format, args...)
}

// Println withotu applied formatters.
func Println(s ...any) {
	logger.printLog(false, syslog.LOG_INFO, os.Stdout, s...)
}

// Info ...
func Info(s ...any) {
	logger.Info(s...)
}

// Infof
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Debug ...
func Debug(s ...any) {
	logger.Debug(s...)
}

// Debugf ...
func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

// Warn ...
func Warn(s ...any) {
	logger.Warn(s...)
}

// Warnf ...
func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Alert ...
func Alert(s ...any) {
	logger.Alert(s...)
}

// Alertf ...
func Alertf(format string, args ...any) {
	logger.Alertf(format, args...)
}

// Notice ...
func Notice(s ...any) {
	logger.Notice(s...)
}

// Noticef ...
func Noticef(format string, args ...any) {
	logger.Noticef(format, args...)
}

// Error ...
func Error(s ...any) {
	logger.Error(s...)
}

// Errorf ...
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Critical ...
func Critical(s ...any) {
	logger.Critical(s...)
}

// Criticalf ...
func Criticalf(format string, args ...any) {
	logger.Criticalf(format, args...)
}

// Emergency ...
func Emergency(s ...any) {
	logger.Emergency(s...)
}

// Emergencyf ...
func Emergancyf(format string, args ...any) {
	logger.Emergancyf(format, args...)
}

func newSyslog(appName string, flags syslog.Priority) (*syslog.Writer, error) {
	if flags == 0 {
		flags = syslog.LOG_INFO | syslog.LOG_SYSLOG
	}

	return syslog.New(flags, appName)
}

func formatLog(level syslog.Priority, msg string, args ...any) string {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return fmt.Sprintf(
		"%-7s | %s | %s | %s",
		LogLevelToString(level),
		time.Now().Format("MST"),
		time.Now().Format("02/01/2006 15:04:05.000"),
		msg,
	)
}

// LogLevelToString converts syslog priority (log level) as readable string.
func LogLevelToString(level syslog.Priority) string {
	switch level {
	case syslog.LOG_EMERG:
		return "EMERG"
	case syslog.LOG_ALERT:
		return "ALERT"
	case syslog.LOG_CRIT:
		return "CRIT"
	case syslog.LOG_ERR:
		return "ERROR"
	case syslog.LOG_WARNING:
		return "WARNING"
	case syslog.LOG_NOTICE:
		return "NOTICE"
	case syslog.LOG_INFO:
		return "INFO"
	case syslog.LOG_DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
