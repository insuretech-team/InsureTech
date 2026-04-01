//go:build !windows
// +build !windows

package logger

import "os"

func enableVirtualTerminalProcessing() {
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}
	if os.Getenv("COLORTERM") == "" {
		os.Setenv("COLORTERM", "truecolor")
	}
}
