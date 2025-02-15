package errorutils

import (
	"bytes"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type logAction int

const (
	Redirect logAction = iota
	Suppress
)

func Change(l logrus.Level, a logAction) {
	customFormatter.Actions[l] = a
}

func ChangeWriter(l logrus.Level, w io.Writer) {
	customFormatter.Writers[l] = w
}

func init() {
	logrus.SetFormatter(customFormatter)
	//fix duplication issue
	// logrus.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
	// 	Writer: os.Stderr,
	// 	LogLevels: []logrus.Level{
	// 		logrus.PanicLevel,
	// 		logrus.FatalLevel,
	// 		logrus.ErrorLevel,
	// 		logrus.WarnLevel,
	// 	},
	// })
	// logrus.AddHook(&writer.Hook{ // Send info and debug logs to stdout
	// 	Writer: os.Stdout,
	// 	LogLevels: []logrus.Level{
	// 		logrus.InfoLevel,
	// 		logrus.DebugLevel,
	// 	},
	// })
}

type MyFormatter struct {
	F               logrus.Formatter
	Actions         map[logrus.Level]logAction
	TimestampFormat string
	Writers         map[logrus.Level]io.Writer
}

var customFormatter = &MyFormatter{
	F:       nil,
	Actions: make(map[logrus.Level]logAction),
	//human time
	TimestampFormat: "2006-01-02 15:04:05",
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if action, ok := f.Actions[entry.Level]; ok {
		switch action {
		case Redirect:
			//alternatively user writter hooks.
			f.Writers[entry.Level].Write([]byte(fmt.Sprintf("%s [%s] %s\n", entry.Time.Format(f.TimestampFormat), makeupL(entry.Level, false), entry.Message)))
			return nil, nil
		case Suppress:
			return nil, nil
		}
	}
	b.WriteString(fmt.Sprintf("%s [%s] %s\n", entry.Time.Format(f.TimestampFormat), makeupL(entry.Level, true), entry.Message))
	return b.Bytes(), nil
}
