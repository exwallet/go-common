package main

import (
	"github.com/exwallet/go-common/log"
)

func main() {

	_log := log.NewLogger()

	fileConfig := &log.FileConfig{
		Filename: "./test.log",
		LevelFileName: map[int]string{
			log.LOGGER_LEVEL_ERROR: "./error.log",
			log.LOGGER_LEVEL_INFO:  "./info.log",
			log.LOGGER_LEVEL_DEBUG: "./debug.log",
		},
		MaxSize:    1024 * 1024,
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	_log.Attach(log.LOGGER_LEVEL_DEBUG, fileConfig)
	_log.SetAsync()

	i := 0
	for {
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
		_log.Emergency("this is a emergency log!")
		_log.Alert("this is a alert log!")
		_log.Critical("this is a critical log!")

		i += 1
		if i == 21000 {
			break
		}
	}

	_log.Flush()
}
