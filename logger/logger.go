// Package logger defines the Cloudinary Logger
package logger

import (
	"log"
)

type Level int8

const NONE Level = 0
const ERROR Level = 1
const DEBUG Level = 2

type LogWriter interface {
	Debug(v ...interface{})
	Error(v ...interface{})
}

type GoLog struct{}

func (g *GoLog) Debug(v ...interface{}) {
	log.Println("cloudinary debug", v)
}

func (g *GoLog) Error(v ...interface{}) {
	log.Println("cloudinary error", v)
}

type Logger struct {
	Writer LogWriter
	level  Level
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) Error(v ...interface{}) {
	if l.level >= ERROR {
		l.Writer.Error(v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level == DEBUG {
		l.Writer.Debug(v...)
	}
}

func New() *Logger {
	return &Logger{
		Writer: &GoLog{},
		level:  ERROR,
	}
}
