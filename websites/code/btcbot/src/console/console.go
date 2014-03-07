// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package console

import (
	"fmt"

	"syscall"
	"unsafe"
)

const (
	/*
	   十六进制，高四位为背景色低4位为前景色
	   0 = 黑色       8 = 灰色
	   1 = 蓝色       9 = 淡蓝色
	   2 = 绿色       A = 淡绿色
	   3 = 湖蓝色     B = 淡浅绿色
	   4 = 红色       C = 淡红色
	   5 = 紫色       D = 淡紫色
	   6 = 黄色       E = 淡黄色
	   7 = 白色       F = 亮白色
	*/
	COLOR_BLANK     = 0x00
	COLOR_BLUE      = 0x01
	COLOR_GREEN     = 0x02
	COLOR_RED       = 0x04
	COLOR_INTENSITY = 0x08

	FOREGROUND_BLUE      = 0x01
	FOREGROUND_GREEN     = 0x02
	FOREGROUND_RED       = 0x04
	FOREGROUND_INTENSITY = 0x08

	BACKGROUND_BLUE      = 0x10
	BACKGROUND_GREEN     = 0x20
	BACKGROUND_RED       = 0x40
	BACKGROUND_INTENSITY = 0x80
)

type COORD struct {
	X int16
	Y int16
}

type SMALL_RECT struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}
type CONSOLE_SCREEN_BUFFER_INFO struct {
	dwSize              COORD
	dwCursorPosition    COORD
	wAttributes         uint16
	srWindow            SMALL_RECT
	dwMaximumWindowSize COORD
}

var (
	kernel32, _ = syscall.LoadLibrary("kernel32.dll")

	_SetConsoleTitle, _            = syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
	_SetConsoleTextAttribute, _    = syscall.GetProcAddress(kernel32, "SetConsoleTextAttribute")
	_GetConsoleScreenBufferInfo, _ = syscall.GetProcAddress(kernel32, "GetConsoleScreenBufferInfo")
	_FillConsoleOutputCharacter, _ = syscall.GetProcAddress(kernel32, "FillConsoleOutputCharacterA")
	_FillConsoleOutputAttribute, _ = syscall.GetProcAddress(kernel32, "FillConsoleOutputAttribute")
	_SetConsoleCursorPosition, _   = syscall.GetProcAddress(kernel32, "SetConsoleCursorPosition")
)

func SetConsoleTitle(title string) int {

	ret, _, callErr := syscall.Syscall(
		uintptr(_SetConsoleTitle),
		1,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		0,
		0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}
	return int(ret)
}

func SetColor(color uint8) int {

	handler := syscall.Stdout
	ret, _, callErr := syscall.Syscall(
		uintptr(_SetConsoleTextAttribute),
		2,
		uintptr(handler),
		uintptr(color),
		0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}
	return int(ret)
}

func Clear() int {

	handler := syscall.Stdout
	buffer := CONSOLE_SCREEN_BUFFER_INFO{}
	/**  fill screen with blanks  Character**/
	ret, _, callErr := syscall.Syscall(
		uintptr(_GetConsoleScreenBufferInfo),
		2,
		uintptr(handler),
		uintptr(unsafe.Pointer(&buffer)),
		0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}
	dwConSize := buffer.dwSize.X * buffer.dwSize.Y

	ret, _, callErr = syscall.Syscall6(
		uintptr(_FillConsoleOutputCharacter),
		5,
		uintptr(handler),
		uintptr(0x20),
		uintptr(dwConSize),
		0,
		uintptr(unsafe.Pointer(&[...]uint32{0})),
		0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}
	/**  fill Attribut **/
	ret, _, callErr = syscall.Syscall(
		uintptr(_GetConsoleScreenBufferInfo),
		2,
		uintptr(handler),
		uintptr(unsafe.Pointer(&buffer)),
		0)
	if callErr != 0 {
		fmt.Println("zzzzzzzzzz", callErr)
	}

	ret, _, callErr = syscall.Syscall6(
		uintptr(_FillConsoleOutputAttribute),
		5,
		uintptr(handler),
		uintptr(buffer.wAttributes),
		uintptr(dwConSize),
		0,
		uintptr(unsafe.Pointer(&[...]uint32{0})),
		0)
	if callErr != 0 {
		fmt.Println("zzzzzzzzzz", callErr)
	}
	/** set cursor position **/
	ret, _, callErr = syscall.Syscall(
		uintptr(_SetConsoleCursorPosition),
		2,
		uintptr(handler),
		uintptr(0),
		0)
	if callErr != 0 {
		fmt.Println("callErr", callErr)
	}

	return int(ret)
}
func FreeLib() {
	syscall.FreeLibrary(kernel32)
}
func init() {
	defer FreeLib()

	SetConsoleTitle("★★Robot★★")
}
