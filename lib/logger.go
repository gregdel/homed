package homed

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gregdel/homed/lib/config"
)

const (
	colorReset = "\x1b[0m"
	colorDebug = "\x1b[2;36m"
	colorInfo  = "\x1b[32m"
	colorWarn  = "\x1b[33m"
	colorError = "\x1b[31m"
)

func newLogger(c *config.Config, output io.Writer) *slog.Logger {
	level := slog.LevelInfo
	if c.Debug {
		level = slog.LevelDebug
	}

	if shouldColorizeLogs(c, output) {
		output = colorWriter{writer: output}
	}

	return slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceLogAttr(c),
	}))
}

func replaceLogAttr(c *config.Config) func([]string, slog.Attr) slog.Attr {
	timestamps := c.Debug
	if c.Logging.Timestamps != nil {
		timestamps = *c.Logging.Timestamps
	}

	return func(_ []string, attr slog.Attr) slog.Attr {
		if !timestamps && attr.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return attr
	}
}

func shouldColorizeLogs(c *config.Config, output io.Writer) bool {
	if !c.Debug || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	file, ok := output.(*os.File)
	if !ok {
		return false
	}

	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

type colorWriter struct {
	writer io.Writer
}

func (w colorWriter) Write(p []byte) (int, error) {
	line := string(p)
	line = replaceLevel(line, "level=DEBUG", colorDebug)
	line = replaceLevel(line, "level=INFO", colorInfo)
	line = replaceLevel(line, "level=WARN", colorWarn)
	line = replaceLevel(line, "level=ERROR", colorError)

	_, err := io.WriteString(w.writer, line)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func replaceLevel(line, level, color string) string {
	return strings.Replace(line, level, "level="+color+level[len("level="):]+colorReset, 1)
}
