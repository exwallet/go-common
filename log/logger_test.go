package log

import (
	"fmt"
	"testing"
)

func A() {
	log := NewLogger()
	fc := &FileConfig{
		Filename:      "main.log",
		LevelFileName: nil,
		MaxSize:       0,
		MaxLine:       0,
		DateSlice:     FILE_SLICE_DATE_DAY,
		JsonFormat:    false,
		Format:        "",
	}
	e := log.Attach(LOGGER_LEVEL_INFO, fc)
	fmt.Printf("err: %+v\n", e)
	SetDefaultLogger(log)

	Debug("debug")
	Info("info")
}

func B() {
	Warn("func B ")
}

func TestNewLogger(t *testing.T) {
	A()
	B()

}

//
//func TestLogger_loggerMessageFormat(t *testing.T) {
//
//	loggerMsg := &loggerMessage{
//		Timestamp:         time.Now().Unix(),
//		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
//		Millisecond:       time.Now().UnixNano() / 1e6,
//		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
//		Level:             LOGGER_LEVEL_DEBUG,
//		LevelString:       "debug",
//		Body:              "logger console adapter test",
//		File:              "console_test.go",
//		Line:              77,
//		Function:          "TestAdapterConsole_WriteJsonFormat",
//	}
//
//	format := "%millisecond_format% [%level_string%] [%file%:%line%] %body%"
//	str := loggerMessageFormat(format, loggerMsg)
//
//	fmt.Println(str)
//}
