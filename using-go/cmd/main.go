package main

import (
	pkgsrvr "using-kit/using-go/server"
	"using-kit/using-go/thingsvc"

	log "github.com/go-kit/kit/log/logrus"
	"github.com/sirupsen/logrus"
)

func main() {
	logrusLogger := logrus.New()
	logrusLogger.Formatter = &logrus.TextFormatter{TimestampFormat: "02-01-2006 15:04:05", FullTimestamp: true}
	logger := log.NewLogger(logrusLogger)
	ts := thingsvc.NewThingSvc()

	err := pkgsrvr.Run(ts)
	if err != nil {
		_ = logger.Log("error:", err)
	}
}
