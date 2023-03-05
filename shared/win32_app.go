package shared

import (
	"runtime"
	"unsafe"

	"github.com/bbredesen/win32-toolkit"
	"golang.org/x/sys/windows"
)

// To use this application template, get a new instance from this function, call
// app.Initialize("My Window Title") to create a window and initialize Vulkan,
// and then call app.MainLoop(). The app will gracefully shut down when the
// window is closed.
func NewWin32App(c chan WindowMessage) *Win32App {
	globalChannel = c
	return &Win32App{
		winMsgs:   c,
		ClassName: "defaultClass",
	}
}

type Win32App struct {
	isInitialized bool

	windowTitle string

	ClassName string

	// Windows handles and messages
	winMsgs   <-chan WindowMessage
	HInstance win32.HInstance
	HWnd      win32.HWnd
}

func (app *Win32App) GetRequiredInstanceExtensions() []string {
	return []string{"VK_KHR_surface", "VK_KHR_win32_surface"}
}

func (app *Win32App) createAndLoop(quitChan chan int) {
	// The thread that creates the window also has to run the message loop. You
	// can't createWindow and then go messageLoop, or the window simply freezes
	// This ensures that the spawned goroutine is 1-to-1 with the current
	// thread.

	// This ensures that the spawned goroutine is 1-to-1 with the current
	// thread.
	runtime.LockOSThread()

	app.HInstance, app.HWnd = app.createWindow()
	quitChan <- messageLoop(app.HWnd)
}

func (app *Win32App) Initialize(windowTitle string) {
	quitChan := make(chan int)
	app.windowTitle = windowTitle

	go app.createAndLoop(quitChan)

	app.isInitialized = true

}

func messageLoop(hWnd win32.HWnd) int {
	var msg = win32.MSG{} //win32.HWnd(0)
	for {
		code, _ := win32.GetMessageW(&msg, win32.HWnd(0), 0, 0)
		if code == 0 {
			break
		}

		win32.TranslateMessage(&msg)
		win32.DispatchMessage(&msg)
	}
	return 0
}

// Shutdown is the reverse of Initialize
func (app *Win32App) Shutdown() {
	app.isInitialized = false
}

func (app *Win32App) IsInitialized() bool { return app.isInitialized }

func (app *Win32App) createWindow() (win32.HInstance, win32.HWnd) {

	hInstance, err := win32.GetModuleHandleExW(0, "")
	if err != 0 {
		panic(err)
	}

	cursor, err := win32.LoadCursor(0, win32.IDC_ARROW)

	fn := wndProc

	wndClass := win32.WNDCLASSEXW{
		Style:      (win32.CS_OWNDC | win32.CS_HREDRAW | win32.CS_VREDRAW),
		WndProc:    windows.NewCallback(fn),
		Instance:   hInstance,
		Cursor:     cursor,
		Background: win32.HBrush(win32.COLOR_WINDOW + 1),
		ClassName:  windows.StringToUTF16Ptr(app.ClassName),
	}

	wndClass.Size = uint32(unsafe.Sizeof(wndClass))

	if _, err := win32.RegisterClassExW(&wndClass); err != 0 {
		panic(err)
	}

	hWnd, err := win32.CreateWindowExW(
		0,
		app.ClassName,
		app.windowTitle,
		(win32.WS_VISIBLE | win32.WS_OVERLAPPEDWINDOW),
		win32.SW_USE_DEFAULT,
		win32.SW_USE_DEFAULT,
		800,
		600,
		0,
		0,
		hInstance,
		nil,
	)

	if err != 0 {
		panic(err)
	}

	return hInstance, hWnd
}

// Work remains to be done on handling messages. The original intent here was to abstract the messages so that the same
// message channel could be used on different platforms. As it stands now, this just trivially wraps the message and
// includes hwnd in the message passed to your app, which does tie everything to the Win32 API.
//
// I should look at messaging on some of the well-known input/message systems. While we don't need a full bubble-up
// mechanism like DOM Events/Javascript, nor a callback method similar to Qt's signals/slots mechanism, modeling on the
// event message formats themselves would be useful.
func wndProc(hwnd win32.HWnd, msg win32.Msg, wParam, lParam uintptr) uintptr {
	// switch msg {
	// case win32.WM_CREATE:
	// 	fmt.Printf("Message: CREATE\n")
	// 	break

	// case win32.WM_PAINT:
	// 	fmt.Println("Message: PAINT")
	// 	win32.ValidateRect(hwnd, nil)
	// 	break

	// case win32.WM_MOUSEMOVE:
	// 	fmt.Printf("Message MOUSEMOVE at %d, %d\n", lParam&0xFFFF, lParam>>16)

	// case win32.WM_CLOSE:
	// 	fmt.Println("Message: CLOSE")
	// 	win32.DestroyWindow(hwnd)
	// 	break

	// case win32.WM_DESTROY:
	// 	fmt.Println("Message: DESTROY")
	// 	win32.PostQuitMessage(0)
	// 	break

	// default:
	// 	// fmt.Printf("%s\n", msg.String())
	// 	return win32.DefWindowProcW(hwnd, msg, wParam, lParam)
	// }
	switch msg {
	case win32.WM_CREATE:
		// fmt.Printf("Message: CREATE\n")
		globalChannel <- WindowMessage{
			Text: "CREATE",
			HWnd: hwnd,
		}

	case win32.WM_PAINT:
		// fmt.Println("WM_PAINT")
		globalChannel <- WindowMessage{
			Text: "PAINT",
			HWnd: hwnd,
		}
		win32.ValidateRect(hwnd, nil)

	case win32.WM_SIZE:
		// fmt.Printf("WM_SIZE: %d x %d\n", lParam&0xFFFF, lParam>>16)
		globalChannel <- WindowMessage{
			Text: "SIZE",
			HWnd: hwnd,
		}
		// win32.ValidateRect(hwnd, nil)
	case win32.WM_ENTERSIZEMOVE:
		// fmt.Printf("WM_ENTERSIZEMOVE\n")
		globalChannel <- WindowMessage{
			Text: "ENTERSIZEMOVE",
			HWnd: hwnd,
		}

	case win32.WM_EXITSIZEMOVE:
		// fmt.Printf("WM_EXITSIZEMOVE\n")
		globalChannel <- WindowMessage{
			Text: "EXITSIZEMOVE",
			HWnd: hwnd,
		}

	// case win32.WM_SIZING:
	// 	fmt.Printf("WM_SIZING: %d x %d\n", lParam&0xFFFF, lParam>>16)
	// 	globalChannel <- WindowMessage{
	// 		Text:      "PAINT",
	// 		HWnd:      hwnd,
	// 		HInstance: hInstance,
	// 	}
	// 	win32.ValidateRect(hwnd, nil)

	case win32.WM_CLOSE:
		globalChannel <- WindowMessage{
			Text: "CLOSE",
			HWnd: hwnd,
		}
		win32.DestroyWindow(hwnd)

	case win32.WM_DESTROY:
		globalChannel <- WindowMessage{
			Text: "DESTROY",
			HWnd: hwnd,
		}
		win32.PostQuitMessage(0)

	default:
		return win32.DefWindowProcW(hwnd, msg, wParam, lParam)
	}

	return 0
}
