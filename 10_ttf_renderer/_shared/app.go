package shared

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/bbredesen/go-vk"
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
		winMsgs: c,
	}
}

type Win32App struct {
	isInitialized bool

	windowTitle string

	// Windows handles and messages
	winMsgs   <-chan WindowMessage
	HInstance win32.HInstance
	HWnd      win32.HWnd

	lastFrameTime time.Time

	ClassName     string
	Width, Height uint32
}

func (app *Win32App) GetRequiredInstanceExtensions() []string {
	return []string{vk.KHR_SURFACE_EXTENSION_NAME, vk.KHR_WIN32_SURFACE_EXTENSION_NAME}
}

func (app *Win32App) createAndLoop(quitChan chan int) {
	// The OS thread that creates the window also has to run the message loop. You
	// can't createWindow and then go messageLoop, or the window simply freezes once the Go runtime attempts to run this
	// goroutine on a different OS thread.
	// This call ensures that the spawned goroutine is 1-to-1 with the current thread.
	runtime.LockOSThread()

	app.createWindow()
	quitChan <- messageLoop(app.HWnd)
}

func (app *Win32App) Initialize(windowTitle string) {
	quitChan := make(chan int)
	app.windowTitle = windowTitle

	go app.createAndLoop(quitChan)

	// Create and loop will fire the CREATE message once the window is created. Vulkan initialization will crash if HWnd
	// is not set yet, so we wait here until the window is created before moving forward.
	createMessage := <-app.winMsgs
	if createMessage.Text != "CREATE" {
		panic("Did not get CREATE as first message: " + createMessage.Text)
	}
	app.HWnd = createMessage.HWnd

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

func (app *Win32App) createWindow() {
	if app.Width == 0 {
		app.Width = 1280
	}
	if app.Height == 0 {
		app.Height = 1024
	}

	className := app.ClassName //"ttf-renderer"

	var err win32.Win32Error
	app.HInstance, err = win32.GetModuleHandleExW(0, "")
	if err != 0 {
		panic(err)
	}

	cursor, err := win32.LoadCursor(0, win32.IDC_ARROW)

	fn := wndProc

	wndClass := win32.WNDCLASSEXW{
		Style:      (win32.CS_OWNDC | win32.CS_HREDRAW | win32.CS_VREDRAW),
		WndProc:    windows.NewCallback(fn),
		Instance:   app.HInstance,
		Cursor:     cursor,
		Background: win32.HBrush(win32.COLOR_WINDOW + 1),
		ClassName:  windows.StringToUTF16Ptr(className),
	}

	wndClass.Size = uint32(unsafe.Sizeof(wndClass))

	if _, err := win32.RegisterClassExW(&wndClass); err != 0 {
		panic(err)
	}

	app.HWnd, err = win32.CreateWindowExW(
		0,
		className,
		app.windowTitle,
		(win32.WS_VISIBLE | win32.WS_OVERLAPPEDWINDOW),
		win32.SW_USE_DEFAULT,
		win32.SW_USE_DEFAULT,
		app.Width,
		app.Height,
		0,
		0,
		app.HInstance,
		nil,
	)

	if err != 0 {
		panic(err)
	}

}

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

	case win32.WM_CHAR:
		globalChannel <- WindowMessage{
			Text:      "CHAR",
			HWnd:      hwnd,
			Character: rune(wParam),
		}
	case win32.WM_KEYDOWN:
		globalChannel <- WindowMessage{
			Text:     "KEYDOWN",
			HWnd:     hwnd,
			KeyCode:  byte(wParam),
			IsRepeat: lParam&(1<<30) != 0,
		}
	case win32.WM_KEYUP:
		globalChannel <- WindowMessage{
			Text:    "KEYUP",
			HWnd:    hwnd,
			KeyCode: byte(wParam),
		}

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

type ProcessInputFunc func(keys map[byte]bool, deltaT time.Duration)
type TickFunc func(deltaT time.Duration)
type DrawFunc func()

func (app *Win32App) DefaultMainLoop(fnInput ProcessInputFunc, fnTick TickFunc, fnDraw DrawFunc) {

	// Read any system messages...input, resize, window close, etc.
	for {
	innerLoop:
		for {
			select {
			case msg := <-app.winMsgs:
				// fmt.Println(msg.Text)
				switch msg.Text {
				case "KEYDOWN":
					setAutoRepeat(msg.KeyCode)
					// Some kind of non-repeat option for a keycode would be useful...handle a single keypress, possibly
					// let the OS handle the input delay and watch msg.IsRepeat for certain keys
				case "KEYUP":
					clearAutoRepeat(msg.KeyCode)
				case "DESTROY":
					// Break out of the loop
					return

				}
			default:
				// Pull everything off the queue, then continue the outer loop
				break innerLoop // "break" will break out of the select statement, not the loop, so we have to use a break label
			}

		}

		deltaT := time.Since(app.lastFrameTime)
		zeroTime := time.Time{}
		if app.lastFrameTime == zeroTime {
			deltaT = 0
		}

		app.lastFrameTime = time.Now()

		fnInput(keyAutoRepeat(), deltaT)
		fnTick(deltaT)
		fnDraw()
	}
}

func DefaultIgnoreInput(map[byte]bool, time.Duration) {}
func DefaultIgnoreTick(time.Duration)                 {}
func DefaultIgnoreDraw()                              {}

var autoRepeater map[byte]bool

func setAutoRepeat(keyCode byte) {
	autoRepeater[keyCode] = true
}
func clearAutoRepeat(keyCode byte) {
	delete(autoRepeater, keyCode)
}

func keyAutoRepeat() map[byte]bool {
	if autoRepeater == nil {
		autoRepeater = make(map[byte]bool)
	}

	return autoRepeater
}
