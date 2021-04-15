package logger

import (
	"fmt"
	"testing"
)

type MockLogger struct {
	ErrorMessages map[int][]string
	DebugMessages map[int][]string
}

func (m MockLogger) Debug(v ...interface{}) {
	m.DebugMessages[len(m.DebugMessages)] = toStringsSlice(v...)
}

func (m MockLogger) Error(v ...interface{}) {
	m.ErrorMessages[len(m.ErrorMessages)] = toStringsSlice(v...)
}

func TestLogger_LogLevel_Error(t *testing.T) {
	mock := MockLogger{
		ErrorMessages: map[int][]string{},
		DebugMessages: map[int][]string{},
	}
	log := Logger{Writer: mock}
	log.SetLevel(ERROR)

	log.Error("Test error")
	log.Debug("Test debug")

	if len(mock.ErrorMessages) != 1 {
		t.Error("Logger should log error messages with level ERROR")
	}

	if mock.ErrorMessages[0][0] != "Test error" {
		t.Errorf("Error log message should be \"Test error\", %v given", mock.ErrorMessages[0][0])
	}

	if len(mock.DebugMessages) > 0 {
		t.Error("Logger should not log debug messages with level ERROR")
	}
}

func TestLogger_LogLevel_Debug(t *testing.T) {
	mock := MockLogger{
		ErrorMessages: map[int][]string{},
		DebugMessages: map[int][]string{},
	}
	log := Logger{Writer: mock}
	log.SetLevel(DEBUG)

	log.Error("Test error")
	log.Debug("Test debug")

	if len(mock.ErrorMessages) != 1 {
		t.Error("Logger should log error messages with level DEBUG")
	}

	if mock.ErrorMessages[0][0] != "Test error" {
		t.Errorf("Error log message should be \"Test error\", %v given", mock.ErrorMessages[0][0])
	}

	if len(mock.DebugMessages) != 1 {
		t.Error("Logger should log debug messages with level DEBUG")
	}

	if mock.DebugMessages[0][0] != "Test debug" {
		t.Errorf("Error log message should be \"Test debug\", %v given", mock.DebugMessages[0][0])
	}
}

func TestLogger_LogLevel_None(t *testing.T) {
	mock := MockLogger{
		ErrorMessages: map[int][]string{},
		DebugMessages: map[int][]string{},
	}
	log := Logger{Writer: mock}
	log.SetLevel(NONE)

	log.Error("Test error")
	log.Debug("Test debug")

	if len(mock.ErrorMessages) > 0 {
		t.Error("Logger should not log error messages with level NONE")
	}

	if len(mock.DebugMessages) > 0 {
		t.Error("Logger should not log debug messages with level NONE")
	}
}

func toStringsSlice(v ...interface{}) []string {
	var res []string

	for _, value := range v {
		res = append(res, fmt.Sprintf("%v", value))
	}

	return res
}
