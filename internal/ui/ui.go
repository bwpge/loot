package ui

import "github.com/fatih/color"

var (
	ColorErr  = color.New(color.FgRed).Add(color.Bold)
	ColorWarn = color.New(color.FgYellow).Add(color.Bold)
)

func ID(id string) string {
	return color.CyanString(id)
}

func Cli(v string) string {
	return color.YellowString(v)
}

func Header(h string) string {
	return color.YellowString(h)
}

func Value(v string) string {
	return color.MagentaString(v)
}

func Comment(v string) string {
	return color.BlackString(v)
}
