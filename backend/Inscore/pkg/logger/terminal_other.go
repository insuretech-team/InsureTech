//go:build !windows
// +build !windows

package logger

import "os"

// enableVirtualTerminalProcessing ensures ANSI colors behave consistently on Unix-like systems
// by setting reasonable defaults only when not already provided by the environment.
// We avoid any terminal mode syscalls since Unix terminals support ANSI by default.
func enableVirtualTerminalProcessing() {
    if os.Getenv("TERM") == "" {
        os.Setenv("TERM", "xterm-256color")
    }
    if os.Getenv("COLORTERM") == "" {
        os.Setenv("COLORTERM", "truecolor")
    }
}
