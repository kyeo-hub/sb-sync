package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

type LogEntry struct {
	Time    string                 `json:"time"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

var (
	logger Logger
	once   sync.Once
)

type Logger struct {
	mu     sync.Mutex
	level  Level
	output io.Writer
}

func Init(level Level) {
	once.Do(func() {
		logger = Logger{
			level:  level,
			output: os.Stdout,
		}
	})
}

func SetLevel(level Level) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.level = level
}

func SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.output = output
}

func log(level Level, msg string, fields map[string]interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if level < logger.level {
		return
	}

	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Level:   levelNames[level],
		Message: msg,
		Fields:  fields,
	}

	data, _ := json.Marshal(entry)
	logger.output.Write(append(data, '\n'))
}

func Debug(msg string, fields map[string]interface{}) {
	log(DEBUG, msg, fields)
}

func Info(msg string, fields map[string]interface{}) {
	log(INFO, msg, fields)
}

func Warn(msg string, fields map[string]interface{}) {
	log(WARN, msg, fields)
}

func Error(msg string, fields map[string]interface{}) {
	log(ERROR, msg, fields)
}

func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...), nil)
}

func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...), nil)
}

func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...), nil)
}

func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...), nil)
}
