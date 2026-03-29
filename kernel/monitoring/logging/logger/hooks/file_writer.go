package hooks

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

var fileLevelLabels = map[logrus.Level]string{
	logrus.TraceLevel: "TRC",
	logrus.DebugLevel: "DBG",
	logrus.InfoLevel:  "INF",
	logrus.WarnLevel:  "WRN",
	logrus.ErrorLevel: "ERR",
	logrus.PanicLevel: "PNC",
	logrus.FatalLevel: "FTL",
}

type FileWriterHook struct {
	logFolder  string
	logName    string
	fileHandle *os.File
}

func NewFileWriterHook(logFolder string, logName string) *FileWriterHook {
	logFile := path.Join(logFolder, logName+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to log to file, using default stderr %s", err)
	}

	return &FileWriterHook{
		logFolder:  logFolder,
		logName:    logName,
		fileHandle: file,
	}
}

func (hook *FileWriterHook) Fire(entry *logrus.Entry) error {
	line := hook.formatEntry(entry)
	return hook.append(line)
}

func (hook *FileWriterHook) formatEntry(entry *logrus.Entry) string {
	var b bytes.Buffer

	ts := entry.Time.Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(&b, "%s ", ts)

	label := fileLevelLabels[entry.Level]
	if label == "" {
		label = strings.ToUpper(entry.Level.String()[:3])
	}
	fmt.Fprintf(&b, "[%s] ", label)

	if logField, ok := entry.Data["log"]; ok {
		if scope, ok := logField.(string); ok {
			fmt.Fprintf(&b, "%s > ", scope)
		}
	}

	b.WriteString(entry.Message)
	b.WriteByte('\n')

	return b.String()
}

func (hook *FileWriterHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FileWriterHook) Close() error {
	if hook.fileHandle != nil {
		return hook.fileHandle.Close()
	}
	return nil
}

func (hook *FileWriterHook) append(s string) error {
	_, err := fmt.Fprint(hook.fileHandle, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write entry, %s", s)
		return err
	}
	return nil
}
