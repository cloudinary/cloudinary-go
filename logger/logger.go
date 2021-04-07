// Package logger defines the Cloudinary Logger
package logger

import (
	"log"
)

type Level int8

const NONE Level = 0
const ERROR Level = 1
const DEBUG Level = 2

type LoggerFunction func(...interface{})
type Logger struct {
	ErrorLogger LoggerFunction
	DebugLogger LoggerFunction
	level       Level
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) Error(v ...interface{}) {
	if l.level >= ERROR {
		l.ErrorLogger(v)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level == DEBUG {
		l.DebugLogger(v)
	}
}

func New() *Logger {
	return &Logger{
		ErrorLogger: func(v ...interface{}) {
			log.Println("cloudinary error", v)
		},
		DebugLogger: func(v ...interface{}) {
			log.Println("cloudinary debug", v)
		},
		level: ERROR,
	}
}
