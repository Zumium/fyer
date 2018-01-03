package log

import (
	"github.com/Zumium/fyer/cfg"
	"github.com/op/go-logging"
	"os"
)

var loggingFormat = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

func Init() error {
	stdoutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutFormatter := logging.NewBackendFormatter(stdoutBackend, loggingFormat)
	stdoutLeveled := logging.AddModuleLevel(stdoutFormatter)
	stdoutLeveled.SetLevel(cfg.LogLevel(), "")
	logging.SetBackend(stdoutLeveled)
	return nil
}
