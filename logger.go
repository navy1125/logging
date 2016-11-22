package logging

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"

	"git.code4.in/mobilegameserver/unibase/unitime"
)

type Record struct {
	Time    time.Time
	Level   logLevel
	Message string
}

type Emitter interface {
	Emit(string, *Record)
}

func init() {
}

type LogServer func(id uint64, name, classname, servername string, level, timestamp uint32, log string)
type Logger struct {
	Name           string
	Handlers       map[string]Emitter
	LogServer      LogServer
	logServerLevel logLevel
	log2server     bool
	ChanLogRecord  chan *Record
}

func NewLogger() *Logger {
	ret := &Logger{
		Handlers:       make(map[string]Emitter),
		logServerLevel: ERROR,
		ChanLogRecord:  make(chan *Record, 1024), //日志多线程
	}
	go func() {
		for true {
			if rd, ok := <-ret.ChanLogRecord; ok == true {
				ret.LogWrite(rd)
			} else {
				break
			}
		}
	}()
	return ret
}

var DefaultLogger = NewLogger()

func (l *Logger) AddLoggerServerFunc(fun LogServer) {
	l.LogServer = fun
}
func (l *Logger) SetLogServerLevel(level logLevel) {
	l.logServerLevel = level
}
func (l *Logger) AddHandler(name string, h Emitter) {
	oldHandler, ok := l.Handlers[name]
	if ok {
		closer, ok := oldHandler.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}
	l.Handlers[name] = h
	l.Name = name
	if DefaultLogger.Name == "" {
		DefaultLogger.Name = name
	}
}

var ()

func (l *Logger) LogWrite(rd *Record) {
	for _, h := range l.Handlers {
		h.Emit(l.Name, rd)
	}
}
func (l *Logger) Log(level logLevel, format string, values ...interface{}) {
	rd := &Record{
		Time:    unitime.Time.Now(),
		Level:   level,
		Message: fmt.Sprintf(format, values...),
	}
	//l.ChanLogRecord <- rd
	l.LogWrite(rd)
	if l.LogServer != nil && l.log2server == false && l.logServerLevel >= level {
		l.log2server = true
		l.LogServer(0, "", "", "", uint32(level), uint32(unitime.Time.Sec()), rd.Message)
		l.log2server = false
	}
}

func (l *Logger) Debug(format string, values ...interface{}) {
	l.Log(DEBUG, format, values...)
}

func (l *Logger) Info(format string, values ...interface{}) {
	l.Log(INFO, format, values...)
}

func (l *Logger) Warning(format string, values ...interface{}) {
	l.Log(WARNING, format, values...)
}

func (l *Logger) Error(format string, values ...interface{}) {
	l.Log(ERROR, format, values...)
}
func (l *Logger) Output(calldepth int, s string) error {

	l.Log(logLevel(calldepth), s)
	return nil
}

func (l *Logger) ResetLogLevel(level string) {
	for _, e := range l.Handlers {
		if h, ok := e.(*Handler); ok {
			h.SetLevelString(level)
		}
	}
}

//打印日志用，根据回退堆栈层级获取文件名和行号信息
//参数：需要回退的堆栈层数
func GetLogBtInfo(level int) string {
	if level < 0 { //参数错误
		return ""
	}
	format := ""
	level += 1 //函数自身占一层
	_, file, line, ok := runtime.Caller(level)
	if ok == true {
		file = filepath.Base(file)
		prefix := fmt.Sprintf("[%s:%d] ", file, line)
		format = prefix + format
	}
	return format
}

func AddHandler(name string, h Emitter) {
	DefaultLogger.AddHandler(name, h)
}
func AddLoggerServerFunc(fun LogServer) {
	DefaultLogger.LogServer = fun
}

func Log(level logLevel, format string, values ...interface{}) {
	DefaultLogger.Log(level, format, values...)
}

func Debug(format string, values ...interface{}) {
	//format = GetLogBtInfo(1) + format //回退一层到原始栈
	DefaultLogger.Log(DEBUG, format, values...)
}

func Info(format string, values ...interface{}) {
	//format = GetLogBtInfo(1) + format //回退一层到原始栈
	DefaultLogger.Log(INFO, format, values...)
}

func Warning(format string, values ...interface{}) {
	//format = GetLogBtInfo(1) + format //回退一层到原始栈
	DefaultLogger.Log(WARNING, format, values...)
}

func Error(format string, values ...interface{}) {
	//format = GetLogBtInfo(1) + format //回退一层到原始栈
	DefaultLogger.Log(ERROR, format, values...)
}

func ResetLogLevel(level string) {
	DefaultLogger.ResetLogLevel(level)
}
