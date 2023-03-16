//go:build darwin

package shared

/*
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore
#include "go_bridge.h"
*/
import "C"

import (
	"unsafe"
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

func (app *darwinApp) GetHandleForSurface() uintptr {
	return uintptr(app.caLayer)
}
