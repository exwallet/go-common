package main

import (
	"github.com/exwallet/go-common/gologger"
)

func main() {

	log := gologger.NewLogger()

	apiConfig := &gologger.ApiConfig{
		Url:        "http://127.0.0.1:8081/index.php",
		Method:     "GET",
		Headers:    map[string]string{},
		IsVerify:   false,
		VerifyCode: 0,
	}
	log.Attach(gologger.LOGGER_LEVEL_DEBUG, apiConfig)
	log.SetAsync()

	log.Emergency("this is a emergency log!")
	log.Alert("this is a alert log!")

	log.Flush()
}
