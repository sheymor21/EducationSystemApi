package utilities

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	sync2 "sync"
)

var Log *logrus.Logger

type CustomFormatter struct {
}

func init() {
	var sync sync2.Once
	sync.Do(func() {
		Log = logrus.New()
		Log.SetFormatter(&CustomFormatter{})
	})
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s [%s] %s: %s path=%s duration=%s\n",
		timestamp, entry.Level, entry.Message,
		entry.Data["method"], entry.Data["path"], entry.Data["duration"])
	b.WriteString(logLine)
	return b.Bytes(), nil
}
