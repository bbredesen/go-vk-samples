//go:build darwin

package shared

/*
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore
#include "delegate/go_bridge.h"
*/
import "C"

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
)

func NewApp() (App, error) {
	app := &darwinApp{
		newSharedApp(),
		nil,
	}

	return app, nil
}

type darwinApp struct {
	sharedApp

	caLayer unsafe.Pointer
}

// Run() must be called from the main thread and it is blocking until the window
// is closed.  You should start a goroutine to read the message channel provided
// by this App. That channel will be closed after the window is closed.
func (app *darwinApp) Run() error {
	if app.reqWidth < 0 || app.reqHeight < 0 {
		app.reqWidth = 640
		app.reqHeight = 480
	}
	app.caLayer = C.initCocoaWindow(C.int(app.reqWidth), C.int(app.reqHeight), C.int(app.reqLeft), C.int(app.reqTop))

	C.runCocoaWindow()
	return nil
}

func (app *darwinApp) GetRequiredInstanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.EXT_METAL_SURFACE_EXTENSION_NAME,
	}
}

func (app *darwinApp) DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error) {
	ci := vk.MetalSurfaceCreateInfoEXT{
		PLayer: (*vk.CAMetalLayer)(app.caLayer),
	}

	return vk.CreateMetalSurfaceEXT(instance, &ci, nil)
}
