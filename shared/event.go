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
	// *KeyEvent
	*SystemEvent
}

type MouseEvent struct {
	TriggerButtonMask    uint8  // Only one bit will be set, and only for ButtonDown, ButtonUp, and Click events
	ButtonsMask          uint8  // Potentially multiple bits will be set, showing the state of all buttons at the time of this event
	LocationX, LocationY uint16 // Event location in the window active area. (0,0) is [system dependent? Always lower left?]
}

type SystemEvent struct {
	HandleForSurface uintptr
}
