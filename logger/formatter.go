package logger

import (
	"github.com/blckur/blckur/colorize"
	"github.com/Sirupsen/logrus"
	"fmt"
)

var (
	blueArrow = colorize.ColorString("▶", colorize.BlueBold, colorize.None)
	whiteDiamond = colorize.ColorString("◆", colorize.WhiteBold, colorize.None)
)

type Formatter struct {}

func (f *Formatter) Format(entry *logrus.Entry) (output []byte, err error) {
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
			colorize.PurpleBold, colorize.None))
	}

	if error != "" {
		msg += "\n" + colorize.ColorString(error, colorize.Red, colorize.None)
	}

	output = []byte(msg)

	return
}
