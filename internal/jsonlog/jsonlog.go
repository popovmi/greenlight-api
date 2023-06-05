package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (self *Logger) PrintDebug(message string, properties map[string]string) {
	self.print(LevelDebug, message, properties)
}

func (self *Logger) PrintInfo(message string, properties map[string]string) {
	self.print(LevelInfo, message, properties)
}

func (self *Logger) PrintError(err error, properties map[string]string) {
	self.print(LevelError, err.Error(), properties)
}

func (self *Logger) PrintFatal(err error, properties map[string]string) {
	self.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (self *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < self.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	self.mu.Lock()
	defer self.mu.Unlock()

	return self.out.Write(append(line, '\n'))
}

func (self *Logger) Write(message []byte) (n int, err error) {
	return self.print(LevelError, string(message), nil)
}
