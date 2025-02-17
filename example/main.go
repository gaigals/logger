package main

import (
	"log"
	"log/syslog"

	"github.com/gaigals/logger"
)

// Set as false if this is docker build. Syslogs are not supported on docker.
const isLocalBuild = false

// Disable stdout and stderr output. Does not impact syslogger.
const disableStdOutput = false

func main() {
	err := logger.NewGlobalLogger(
		"test",
		"global.log",
		isLocalBuild,
		disableStdOutput,
		syslog.LOG_USER,
	)
	if err != nil {
		log.Fatalln(err)
	}

	l, err := logger.NewLogger(
		"test",
		"instance.log",
		isLocalBuild,
		disableStdOutput,
		syslog.LOG_INFO|syslog.LOG_USER,
	)
	if err != nil {
		log.Fatalln(err)
	}

	l2, err := logger.NewLogger(
		"test",
		"",
		isLocalBuild,
		disableStdOutput,
		syslog.LOG_INFO|syslog.LOG_USER,
	)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info("Global Info")
	logger.Error("Global ERORR")

	l.Alert("instance", "alert", 15)
	l.Warn("instance warn")
	l.Critical("instance crit")
	l.Emergancyf("instance emerg=%d", 1)

	l2.Info("instance2", "no file output")

	l2.Println("instance2", "println")
	l2.Printf("instance2 printlf=%d\n", 21)

	l3 := logger.NewLoggerOrFatal(
		"test",
		"test.log",
		isLocalBuild,
		disableStdOutput,
		syslog.LOG_INFO|syslog.LOG_INFO,
	)

	l3.Debug("instance3 debug test")
}
