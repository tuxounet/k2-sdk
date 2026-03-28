package hooks

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type FileWriterHook struct {
	logFolder  string
	logName    string
	fileHandle *os.File
}

func NewFileWriterHook(logFolder string, logName string) *FileWriterHook {
	logFile := path.Join(logFolder, logName+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
	} else {
		fmt.Fprintf(os.Stderr, "Failed to log to file, using default stderr %s", err)
	}

	return &FileWriterHook{
		logFolder:  logFolder,
		logName:    logName,
		fileHandle: file,
	}

}

func (hook *FileWriterHook) Fire(entry *logrus.Entry) error {

	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	return hook.append(line)

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
