package shared

type EventType int

const (
	ET_None EventType = iota

	ET_Mouse_Scroll
	ET_Mouse_Move
	ET_Mouse_Drag
	ET_Mouse_ButtonDown
	ET_Mouse_ButtonUp
	ET_Mouse_Click

	ET_Key_Down
	ET_Key_Up
	ET_Key_Repeat

	// Sent when the window has been created, but not neccesarily visible. Guaranteed to be the first message sent.
	ET_Sys_Created
	// Sent when the window has closed. Guaranteed to be the last event sent.
	ET_Sys_Closed

	ET_Sys_Minimize
	ET_Sys_UnMinimize
	ET_Sys_Maximize
	ET_Sys_LostFocus
	ET_Sys_RecieveFocus
	ET_Sys_ResizeStart
	ET_Sys_ResizeProgress
	ET_Sys_ResizeComplete

	ET_Maximum
)

type EventMessage struct {
	Type EventType

	*MouseEvent
	*KeyEvent
	*SystemEvent
}

type MouseEvent struct {
	TriggerButtonMask    MouseBtnBitFlags // Only one bit will be set, and only for ButtonDown, ButtonUp, and Click events
	ButtonsMask          MouseBtnBitFlags // Potentially multiple bits will be set, showing the state of all buttons at the time of this event
	LocationX, LocationY uint16           // Event location in the window active area. (0,0) is [system dependent? Always lower left?]
	Modifiers            KeyModBitFlags
}

type KeyEvent struct {
	KeyCode uint16
	// Rune is the UTF-8 rune that the OS interprets this keypress as. It can be different from the key code
	// e.g. if shift is pressed, if the user has a key mapping configured, or if a combination of keys results
	// in a (non-ascii) unicode rune according to the operating system.
	Rune      rune
	Modifiers KeyModBitFlags
}

type KeyModBitFlags uint32

const (
	KeyModNone KeyModBitFlags = 0

	// These are intentionally not using iota for the bitshift.
	// These values are manually dupliated in go_bridge.h

	KeyModLeftShift  KeyModBitFlags = 1 << 1
	KeyModRightShift KeyModBitFlags = 1 << 2
	KeyModLeftCtrl   KeyModBitFlags = 1 << 3
	KeyModRightCtrl  KeyModBitFlags = 1 << 4
	KeyModLeftAlt    KeyModBitFlags = 1 << 5 // LeftAlt is the same as left option on Mac keyboard
	KeyModRightAlt   KeyModBitFlags = 1 << 6 // RightAlt is the same as right option on Mac keyboard
	KeyModLeftMeta   KeyModBitFlags = 1 << 7 // LeftMeta is the left command key on Mac and the left Windows key on win32
	KeyModRightMeta  KeyModBitFlags = 1 << 8 // RightMeta is the right command key on Mac and the right Windows key on win32

	//export KeyModAnyShift
	KeyModAnyShift KeyModBitFlags = KeyModLeftShift | KeyModRightShift
	KeyModAnyCtrl  KeyModBitFlags = KeyModLeftCtrl | KeyModRightCtrl
	KeyModAnyAlt   KeyModBitFlags = KeyModLeftAlt | KeyModRightAlt
	KeyModAnyMeta  KeyModBitFlags = KeyModLeftMeta | KeyModRightMeta
)

type MouseBtnBitFlags uint8

const (
	MouseBtnNone   MouseBtnBitFlags = 0
	MouseBtnLeft   MouseBtnBitFlags = 1 << 0
	MouseBtnRight  MouseBtnBitFlags = 1 << 1
	MouseBtnMiddle MouseBtnBitFlags = 1 << 2
)

// SystemEvent todo...window handle, process handle, others?
type SystemEvent struct {
	HandleForSurface          uintptr
	WindowWidth, WindowHeight uint32
}
