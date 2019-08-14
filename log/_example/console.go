package main

import (
	"github.com/exwallet/go-common/log"
)

func main() {

	_log := log.NewLogger()
	//default attach console, detach console
	_log.Detach("console")

	consoleConfig := &log.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}

	_log.Attach(log.LOGGER_LEVEL_DEBUG, consoleConfig)
	_log.SetAsync()
	_log.Emergency("this is a emergency log!")
	_log.Alert("this is a alert log!")
	_log.Critical("this is a critical log!")
	_log.Error("this is a error log!")
	_log.Warning("this is a warning log!")
	_log.Notice("this is a notice log!")
	_log.Info("this is a info log!")
	_log.Debug("this is a debug log!")
	_log.Emergency("this is a emergency log!")
	_log.Notice("this is a notice log!")
	_log.Info("this is a info log!")
	_log.Debug("this is a debug log!")
	_log.Emergency("this is a emergency %d log!", 10)
	_log.Alert("this is a alert %s log!", "format")
	_log.Critical("this is a critical %s log!", "format")
	_log.Error("this is a error %s log!", "format")
	_log.Warning("this is a warning %s log!", "format")
	_log.Notice("this is a notice %s log!", "format")
	_log.Info("this is a info %s log!", "format")
	_log.Debug("this is a debug %s log!", "format")
	_log.Flush()
}
