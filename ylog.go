package ylog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

type LogFileWriter struct {
	lastHour int64
	file     *os.File
	mu       sync.Mutex
	Path     string
}

func NewLogFileWriter(path string) (fileWriter *LogFileWriter) {
	return &LogFileWriter{
		Path: path,
	}
}

// Write implement io.Writer
func (l *LogFileWriter) Write(p []byte) (n int, err error) {
	err = l.setCurrentLogFile()
	if err != nil {
		return
	}
	return l.file.Write(p)
}

func (l *LogFileWriter) setCurrentLogFile() (err error) {
	currentTime := time.Now()
	if l.file == nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.file == nil {
			l.file, err = createFile(&l.Path, &currentTime)
			l.lastHour = getTimeHour(&currentTime)
		}
		return
	}
	currentHour := getTimeHour(&currentTime)
	if l.lastHour != currentHour {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.lastHour != currentHour {
			_ = l.file.Close()
			l.file, err = createFile(&l.Path, &currentTime)
			l.lastHour = getTimeHour(&currentTime)
		}
	}
	return
}

type YLogger struct {
	level  LogLevel
	logger *log.Logger
}

func NewYLogger(level LogLevel, writer io.Writer) *YLogger {
	logger := log.New(writer, "", log.Ldate|log.Lmicroseconds)
	return &YLogger{
		level:  level,
		logger: logger,
	}
}

func (l *YLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *YLogger) canLog(level LogLevel) bool {
	return l.level <= level
}

func (l *YLogger) log(level LogLevel, levelName string, v ...interface{}) {
	if l.canLog(level) {
		funcName, fileName, line, _ := runtime.Caller(2)
		modIndex := strings.Index(runtime.FuncForPC(funcName).Name(), "/")
		codeFile := fileName[modIndex:]
		fullName := fmt.Sprintf("%s:%d", codeFile, line)
		v = append([]interface{}{levelName, fullName, "-->"}, v...)
		l.logger.Output(3, fmt.Sprintln(v...))
	}
}

// 日志包方法

var logWriter = NewLogFileWriter("logs")
var logger = NewYLogger(LevelInfo, io.MultiWriter(os.Stdout, logWriter))

func SetLevel(level LogLevel) {
	logger.SetLevel(level)
}

func GetLogWriter() io.Writer {
	return logWriter
}

func Debug(v ...interface{}) {
	logger.log(LevelDebug, "[DEBUG]", v...)
}
func DebugF(format string, v ...interface{}) {
	logger.log(LevelDebug, "[DEBUG]", fmt.Sprintf(format, v...))
}
func Info(v ...interface{}) {
	logger.log(LevelInfo, "[INFO] ", v...)
}
func InfoF(format string, v ...interface{}) {
	logger.log(LevelInfo, "[INFO] ", fmt.Sprintf(format, v...))
}
func Warn(v ...interface{}) {
	logger.log(LevelWarn, "[WARN] ", v...)
}
func WarnF(format string, v ...interface{}) {
	logger.log(LevelWarn, "[WARN] ", fmt.Sprintf(format, v...))
}
func Error(v ...interface{}) {
	logger.log(LevelError, "[ERROR]", v...)
}
func ErrorF(format string, v ...interface{}) {
	logger.log(LevelError, "[ERROR]", fmt.Sprintf(format, v...))
}

func getTimeHour(t *time.Time) int64 {
	return t.Unix() / 3600
}

func getFileName(t *time.Time) string {
	return t.Format("2006-01-02_15")
}

func createFile(path *string, t *time.Time) (file *os.File, err error) {
	dir := filepath.Join(*path, t.Format("200601"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
		if err != nil {
			return nil, err
		}
	}

	fileName := filepath.Join(dir, getFileName(t)+".txt")
	file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	return
}
