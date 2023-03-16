package shared

import (
	"github.com/bbredesen/go-vk"
)

var globalChannel = make(chan EventMessage, 64)

type App interface {
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
