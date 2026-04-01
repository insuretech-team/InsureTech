//go:build windows
// +build windows

package logger

import (
	"syscall"
	"unsafe"
)

func enableVirtualTerminalProcessing() {
	const enableVirtualTerminalProcessing uint32 = 0x0004
	const stdOutputHandle int32 = -11
	const stdErrorHandle int32 = -12

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
		mode |= enableVirtualTerminalProcessing
		_, _, _ = procSetConsoleMode.Call(h, uintptr(mode))
	}

	enable(stdOutputHandle)
	enable(stdErrorHandle)
}
