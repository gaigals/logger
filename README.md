# Logger

### Installation:
```sh
go get github.com/gaigals/logger@latest
```

### Usage:

```go
appName := "myApp"
enableSyslog := true
disableOsOutput := false // Enable or disable stdout/stderr
syslogFlags := syslog.LOG_USER
logFilePath := "path/to/log/file.log" // Optional.

myLogger, err := logger.NewLogger(
    appName,
    logFilePath,
    enableSyslog,
    disableOsOutput,
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
    disableOsOutput,
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
    disableOsOutput,
    syslogFlags,
)

logger.NewGlobalLoggerOrFatal(
    appName,
    logFilePath,
    enableSyslog,
    disableOsOutput,
    syslogFlags,
)
```

If you don't want to use file as log output, then set logFilePath as empty
string.\
Note, on docker syslog does not work as you cannot make connection to a local
syslog server.


### Optional syslog example

For cases when you are developing in docker but the final product is running
on non-dockerized enviroment.

```go
func newLogger() (*logger.Logger, error) {
    appName := "myApp"
    enableSyslog := true
    disableOsOutput := false // Enable or disable stdout/stderr
    syslogFlags := syslog.LOG_USER
    logFilePath := "path/to/log/file.log" // Optional.

    myLogger, err := logger.NewLogger(
        appName,
        logFilePath,
        enableSyslog,
        disableOsOutput,
        syslogFlags,
    )
    // No errors, syslog enabled.
    if err == nil {
        return myLogger, nil
    }
    // Unexpected error handling.
    if !errors.Is(err, logger.ErrSyslogConnFailed) {
        return nil, fmt.Errorf("unexpected new logger error: %w", err)
    }

    // Syslog connection failed (for example, docker enviroment).
    // Enable logger without syslog connection.
    myLogger, err := logger.NewLogger(
        appName,
        logFilePath,
        false, // Do not log into a syslog.
        disableOsOutput,
        syslogFlags,
    )
    if err != nil {
        return nil, fmt.Errorf(
            "unexpected new logger without syslog error: %w",
            err,
        )
    }

    return myLogger, nil
}

func main() {
    myLogger, err := newLogger()
    if err != nil {
        log.Fatalln(err)
    }

    myLogger.Info("test message", 1)
}

```

This code will ensure that you can get log output in docker container and
syslog output in your host enviroment (non-dockerized enviroment). No additional
steps required.
