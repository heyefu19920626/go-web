package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-web/cli"
	"io"
	"os"
)

func main() {
	if initLogConfig() {
		return
	}
	localCli := cli.LocalCli{}
	localCli.ExecCommand("dir")
}

func initLogConfig() bool {
	logrus.SetFormatter(&LogFormatter{})
	logPath := "log/all.log"
	_, err := os.Stat(logPath)
	if err != nil {
		logrus.Infof("logfile not exits, will create")
		_, err := os.Stat("log")
		if err != nil {
			logrus.Infof("log dir not exits, will create")
			err := os.Mkdir("log", os.ModePerm)
			if err != nil {
				logrus.Errorf("create log dir error. %+v", err)
				return true
			}
			_, err = os.Stat(logPath)
			if err != nil {
				_, err := os.Create(logPath)
				if err != nil {
					logrus.Errorf("create logfile error. %+v", err)
					return true
				}
			}
			logrus.Infof("logfile create success")
		}
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writer := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(writer)
	return false
}

// LogFormatter 自定义logrus的日志格式
type LogFormatter struct {
}

// Format 自定义logrus的日志格式，需要实现该方法
func (m *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf("%s [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}
