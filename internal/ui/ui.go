package ui

import "github.com/fatih/color"

var (
	ColorErr  = color.New(color.FgRed).Add(color.Bold)
	ColorWarn = color.New(color.FgYellow).Add(color.Bold)
)

func Cli(v string) string {
	return color.YellowString(v)
}

func Old(v string) string {
	return color.RedString(v)
}

func New(v string) string {
	return color.GreenString(v)
}

func Comment(v string) string {
	return color.BlackString(v)
}
