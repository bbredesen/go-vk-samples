package shared

/*
#include "delegate/go_bridge.h"
*/
import "C"

//export gonotify_windowCreated
func gonotify_windowCreated(handle uintptr) {
	globalChannel <- EventMessage{
		Type: ET_Sys_Created,
		SystemEvent: &SystemEvent{
			HandleForSurface: handle,
		},
	}
}

//export gonotify_windowWillClose
func gonotify_windowWillClose(handle uintptr) {
	// fmt.Println("window will close notification received")
	globalChannel <- EventMessage{
		Type: ET_Sys_Closed,
	}
	close(globalChannel)
}

//export gonotify_windowResizeStart
func gonotify_windowResizeStart() {
	globalChannel <- EventMessage{
		Type: ET_Sys_ResizeStart,
	}
}

//export gonotify_windowResizeProgress
func gonotify_windowResizeProgress(width, height uint32) {
	globalChannel <- EventMessage{
		Type: ET_Sys_ResizeProgress,
		SystemEvent: &SystemEvent{
			WindowWidth:  width,
			WindowHeight: height,
		},
	}
}

//export gonotify_windowResizeComplete
func gonotify_windowResizeComplete() {
	globalChannel <- EventMessage{
		Type: ET_Sys_ResizeComplete,
	}
}

//export gonotify_keyDown
func gonotify_keyDown(keyCode uint16, keyRune uint32, modifiers uint32) {
	msg := EventMessage{
		Type: ET_Key_Down,
		KeyEvent: &KeyEvent{
			KeyCode:   keyCode,
			Rune:      rune(keyRune),
			Modifiers: KeyModBitFlags(modifiers),
		},
	}
	globalChannel <- msg
}

//export gonotify_keyUp
func gonotify_keyUp(keyCode uint16, keyRune uint32, modifiers uint32) {
	msg := EventMessage{
		Type: ET_Key_Up,
		KeyEvent: &KeyEvent{
			KeyCode:   keyCode,
			Rune:      rune(keyRune),
			Modifiers: KeyModBitFlags(modifiers),
		},
	}
	globalChannel <- msg
}

//export gonotify_mouseDown
func gonotify_mouseDown(button uint8, locationX, locationY uint16, modifiers uint32) {
	// fmt.Println("Button", button, "at", locationX, locationY)
	msg := EventMessage{
		Type: ET_Mouse_ButtonDown,
		MouseEvent: &MouseEvent{
			TriggerButtonMask: button,
			LocationX:         uint16(locationX),
			LocationY:         uint16(locationY),
			Modifiers:         KeyModBitFlags(modifiers),
		},
	}
	globalChannel <- msg
}

//export gonotify_mouseUp
func gonotify_mouseUp(button uint8, locationX, locationY uint16, modifiers uint32) {
	// fmt.Println("Button", button, "at", locationX, locationY)
	msg := EventMessage{
		Type: ET_Mouse_ButtonUp,
		MouseEvent: &MouseEvent{
			TriggerButtonMask: button,
			LocationX:         uint16(locationX),
			LocationY:         uint16(locationY),
			Modifiers:         KeyModBitFlags(modifiers),
		},
	}
	globalChannel <- msg
}
