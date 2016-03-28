package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Ldate         = log.Ldate
	Llongfile     = log.Llongfile
	Lmicroseconds = log.Lmicroseconds
	Lshortfile    = log.Lshortfile
	LstdFlags     = log.LstdFlags
	Ltime         = log.Ltime
)

type LogLevel int64

const (
	LEVEL_DEBUG = LogLevel(1 << iota)
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
	LEVEL_PANIC
)

func (l LogLevel) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "[DEBUG]"
	case LEVEL_INFO:
		return "[INFO ]"
	case LEVEL_WARNING:
		return "[WARN ]"
	case LEVEL_ERROR:
		return "[ERROR]"
	case LEVEL_PANIC:
		return "[PANIC]"
	default:
		return "[     ]"
	}
}

type Logger struct {
	log          *log.Logger
	level        LogLevel
	maxFileSize  uint64
	currFileSize uint64
	filename     string
	mu           sync.Mutex
}

var stdLog = NewLogger(os.Stdout)

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		log:         log.New(w, "", LstdFlags|Lmicroseconds),
		level:       LEVEL_DEBUG,
		maxFileSize: 100 * 1024 * 1024, // 100M
	}
}

func NewFileLog(logFile string) (*Logger, error) {
	lastSlashIndex := strings.LastIndex(logFile, "/")
	if lastSlashIndex != -1 {
		logDir := logFile[:lastSlashIndex]
		err := os.MkdirAll(logDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		filename:    logFile,
		log:         log.New(f, "", LstdFlags|Lmicroseconds),
		level:       LEVEL_DEBUG,
		maxFileSize: 100 * 1024 * 1024, // 100M
	}, nil
}

func SetLogFile(logFile string) error {
	logDir := logFile
	lastSlashIndex := strings.LastIndex(logFile, "/")
	if lastSlashIndex != -1 {
		logDir = logFile[:lastSlashIndex]
	}

	err := os.MkdirAll(logDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	stdLog = &Logger{
		filename:    logFile,
		log:         log.New(f, "", LstdFlags|Lmicroseconds),
		level:       LEVEL_DEBUG,
		maxFileSize: 100 * 1024 * 1024, // 100M
	}
	return nil
}

func SetLevel(level LogLevel) {
	atomic.StoreInt64((*int64)(&stdLog.level), int64(level))
}

func SetFlags(flags int) {
	stdLog.log.SetFlags(flags)
}

func SetMaxFileSize(maxSize uint64) {
	atomic.StoreUint64(&stdLog.maxFileSize, maxSize)
}

func Debug(format string, v ...interface{}) {
	stdLog.output(LEVEL_DEBUG, format, v...)
}
func Debugln(v ...interface{}) {
	stdLog.output(LEVEL_DEBUG, "", v...)
}

func Info(format string, v ...interface{}) {
	stdLog.output(LEVEL_INFO, format, v...)
}
func Infoln(v ...interface{}) {
	stdLog.output(LEVEL_INFO, "", v...)
}

func Warn(format string, v ...interface{}) {
	stdLog.output(LEVEL_WARNING, format, v...)
}
func Warnln(v ...interface{}) {
	stdLog.output(LEVEL_WARNING, "", v...)
}

func Error(format string, v ...interface{}) {
	stdLog.output(LEVEL_ERROR, format, v...)
}
func Errorln(v ...interface{}) {
	stdLog.output(LEVEL_ERROR, "", v...)
}

func Panic(format string, v ...interface{}) {
	stdLog.output(LEVEL_PANIC, format, v...)

	os.Exit(1)
}

func Panicln(v ...interface{}) {
	stdLog.output(LEVEL_PANIC, "", v...)

	os.Exit(1)
}

func (l *Logger) createNewLogFile(logFile string) error {
	logDir := logFile
	lastSlashIndex := strings.LastIndex(logFile, "/")
	if lastSlashIndex != -1 {
		logDir = logFile[:lastSlashIndex]
	}

	err := os.MkdirAll(logDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.log = log.New(f, "", LstdFlags|Lmicroseconds)
	return nil
}

func (l *Logger) SetLevel(level LogLevel) {
	atomic.StoreInt64((*int64)(&l.level), int64(level))
}

func (l *Logger) SetFlags(flags int) {
	l.log.SetFlags(flags)
}

func (l *Logger) SetMaxFileSize(maxSize uint64) {
	atomic.StoreUint64(&l.maxFileSize, maxSize)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.output(LEVEL_DEBUG, format, v...)
}
func (l *Logger) Debugln(v ...interface{}) {
	l.output(LEVEL_DEBUG, "", v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.output(LEVEL_INFO, format, v...)
}
func (l *Logger) Infoln(v ...interface{}) {
	l.output(LEVEL_INFO, "", v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(LEVEL_WARNING, format, v...)
}
func (l *Logger) Warnln(v ...interface{}) {
	l.output(LEVEL_WARNING, "", v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.output(LEVEL_ERROR, format, v...)
}
func (l *Logger) Errorln(v ...interface{}) {
	l.output(LEVEL_ERROR, "", v...)
}

func (l *Logger) Panic(format string, v ...interface{}) {
	l.output(LEVEL_PANIC, format, v...)

	os.Exit(1)
}
func (l *Logger) Panicln(v ...interface{}) {
	l.output(LEVEL_PANIC, "", v...)

	os.Exit(1)
}

func (l *Logger) output(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	var s string
	if format == "" {
		s = fmt.Sprint(v...)
	} else {
		s = fmt.Sprintf(format, v...)
	}

	var b bytes.Buffer
	fmt.Fprint(&b, level, " ", s)
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.filename != "" && l.currFileSize+(uint64)(b.Len()) > l.maxFileSize {
		// create new log file
		var newLogName string
		lastSlashIndex := strings.LastIndex(l.filename, "/")
		if lastSlashIndex != -1 {
			dir := l.filename[:lastSlashIndex+1]
			newLogName += dir
		}
		newLogName += os.Args[0] + "_" + time.Now().Format("2006_01_02T15_04_05") + ".log"
		//if dotIndex := strings.LastIndex(l.filename, "."); dotIndex != -1 {
		//newLogName = l.filename[:dotIndex] + "_" + time.Now().Format("2006_01_02T15_04_05") + ".log"
		//}

		err := l.createNewLogFile(newLogName)
		if err != nil {
			fmt.Println("create new log file failed, ", err)
			return
		}
		l.currFileSize = 0
	}
	l.log.Output(1, b.String())
	l.currFileSize += (uint64)(b.Len()) + 27 // add filed of time
}
