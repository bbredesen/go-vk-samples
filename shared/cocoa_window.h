#ifndef __COCOA_WINDOW_H__
#define __COCOA_WINDOW_H__

#include <Cocoa/Cocoa.h>

@interface GVKWindow : NSWindow
@end

@interface GVKApplicationDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate>

// @property(readonly) NSWindowController *ctrl;
// @property(readonly) NSViewController *vctrl;
// @property(readonly) NSWindow *window;
// @property(readonly) NSView *view;

@end

@interface GVKView : NSView <NSWindowDelegate>
@end

// TODO: ViewController, WindowController needed?
#endif