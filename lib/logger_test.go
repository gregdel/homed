package homed

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/gregdel/homed/lib/config"
)

func TestNewLoggerTimestamps(t *testing.T) {
	t.Parallel()

	enabled := true
	disabled := false

	tests := []struct {
		name     string
		config   *config.Config
		wantTime bool
	}{
		{
			name:     "production default omits timestamp",
			config:   &config.Config{},
			wantTime: false,
		},
		{
			name:     "debug default includes timestamp",
			config:   &config.Config{Debug: true},
			wantTime: true,
		},
		{
			name: "production explicit timestamp",
			config: &config.Config{
				Logging: config.Logging{Timestamps: &enabled},
			},
			wantTime: true,
		},
		{
			name: "debug explicit no timestamp",
			config: &config.Config{
				Debug:   true,
				Logging: config.Logging{Timestamps: &disabled},
			},
			wantTime: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := newLogger(tc.config, &output)

			logger.Info("test message")

			gotTime := strings.Contains(output.String(), "time=")
			if gotTime != tc.wantTime {
				t.Fatalf("time attribute presence = %v, want %v; output: %q",
					gotTime, tc.wantTime, output.String())
			}
		})
	}
}

func TestNewLoggerDebugLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		config     *config.Config
		wantOutput bool
	}{
		{
			name:       "production omits debug logs",
			config:     &config.Config{},
			wantOutput: false,
		},
		{
			name:       "debug includes debug logs",
			config:     &config.Config{Debug: true},
			wantOutput: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := newLogger(tc.config, &output)

			logger.Debug("debug message")

			gotOutput := output.Len() > 0
			if gotOutput != tc.wantOutput {
				t.Fatalf("debug output presence = %v, want %v; output: %q",
					gotOutput, tc.wantOutput, output.String())
			}
		})
	}
}

func TestColorWriterColorsKnownLevels(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	writer := colorWriter{writer: &output}
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	logger.Log(context.Background(), slog.LevelDebug, "debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	got := output.String()
	for _, want := range []string{
		"level=" + colorDebug + "DEBUG" + colorReset,
		"level=" + colorInfo + "INFO" + colorReset,
		"level=" + colorWarn + "WARN" + colorReset,
		"level=" + colorError + "ERROR" + colorReset,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q: %q", want, got)
		}
	}
}
