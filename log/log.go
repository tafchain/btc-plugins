package log

import (
	"bufio"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var (
	gLog = logrus.New()
)

func InitLogger(path, name, level string) {
	e := os.MkdirAll(path, 0755)
	if e != nil {
		panic(e)
	}
	baseLogPath := filepath.Join(path,name)
	opts := make([]rotatelogs.Option, 0)
	err := os.Symlink(baseLogPath+"_%Y%m%d.log", baseLogPath)
	if err == nil {
		opts = append(opts, rotatelogs.WithLinkName(baseLogPath))
	}
	opts = append(opts, rotatelogs.WithMaxAge(7*24*time.Hour), rotatelogs.WithRotationTime(24*time.Hour))
	writer, err := rotatelogs.New(
		baseLogPath+"_%Y%m%d.log", opts...,
	)
	if err != nil {
		panic("config local file system logger error. ")
	}
	ew, err := rotatelogs.New(
		baseLogPath+"_err"+"_%Y%m%d.log", opts...,
	)
	if err != nil {
		panic("config local file system logger error. ")
	}

	switch level {

	case "debug":
		gLog.SetLevel(logrus.DebugLevel)
		gLog.SetOutput(os.Stderr)
	case "info":
		setNull()
		gLog.SetLevel(logrus.InfoLevel)
	case "warn":
		setNull()
		gLog.SetLevel(logrus.WarnLevel)
	case "error":
		setNull()
		gLog.SetLevel(logrus.ErrorLevel)
	default:
		setNull()
		gLog.SetLevel(logrus.InfoLevel)
	}

	g := &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	}

	lfHook := NewHook(WriterMap{
		logrus.DebugLevel: writer, 
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: ew,
		logrus.FatalLevel: ew,
		logrus.PanicLevel: ew,
	}, g)
	gLog.AddHook(lfHook)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	gLog.SetOutput(writer)
}

func Trace(args ...interface{}) {
	gLog.Trace(args...)
}

func Debug(args ...interface{}) {
	gLog.Debug(args...)
}

func Info(args ...interface{}) {
	gLog.Info(args...)
}

func Warn(args ...interface{}) {
	gLog.Warn(args...)
}

func Error(args ...interface{}) {
	gLog.Error(args...)
}

func Panic(args ...interface{}) {
	gLog.Panic(args...)
}

func Fatal(args ...interface{}) {
	gLog.Fatal(args...)
}

func Tracef(format string, args ...interface{}) {
	gLog.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	gLog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	gLog.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	gLog.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	gLog.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	gLog.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	gLog.Fatalf(format, args...)
}

func Traceln(args ...interface{}) {
	gLog.Traceln(args...)
}

func Debugln(args ...interface{}) {
	gLog.Debugln(args...)
}

func Infoln(args ...interface{}) {
	gLog.Infoln(args...)
}

func Warnln(args ...interface{}) {
	gLog.Warnln(args...)
}

func Errorln(args ...interface{}) {
	gLog.Errorln(args...)
}

func Panicln(args ...interface{}) {
	gLog.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	gLog.Fatalln(args...)
}

func IsDebug() bool {
	return gLog.GetLevel() == logrus.DebugLevel
}
