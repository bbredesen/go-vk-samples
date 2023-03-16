# package shared: Cross Platform Window Interface

This package implements an extremely simple interface for opening a window on
different OS platforms, and for receiving certain events from the user and
window system:

* User events: Key press, mouse movement, mouse enter/exit
* Window events: resize, minimize, lose/gain focus

To use this pacakge, do:

```go
    app := shared.NewApp()
    ch := app.GetEventChannel()
    
    go myCustomEventLoop(ch) // You have to write this, see the samples for examples

    app.Run() // Run() must execute on the main thread and it blocks until the window is closed
```

## Events

Events are fed into a channel which can be read by user code. The recommended
procedure for each frame is to read and all events (i.e. empty the channel)
before rendering the current frame. This package does not attempt to bubble or
delgate events to components in any way.

```go
type UIEvent struct {
    Type EventType

    *MouseEvent
    *KeyEvent
    *SystemEvent
}

type EventType int

const (
    ET_Mouse_Scroll EventType = iota
    ET_Mouse_Move
    ET_Mouse_Drag
    ET_Mouse_ButtonDown
    ET_Mouse_ButtonUp
    ET_Mouse_Click

    ET_Key_Down
    ET_Key_Up
    ET_Key_Repeat

    ET_Sys_Create
    ET_Sys_Minimize
    ET_Sys_UnMinimize
    ET_Sys_Maximize
    ET_Sys_Close
    ET_Sys_LostFocus
    ET_Sys_RecieveFocus
    ET_Sys_ResizeStart
    ET_Sys_ResizeProgress
    ET_Sys_ResizeComplete
)
```

### Event Detail Models

```go
type MouseEvent struct {
    TriggerButtonMask uint8 // Only one bit will be set, and only for ButtonDown, ButtonUp, and Click events
    ButtonsMask uint8 // Potentially multiple bits will be set, showing the state of all buttons at the time of this event
    LocationX, LocationY uint16 // Event location in the window active area. (0,0) is [system dependent? Always lower left?]
}

type KeyEvent struct {
    KeyCode uint8
    Rune rune // System-dependent UTF-8 rune associated with the key, if any. 
    ModifiersMask uint8 // Control, alt, shift, etc.
}
```

### System Event Model
```go
type SystemEvent struct {
    WindowIdentifier uintptr // HWnd on Windows, NSWindow* on Cocoa, etc.
    ViewSizeX, ViewSizeY uint16 // Will NOT be changed to 0, 0 on a minimize event
}
```

## Message Dispatch Model

1. Fetch a message from the OS.
2. Build the internal message based on the OS-provided detail
3. Send the message to the message channel.

The message channel will be closed by the sender once the window has been
closed. A final ET_Sys_Close message will be added to the channel before
closing. It is up to the programmer to complete any processing and clean up
after themselves before exiting the program.

### Windows

The Windows message queue lends itself well to the dispatch model above. Windows requires the developer to continuously loop and call GetMessage, then translate and dispatch that message. Windows will then call the wndProc callback with the message details for the appropriate window (assuming a multi-window application). All of this can be done in Go through Cgo and functions from the x/sys/windows package.

### Mac/Cocoa

Cocoa events are somewhat different in that the NSView directly receives events through methods on the view, like mouseMoved and keyDown. For Mac
