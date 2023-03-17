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

	GetRequiredInstanceExtensions() []string

	DelegateCreateSurface(instance vk.Instance) vk.SurfaceKHR
}

type sharedApp struct{}

func (app *sharedApp) GetEventChannel() <-chan EventMessage {
	return globalChannel
}
