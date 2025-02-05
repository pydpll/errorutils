// lifted from rs/zerolog console.go
package console

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
)

// empty string means stderr.
func NewWriter(path string, disabled bool) zerolog.Logger {
	type lI struct {
		color int
		text  string
	} //levelInfo
	LevelColors := map[zerolog.Level]lI{
		zerolog.TraceLevel: {colorBlue, "TRC"},
		zerolog.InfoLevel:  {colorGreen, "INF"},
		zerolog.DebugLevel: {colorBlue, "DBG"},
		zerolog.WarnLevel:  {colorYellow, "WRN"},
		zerolog.ErrorLevel: {colorRed, "ERR"},
		zerolog.FatalLevel: {colorRed, "FAT"},
		zerolog.PanicLevel: {colorRed, "PAN"},
	}
	colorize := func(s interface{}, c int, disabled bool) string {
		e := os.Getenv("NO_COLOR")
		if e != "" || c == 0 {
			disabled = true
		}

		if disabled {
			return fmt.Sprintf("%s", s)
		}
		return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
	}

	var target *os.File = os.Stderr

	if path != "" {
		var err error
		target, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic("no output for logger: " + err.Error())
		}
	}
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        target,
		TimeFormat: time.RFC822,
		FormatLevel: func(i interface{}) string {
			x := LevelColors[zerolog.GlobalLevel()]
			text := strings.ToUpper(fmt.Sprintf("[%s]", x.text))
			return colorize(text, x.color, disabled)
		},
	}).Level(zerolog.DebugLevel).With().Timestamp().Logger()

}
