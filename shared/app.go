package shared

import (
	"github.com/bbredesen/go-vk"
)

var globalChannel = make(chan EventMessage, 64)

type App interface {
	// Specify the width and height of the window before opening, and the location via top and left offset.
	// For example, (800, 600, 20, 20) will position an 800 x 600 pixel window 20 pixels from the top and from the left
	// edge of the screen. Note that the renderable surface may not be the full width and height depending on
	// the underlying window system. If not called, then the default window is 640x480 at an arbitrary screen location.
	// Setting these values to a negative number will trigger the default behavior.
	SetWindowParams(width, height, left, top int)

	GetEventChannel() <-chan EventMessage
	// GetHandleForSurface returns the OS-specific handle that can be used to create a vk.SurfaceKHR.
	// The handle returned by this function will be 0 before Run() has been called and after Run() completes.
	// The handle will also be embedded in the window creation message (ET_Sys_Created).
	// GetHandleForSurface() uintptr
	Run() error

	// OkToClose notifies the window system that it is ok to close the window at this point in time. This should be
	// called by the app once a ET_Sys_Closed message is received. The window system will cause this message to be sent
	// when the user has requested to close the window. The OkToClose callback is needed to avoid destroying the window
	// in the middle of a draw call, which will cause a crash instead of a clean shutdown.
	// The app should exit the message loop after calling this function.
	OkToClose(handle uintptr)

	GetRequiredInstanceExtensions() []string

	DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error)
}

type sharedApp struct {
	title                                string
	reqWidth, reqHeight, reqLeft, reqTop int
}

func newSharedApp(windowTitle string) sharedApp {
	return sharedApp{
		title:     windowTitle,
		reqWidth:  -1,
		reqHeight: -1,
		reqLeft:   -1,
		reqTop:    -1,
	}

}

func (app *sharedApp) GetEventChannel() <-chan EventMessage {
	return globalChannel
}

func (app *sharedApp) SetWindowParams(width, height, left, top int) {
	app.reqWidth = width
	app.reqHeight = height
	app.reqLeft = left
	app.reqTop = top
}
