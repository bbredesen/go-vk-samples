#import "cocoa_window.h"
#include <objc/objc.h>
#import "go_bridge.h"

#import <MetalKit/MetalKit.h>

@implementation GVKApplicationDelegate:NSObject 
// - (void)applicationWillFinishLaunching:(NSNotification *)notification {
//     printf("*** delegate will finish launching\n");
// }

-(void) windowWillClose:(NSNotification *) notification
{
    gonotify_windowWillClose((uintptr_t) notification.object);
    [NSApp stop:(nil)];
}

@end


@implementation GVKView

- (GVKView*) initWithFrame:(NSRect) frame
{
    [super initWithFrame:frame];

    [[NSColor blackColor] set];
    NSRectFill([self bounds]);
    
    self.wantsLayer = YES;
    self.layer = [CAMetalLayer layer];

    return self;
}

-(BOOL) acceptsFirstResponder {
    return YES;
}

-(void) keyDown:(NSEvent *)event {
    gonotify_keyDown(event.keyCode); // more needed here, like modifer keys
}

-(void) mouseDown:(NSEvent *)event {
    gonotify_mouseDown(event.buttonNumber, event.locationInWindow.x, event.locationInWindow.y);
}
-(void) mouseUp:(NSEvent *)event {
    gonotify_mouseUp(event.buttonNumber, event.locationInWindow.x, event.locationInWindow.y);
}
-(void) rightMouseUp:(NSEvent *)event {
    gonotify_mouseDown(event.buttonNumber, event.locationInWindow.x, event.locationInWindow.y);
}
-(void) rightMouseDown:(NSEvent *)event {
    gonotify_mouseUp(event.buttonNumber, event.locationInWindow.x, event.locationInWindow.y);
}


@end
