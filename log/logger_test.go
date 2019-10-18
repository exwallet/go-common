package log

import (
	"fmt"
	"testing"
)

func A() {
	_log := NewLogger()
	fc := &FileConfig{
		Filename:      "main.log",
		LevelFileName: nil,
		MaxSize:       0,
		MaxLine:       0,
		DateSlice:     FILE_SLICE_DATE_DAY,
		JsonFormat:    false,
		Format:        "",
	}
	e := _log.Attach(LOGGER_LEVEL_INFO, fc)
	fmt.Printf("err: %+v\n", e)
	SetDefaultLogger(_log)
	ConsoleOutOn()

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
