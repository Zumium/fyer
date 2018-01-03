package log

import (
	"github.com/op/go-logging"
	"testing"
	"github.com/spf13/viper"
)

func TestLogging(t *testing.T) {
	viper.Set("log_level", "INFO")
	Init()

	testLogger := logging.MustGetLogger("test")
	testLogger.Debug("this is a debug")
	testLogger.Info("this is an info")
	testLogger.Notice("this is a notice")
	testLogger.Warning("this is a warning")
	testLogger.Error("this is an error")
	testLogger.Critical("this is a critical")
}
