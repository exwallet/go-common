package main

import (
	"github.com/exwallet/go-common/log"
)

func main() {

	_log := log.NewLogger()

	apiConfig := &log.ApiConfig{
		Url:        "http://127.0.0.1:8081/index.php",
		Method:     "GET",
		Headers:    map[string]string{},
		IsVerify:   false,
		VerifyCode: 0,
	}
	_log.Attach(log.LOGGER_LEVEL_DEBUG, apiConfig)
	_log.SetAsync()
	_log.Emergency("this is a emergency log!")
	_log.Alert("this is a alert log!")
	_log.Flush()
}
