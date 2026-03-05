//go:build windows
// +build windows

package logger

import (
	"syscall"
	"unsafe"
)

// enableVirtualTerminalProcessing enables ANSI colors in Windows console
func enableVirtualTerminalProcessing() {
	// Attempt to enable ANSI escape sequence processing on Windows consoles
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x0004
	const STD_OUTPUT_HANDLE int32 = -11
	const STD_ERROR_HANDLE int32 = -12

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle := kernel32.NewProc("GetStdHandle")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	enable := func(stdHandle int32) {
		h, _, _ := procGetStdHandle.Call(uintptr(stdHandle))
		if h == 0 {
			return
		}
		var mode uint32
		_, _, _ = procGetConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode)))
		// Only enable virtual terminal processing for ANSI support to keep behavior identical across OSes
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
		_, _, _ = procSetConsoleMode.Call(h, uintptr(mode))
	}

	enable(STD_OUTPUT_HANDLE)
	enable(STD_ERROR_HANDLE)
}
