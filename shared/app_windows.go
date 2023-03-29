//go:build windows

package shared

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"golang.org/x/sys/windows"
)

/*
#cgo LDFLAGS: -L delegate/ -ldelegate
#include "delegate/go_bridge.h"
*/
import "C"

func NewApp() (App, error) {
	app := &windowsApp{
		sharedApp: newSharedApp(),
		reqWidth:  -1,
		reqHeight: -1,
		reqLeft:   -1,
		reqTop:    -1,
	}

	return app, nil
}

type windowsApp struct {
	sharedApp

	reqWidth, reqHeight, reqLeft, reqTop int

	hInstance windows.Handle
	hWnd      windows.HWND
}

func (app *windowsApp) GetRequiredInstanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.KHR_WIN32_SURFACE_EXTENSION_NAME,
	}
}

func (app *windowsApp) DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error) {
	ci := vk.Win32SurfaceCreateInfoKHR{
		Hinstance: app.hInstance,
		Hwnd:      app.hWnd,
	}
	return vk.CreateWin32SurfaceKHR(instance, &ci, nil)
}

func (app *windowsApp) Run() error {

	if app.reqWidth < 0 || app.reqHeight < 0 {
		app.reqWidth = 640
		app.reqHeight = 480
	}

	tmp := C.initWin32Window(C.int(app.reqWidth), C.int(app.reqHeight), C.int(app.reqLeft), C.int(app.reqTop))
	app.hWnd = windows.HWND(unsafe.Pointer(tmp))

	C.runWin32Window(C.HWND(unsafe.Pointer(app.hWnd)))

	return nil
}

func (app *windowsApp) OkToClose(handle uintptr) {
	C.wmnotify_okToClose((C.uintptr_t)(handle))
}
