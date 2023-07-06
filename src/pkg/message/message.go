package message

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

type LogLevel int

const (
	WarnLevel LogLevel = iota
	InfoLevel
	DebugLevel
	TraceLevel
)

var NoProgress bool

var Ruleline = strings.Repeat("-", 72)

var LogWriter io.Writer = os.Stderr
var logLevel = InfoLevel

var logFile *os.File

type DebugWriter struct{}

func (d *DebugWriter) Write(raw []byte) (int, error) {
	Debug(string(raw))
	return len(raw), nil
}

func init() {
	pterm.ThemeDefault.SuccessMessageStyle = *pterm.NewStyle(pterm.FgLightCyan)

	pterm.Success.Prefix = pterm.Prefix{
		Text:  " ✔",
		Style: pterm.NewStyle(pterm.FgLightCyan),
	}
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "    ERROR:",
		Style: pterm.NewStyle(pterm.BgLightRed, pterm.FgBlack),
	}

	// LogFile
	ts := time.Now().Format("2006-01-02-15-04-05")

	var err error
	if logFile != nil {
		LogWriter = io.MultiWriter(os.Stderr, logFile)
		pterm.SetDefaultOutput(LogWriter)
	} else {
		if logFile, err = os.CreateTemp("", fmt.Sprintf("ecos-%s-*.log", ts)); err != nil {
			WarnErr(err, "Error saving a log file to a temporary directory")
		} else {
			LogWriter = io.MultiWriter(os.Stderr, logFile)
			pterm.SetDefaultOutput(LogWriter)
			message := fmt.Sprintf("Saving log file to %s", logFile.Name())
			Note(message)
		}
	}
}

func SetLogLevel(lvl LogLevel) {
	logLevel = lvl
	if logLevel >= DebugLevel {
		pterm.EnableDebugMessages()
	}
}

func GetLogLevel() LogLevel {
	return logLevel
}

func Debug(a ...any) {
	debugPrinter(a...)
}

func Debugf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	Debug(message)
}

func Warn(message string) {
	Warnf("%s", message)
}

func Warnf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	pterm.Println()
	pterm.Warning.Println(message)
}

func WarnErr(err error, message string) {
	Debug(err)
	Warn(message)
}

func WarnErrorf(err error, format string, a ...any) {
	Debug(err)
	Warnf(format, a...)
}

func FatalErr(err any, message string) {
	Debug(err)
	errorPrinter().Println(message)
	Debug(string(debug.Stack()))
	os.Exit(1)
}

func FatalErrorf(err any, format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	FatalErr(err, message)
}

func Info(message string) {
	Infof("%s", message)
}

func Infof(format string, a ...any) {
	if logLevel > 0 {
		message := fmt.Sprintf(format, a...)
		pterm.Info.Println(message)
	}
}

func Success(message string) {
	Successf("%s", message)
}

func Successf(format string, a ...any) {
	pterm.Success.Println(a...)
}

func Question(query string) {
	Questionf("%s", query)
}

func Questionf(format string, a ...any) {
	pterm.Println()
	query := fmt.Sprintf(format, a...)
	pterm.FgLightBlue.Println(query)
}

func Note(a ...any) {
	Notef("%s", a)
}

func Notef(format string, a ...any) {
	pterm.Println()
	note := fmt.Sprintf(format, a...)
	notePrefix := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  "NOTE",
		},
	}
	notePrefix.Println(note)
}

func Title(title string, help string) {
	titleFormatted := pterm.FgBlack.Sprint(pterm.BgWhite.Sprintf(" %s ", title))
	helpFormatted := pterm.FgGray.Sprint(help)
	pterm.Printfln("%s %s", titleFormatted, helpFormatted)
}

func HeaderInfof(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	pterm.Println()
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		WithMargin(2).
		Printfln(message)
}

func HorizontalRule() {
	pterm.Println()
	pterm.Println(Ruleline)
}

func JSONValue(value any) string {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		Debug(err, "ERROR marshalling json")
	}
	return string(bytes)
}

func debugPrinter(a ...any) {
	printer := pterm.Debug.WithShowLineNumber(logLevel > 2).WithLineNumberOffset(2)
	now := time.Now().Format(time.RFC3339)
	a = append([]any{now, " - "}, a...)

	printer.Println(a...)

	pterm.Debug.
		WithShowLineNumber(true).
		WithLineNumberOffset(2).
		WithDebugger(false).
		WithWriter(logFile).
		Println(a...)
}

func errorPrinter() *pterm.PrefixPrinter {
	return pterm.Error.WithShowLineNumber(logLevel > 2).WithLineNumberOffset(2)
}
