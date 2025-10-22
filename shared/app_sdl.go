package shared

/*
#include <stdlib.h>
#include "vulkan/vulkan.h"
*/
import "C"
import (
	"github.com/bbredesen/go-vk"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
	"time"
	"unsafe"
)

func NewApp(windowTitle string) (App, error) {
	app := &sdlApp{
		newSharedApp(windowTitle),
		nil,
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	return app, nil
}

type sdlApp struct {
	sharedApp

	window *sdl.Window
}

// Run() must be called from the main thread and it is blocking until the window
// is closed.  You should start a goroutine to read the message channel provided
// by this App. That channel will be closed after the window is closed.
func (app *sdlApp) Run() error {
	runtime.LockOSThread()

	if app.reqWidth < 0 || app.reqHeight < 0 {
		app.reqWidth = 640
		app.reqHeight = 480
	}

	var err error
	app.window, err = sdl.CreateWindow(app.title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(app.reqWidth), int32(app.reqHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_VULKAN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	globalChannel <- EventMessage{Type: ET_Sys_Created}

	// this event polling loop is better be run in the main thread once per each app.drawFrame()
	for event := sdl.PollEvent(); true; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			globalChannel <- EventMessage{Type: ET_Sys_Closed}
			break
		case *sdl.WindowEvent:
			switch e.Event {
			case sdl.WINDOWEVENT_MINIMIZED:
				globalChannel <- EventMessage{Type: ET_Sys_Minimize}
			case sdl.WINDOWEVENT_RESTORED:
				globalChannel <- EventMessage{Type: ET_Sys_UnMinimize}
			case sdl.WINDOWEVENT_RESIZED:
				// resize seems atomic to the sdl lib, not sure which event type to send
				globalChannel <- EventMessage{Type: ET_Sys_ResizeStart}
				globalChannel <- EventMessage{Type: ET_Sys_ResizeProgress}
				globalChannel <- EventMessage{Type: ET_Sys_ResizeComplete}
			}
		case *sdl.KeyboardEvent:
			sendKeyEvent(e)
		case *sdl.MouseWheelEvent:
			sendMouseWheelEvent(e)
		case *sdl.MouseMotionEvent:
			sendMouseMotion(e)
		case *sdl.MouseButtonEvent:
			sendMouseButton(e)
		}
		if event == nil {
			time.Sleep(time.Millisecond * 100) // a bit less active waiting
		}
	}

	return nil
}

var eventTypeBySdlState = map[uint8]EventType{
	sdl.PRESSED:  ET_Key_Down,
	sdl.RELEASED: ET_Key_Up,
}

func sendKeyEvent(e *sdl.KeyboardEvent) {
	globalChannel <- EventMessage{
		Type: eventTypeBySdlState[e.State],
		KeyEvent: &KeyEvent{
			KeyCode: e.Keysym.Mod,
			// TODO modifiers?
		},
	}
}

var currentlyPressedMouse MouseBtnBitFlags

func sendMouseWheelEvent(e *sdl.MouseWheelEvent) {
	event := basicMouseEvent()
	event.Type = ET_Mouse_Scroll
	// the amount scrolled is e.PreciseY
	globalChannel <- event
}

func sendMouseMotion(e *sdl.MouseMotionEvent) {
	event := basicMouseEvent()
	if currentlyPressedMouse&MouseBtnLeft > 0 {
		event.Type = ET_Mouse_Drag
	} else {
		event.Type = ET_Mouse_Move
	}
	globalChannel <- event
}

var buttonFlagBySdlEnum = map[uint8]MouseBtnBitFlags{
	sdl.BUTTON_LEFT:   MouseBtnLeft,
	sdl.BUTTON_RIGHT:  MouseBtnRight,
	sdl.BUTTON_MIDDLE: MouseBtnMiddle,
}

var buttonEventBySdlEnum = map[uint8]EventType{
	sdl.PRESSED:  ET_Mouse_ButtonDown,
	sdl.RELEASED: ET_Mouse_ButtonUp,
}

func sendMouseButton(e *sdl.MouseButtonEvent) {
	buttonFlag := buttonFlagBySdlEnum[e.Button]
	if e.State == sdl.PRESSED {
		currentlyPressedMouse |= buttonFlag
	} else {
		currentlyPressedMouse &^= buttonFlag
	}

	event := basicMouseEvent()
	event.Type = buttonEventBySdlEnum[e.Button]
	event.MouseEvent.TriggerButtonMask = buttonFlag
	globalChannel <- event
}

func basicMouseEvent() EventMessage {
	x, y, _ := sdl.GetMouseState()
	return EventMessage{
		MouseEvent: &MouseEvent{
			LocationX:   uint16(x),
			LocationY:   uint16(y),
			ButtonsMask: currentlyPressedMouse,
		},
	}
}

func (app *sdlApp) GetRequiredInstanceExtensions() []string {
	return append(app.window.VulkanGetInstanceExtensions(), vk.KHR_SURFACE_EXTENSION_NAME)
}

func (app *sdlApp) DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error) {
	pointer, err := app.window.VulkanCreateSurface((*C.VkInstance)(unsafe.Pointer(instance)))
	if err != nil {
		return 0, err
	}
	s := (*vk.SurfaceKHR)(pointer)
	return *s, nil
}

func (app *sdlApp) OkToClose(handle uintptr) {
}
