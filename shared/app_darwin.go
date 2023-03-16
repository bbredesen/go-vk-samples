//go:build darwin

package shared

/*
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore
#include "go_bridge.h"
*/
import "C"

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
)

func NewApp() (App, error) {
	app := &darwinApp{}

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
	app.caLayer = C.initCocoaWindow()

	C.runCocoaWindow()
	return nil
}

func (app *darwinApp) GetRequiredInstanceExtensions() []string {
	return []string{
		vk.KHR_SURFACE_EXTENSION_NAME,
		vk.EXT_METAL_SURFACE_EXTENSION_NAME,
	}
}

func (app *darwinApp) DelegateCreateSurface(instance vk.Instance) vk.SurfaceKHR {
	var r vk.Result
	var surface vk.SurfaceKHR

	ci := vk.MetalSurfaceCreateInfoEXT{
		PLayer: (*vk.CAMetalLayer)(app.caLayer),
	}

	if r, surface = vk.CreateMetalSurfaceEXT(instance, &ci, nil); r != vk.SUCCESS {
		panic(r)
	} else {
		return surface
	}
}
