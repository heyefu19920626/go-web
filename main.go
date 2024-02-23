package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-web/cli"
)

func main() {
	logrus.SetFormatter(&LogFormatter{})
	logrus.Info("Hello word!")
	localCli := cli.LocalCli{}
	localCli.ExecCommand("dir")
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
