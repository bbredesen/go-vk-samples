package shared

import (
	"fmt"

	"github.com/bbredesen/go-vk"

	"github.com/bbredesen/win32-toolkit"
)

var (
	globalChannel chan<- WindowMessage
	hInstance     win32.HInstance
	hWnd          win32.HWnd
)

type WindowMessage struct {
	Text string
	HWnd win32.HWnd

	Wparam, Lparam uint // Specifically defined as 64 bits by Win32 on 64-bit systems.

	Character rune
	KeyCode   byte
	IsRepeat  bool
	// todo
}

func GetWindowExtent(hWnd win32.HWnd) vk.Extent2D {
	var rval vk.Extent2D
	var rect win32.Rect

	result := win32.GetClientRect(hWnd, &rect)

	if result != 0 {
		panic(result)
	}
	rval.Width = uint32(rect.Right - rect.Left)
	rval.Height = uint32(rect.Bottom - rect.Top)
	// fmt.Printf("Client size: %d t x %d w\n", rval.Height, rval.Width)
	return rval

}

func PrettyWin32Msg(msg win32.MSG) string {
	return fmt.Sprintf("{ message: %s, wParam: %.8x, lParam: %.8x, pt: %v }",
		win32.Msg(msg.Message).String(),
		msg.WParam, msg.LParam, msg.Pt,
	)
}
