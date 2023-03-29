# package shared: Cross Platform Window Interface

This package implements an extremely simple interface for opening a window on
different OS platforms, and for receiving certain events from the user and
window system. Not all messages have been implemented yet!

* User events: Key press, mouse movement, mouse enter/exit
* Window events: resize, minimize, lose/gain focus

To use this pacakge, do:

```go
    app := shared.NewApp()
    ch := app.GetEventChannel()
    
    go myCustomEventLoop(ch) // You have to write this, see the samples for examples

    app.Run() // Run() must execute on the main thread and it blocks until the window is closed.
```

## Events

Events are fed into a channel which can be read by user code. The recommended
procedure for each frame is to read all events (i.e. empty the channel)
before rendering the current frame. This package does not attempt to bubble or
delgate events to components in any way.

## Message Dispatch Loop

1. Receive a message from the OS.
2. Build the internal message based on the OS-provided details.
3. Send the message to the message channel.

A final ET_Sys_Close message will be added to the channel when a request to close the window is received, after which
the channel will be close by the sender. It is up to the programmer to clean up any Vulkan resources and call OkToClose
before exiting the event loop. The OkToClose function signals that your application has completed cleanup and will not
attempt to draw another frame, allowing the window to be safely destroyed.
