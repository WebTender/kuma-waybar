package main

import "errors"

type Format int

const (
	PLAIN  Format = 0
	ANSI   Format = 1
	WAYBAR Format = 2
	JSON   Format = 3
	JSONP  Format = 4
)

func parseFormat(format string) (Format, error) {
	switch format {
	case "plain":
		return PLAIN, nil
	case "ansi":
		return ANSI, nil
	case "waybar":
		return WAYBAR, nil
	case "json":
		return JSON, nil
	case "jsonp":
		return JSONP, nil
	}
	return ANSI, errors.New("invalid format: " + format)
}

type Color int

const (
	Red    Color = 1
	Green  Color = 2
	Yellow Color = 3
)

func (c Color) hex() string {
	switch c {
	case Red:
		return "#FF0000"
	case Green:
		return "#00FF00"
	case Yellow:
		return "#FFBF00"
	}
	return "#FF0000"
}

func (c Color) ansi() string {
	switch c {
	case Red:
		return "1"
	case Green:
		return "2"
	case Yellow:
		return "3"
	}
	return "1"
}

func formatCounter(str string, color Color) string {
	if format == PLAIN {
		return str
	}
	if format == ANSI {
		return "\033[38;5;" + color.ansi() + "m" + str + "\033[0m"
	}
	return "<span foreground=\"" + color.hex() + "\">" + str + "</span>"
}
