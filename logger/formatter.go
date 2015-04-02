package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/colorize"
)

var (
	blueArrow    = colorize.ColorString("▶", colorize.BlueBold, colorize.None)
	whiteDiamond = colorize.ColorString("◆", colorize.WhiteBold, colorize.None)
)

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) (output []byte, err error) {
	msg := fmt.Sprintf("%s %s %s", formatLevel(entry.Level), blueArrow,
		entry.Message)

	var error string
	for key, val := range entry.Data {
		if key == "error" {
			error = fmt.Sprintf("%s", val)
			continue
		}

		msg += fmt.Sprintf(" %s %s=%s", whiteDiamond,
			colorize.ColorString(key, colorize.CyanBold, colorize.None),
			colorize.ColorString(fmt.Sprintf("%#v", val),
				colorize.GreenBold, colorize.None))
	}

	if error != "" {
		msg += "\n" + colorize.ColorString(error, colorize.Red, colorize.None)
	}

	if string(msg[len(msg)-1]) != "\n" {
		msg += "\n"
	}

	output = []byte(msg)

	return
}
