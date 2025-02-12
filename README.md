# Logger

### Installation:
```sh
go get github.com/gaigals/logger@latest
```

### Usage:

```go
appName := "myApp"
enableSyslog := true
syslogFlags := syslog.LOG_INFO|syslog.LOG_USER
logFilePath := "path/to/log/file.log" // Optional.

myLogger, err := logger.NewLogger(
    appName,
    logFilePath,
    enableSyslog,
    syslogFlags,
)
if err != nil {
    log.Fatalln(err)
}

myLogger.Info("test message", 1)

// Or setup it globaly.
err = logger.NewGlobalLogger(
    appName,
    logFilePath, // Optional.
    enableSyslog,
    syslogFlags,
)
if err != nil {
    log.Fatalln(err)
}

logger.Infof("test message %d", 2)


// Or you can make logger without error checks using
// These function calls will exit with exit code 1 in case if logger creation
// fails.
myOtherLogger := logger.NewLoggerOrFatal(
    appName,
    logFilePath,
    enableSyslog,
    syslogFlags,
)

logger.NewGlobalLoggerOrFatal(
    appName,
    logFilePath,
    enableSyslog,
    syslogFlags,
)
```

If you don't want to use file as log output, then set logFilePath as empty
string.\
Note, on docker syslog does not work as you cannot make connection to a local
syslog server.
