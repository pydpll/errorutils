package errorutils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

func makeupL(level logrus.Level, color bool) string {
	x := strings.ToUpper(level.String())[0:4]
	n := getColorByLevel(level) //Color Code
	if !color {
		return x
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", n, x)
}
