package main

import (
	"github.com/exwallet/go-common/gologger"
)

func main() {

	log := gologger.NewLogger()
	//default attach console, detach console
	log.Detach("console")

	consoleConfig := &gologger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}

	log.Attach(gologger.LOGGER_LEVEL_DEBUG, consoleConfig)

	log.SetAsync()

	log.Emergency("this is a emergency log!")
	log.Alert("this is a alert log!")
	log.Critical("this is a critical log!")
	log.Error("this is a error log!")
	log.Warning("this is a warning log!")
	log.Notice("this is a notice log!")
	log.Info("this is a info log!")
	log.Debug("this is a debug log!")
	log.Emergency("this is a emergency log!")
	log.Notice("this is a notice log!")
	log.Info("this is a info log!")
	log.Debug("this is a debug log!")
	log.Emergency("this is a emergency %d log!", 10)
	log.Alert("this is a alert %s log!", "format")
	log.Critical("this is a critical %s log!", "format")
	log.Error("this is a error %s log!", "format")
	log.Warning("this is a warning %s log!", "format")
	log.Notice("this is a notice %s log!", "format")
	log.Info("this is a info %s log!", "format")
	log.Debug("this is a debug %s log!", "format")
	log.Flush()
}
