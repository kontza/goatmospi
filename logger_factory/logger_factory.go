package logger_factory

import "github.com/juju/loggo"

var (
// Main logger_factory.
logger    = loggo.GetLogger("goatmospi")
)

func GetLogger() loggo.Logger {
	return logger
}
