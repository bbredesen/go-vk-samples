package shared

import (
	"runtime"
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/sys/windows"
)

func NewGLFWApp(windowTitle string) (App, error) {
	app := &glfwApp{
		sharedApp: newSharedApp(windowTitle),
		reqWidth:  -1,
		reqHeight: -1,
		reqLeft:   -1,
		reqTop:    -1,
	}

	return app, nil
}

type glfwApp struct {
	sharedApp

	reqWidth, reqHeight, reqLeft, reqTop int

	glfwWindow *glfw.Window
	hWnd       windows.HWND
}

func (app *glfwApp) GetRequiredInstanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.KHR_WIN32_SURFACE_EXTENSION_NAME,
	}
}

func (app *glfwApp) DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error) {
	ci := vk.Win32SurfaceCreateInfoKHR{
		Hwnd: app.hWnd,
	}
	return vk.CreateWin32SurfaceKHR(instance, &ci, nil)
}

func (app *glfwApp) Run() error {
	runtime.LockOSThread()

	glfw.Init()

	if app.reqWidth < 0 || app.reqHeight < 0 {
		app.reqWidth = 640
		app.reqHeight = 480
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)

	if window, err := glfw.CreateWindow(app.reqWidth, app.reqHeight, app.title, nil, nil); err != nil {
		panic(err)
	} else {
		app.glfwWindow = window
		app.hWnd = windows.HWND(unsafe.Pointer(window.GetWin32Window()))
	}

	globalChannel <- EventMessage{
		Type: ET_Sys_Created,
		SystemEvent: &SystemEvent{
			HandleForSurface: uintptr(app.hWnd),
			WindowWidth:      uint32(app.reqWidth),
			WindowHeight:     uint32(app.reqHeight),
		},
	}

	//C.runWin32Window(C.HWND(unsafe.Pointer(app.hWnd)))
	for !app.glfwWindow.ShouldClose() {
		// No event callbacks have been set up!
		glfw.PollEvents()
	}

	return nil
}

func (app *glfwApp) OkToClose(handle uintptr) {
	// C.wmnotify_okToClose((C.uintptr_t)(handle))
}
