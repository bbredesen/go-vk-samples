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
	// todo
}

// func windowProc(hwnd win32.HWnd, uMsg win32.Msg, wParam uintptr, lParam uintptr) uintptr {
// 	switch uMsg {
// 	case win32.WM_CREATE:
// 		// fmt.Printf("Message: CREATE\n")
// 		globalChannel <- WindowMessage{
// 			Text:      "CREATE",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}

// 	case win32.WM_PAINT:
// 		// fmt.Println("WM_PAINT")
// 		globalChannel <- WindowMessage{
// 			Text:      "PAINT",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}
// 		win32.ValidateRect(hwnd, nil)

// 	case win32.WM_SIZE:
// 		// fmt.Printf("WM_SIZE: %d x %d\n", lParam&0xFFFF, lParam>>16)
// 		globalChannel <- WindowMessage{
// 			Text:      "SIZE",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}
// 		// win32.ValidateRect(hwnd, nil)
// 	case win32.WM_ENTERSIZEMOVE:
// 		// fmt.Printf("WM_ENTERSIZEMOVE\n")
// 		globalChannel <- WindowMessage{
// 			Text:      "ENTERSIZEMOVE",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}

// 	case win32.WM_EXITSIZEMOVE:
// 		// fmt.Printf("WM_EXITSIZEMOVE\n")
// 		globalChannel <- WindowMessage{
// 			Text:      "EXITSIZEMOVE",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}

// 	// case win32.WM_SIZING:
// 	// 	fmt.Printf("WM_SIZING: %d x %d\n", lParam&0xFFFF, lParam>>16)
// 	// 	globalChannel <- WindowMessage{
// 	// 		Text:      "PAINT",
// 	// 		HWnd:      hwnd,
// 	// 		HInstance: hInstance,
// 	// 	}
// 	// 	win32.ValidateRect(hwnd, nil)

// 	case win32.WM_CLOSE:
// 		globalChannel <- WindowMessage{
// 			Text:      "CLOSE",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}
// 		win32.DestroyWindow(hwnd)

// 	case win32.WM_DESTROY:
// 		globalChannel <- WindowMessage{
// 			Text:      "DESTROY",
// 			HWnd:      hwnd,
// 			HInstance: hInstance,
// 		}
// 		win32.PostQuitMessage(0)

// 	default:
// 		return win32.DefWindowProcW(hwnd, uMsg, wParam, lParam)
// 	}

// 	return 0
// }

// func CreateWin32Window(messageChannel chan<- WindowMessage, windowTitle string) {
// 	runtime.LockOSThread()

// 	globalChannel = messageChannel
// 	// go doCreateWin32()
// 	doCreateWin32(windowTitle)

// 	MessageLoop(hWnd)
// }

// func doCreateWin32(windowTitle string) {
// 	className := "testClass"

// 	lHInstance, werr := win32.GetModuleHandleExW(0, "")
// 	if werr != 0 {
// 		panic(werr)
// 	}

// 	cursor, werr := win32.LoadCursor(win32.HInstance(0), win32.IDC_ARROW)
// 	if werr != 0 {
// 		panic(werr)
// 	}

// 	fn := windowProc

// 	cname, err := windows.UTF16PtrFromString(className)
// 	if err != nil {
// 		panic(err)
// 	}

// 	wcx := win32.WNDCLASSEXW{
// 		Style:      win32.CS_OWNDC | win32.CS_HREDRAW | win32.CS_VREDRAW,
// 		WndProc:    windows.NewCallback(fn),
// 		Instance:   hInstance,
// 		Cursor:     cursor,
// 		Background: win32.HBrush(win32.COLOR_WINDOW + 1),
// 		ClassName:  cname,
// 	}
// 	wcx.Size = uint32(unsafe.Sizeof(wcx))

// 	if _, werr = win32.RegisterClassExW(&wcx); werr != 0 {
// 		panic(werr)
// 	}

// 	lHWnd, werr := win32.CreateWindowExW(
// 		0,
// 		className,
// 		windowTitle,
// 		uint32(win32.WS_VISIBLE|win32.WS_OVERLAPPEDWINDOW),
// 		win32.SW_USE_DEFAULT,
// 		win32.SW_USE_DEFAULT,
// 		800,
// 		800,
// 		0,
// 		0,
// 		lHInstance,
// 		nil,
// 	)

// 	if werr != 0 {
// 		panic(werr)
// 	}

// 	hInstance = lHInstance
// 	hWnd = lHWnd

// 	// fmt.Printf("Goroutine initiated, window handle is 0x%.8x (inst: 0x%.8x)\n", hWnd, hInstance)

// }

// func MessageLoop(hwnd win32.HWnd) {
// 	// fmt.Println("messageLoop starting...")
// 	for DoDispatch(hwnd) {
// 	}
// 	// fmt.Println("MessageLoop exiting")
// }

// func DoDispatch(hwnd win32.HWnd) bool {
// 	msg := win32.MSG{}
// 	// gotMessage, err := peekMessage(&msg, hwnd, 0, 0)
// 	gotMessage, err := win32.GetMessageW(&msg, 0, 0, 0)

// 	// fmt.Printf("gotMesasge / err? %v / %v\n", gotMessage, err != 0)

// 	if err != 0 {
// 		// There is a legitimate (?) error that occurs when the window is closed, which
// 		// invalidates hwnd.
// 		fmt.Printf("non-nil error, msg: %s\n", PrettyWin32Msg(msg))
// 		return false
// 		// panic(err)
// 	}

// 	// if msg.message != lastMsg {
// 	// if msg.message != 160 {
// 	// fmt.Printf("Message: (%d) %s\n", msg.Message, PrettyWin32Msg(msg))
// 	// }
// 	// lastMsg = msg.Message
// 	// }
// 	// fmt.Print("m")

// 	win32.TranslateMessage(&msg)
// 	win32.DispatchMessage(&msg)

// 	return gotMessage != 0
// }

// func GetRequiredInstanceExtensions() []string {
// 	return []string{"VK_KHR_surface", "VK_KHR_win32_surface"}
// }

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
