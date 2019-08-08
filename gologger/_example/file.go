package main

import (
	"github.com/exwallet/go-common/gologger"
)

func main() {

	log := gologger.NewLogger()

	fileConfig := &gologger.FileConfig{
		Filename: "./test.log",
		LevelFileName: map[int]string{
			gologger.LOGGER_LEVEL_ERROR: "./error.log",
			gologger.LOGGER_LEVEL_INFO:  "./info.log",
			gologger.LOGGER_LEVEL_DEBUG: "./debug.log",
		},
		MaxSize:    1024 * 1024,
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	log.Attach(gologger.LOGGER_LEVEL_DEBUG, fileConfig)
	log.SetAsync()

	i := 0
	for {
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
		log.Emergency("this is a emergency log!")
		log.Alert("this is a alert log!")
		log.Critical("this is a critical log!")

		i += 1
		if i == 21000 {
			break
		}
	}

	log.Flush()
}
