package logs

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log *logrus.Logger
)

func InitLog(ravenDSN string) {
	log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(logrus.InfoLevel)

	if ravenDSN != "" {
		log.Info("RavenDSN is not empty. Enable sentry feature.")
		hook, err := logrus_sentry.NewAsyncSentryHook(ravenDSN,
			[]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			})
		hook.Timeout = 0
		if err != nil {
			panic(err)
		}
		log.Hooks.Add(hook)
	}

}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Error(args ...interface{}) {
	log.Error(args)
}

func Warning(args ...interface{}) {
	log.Warning(args)
}

func Panic(args ...interface{}) {
	log.Panic(args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}
