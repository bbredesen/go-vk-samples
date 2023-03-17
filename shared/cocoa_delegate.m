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

-(void) windowWillStartLiveResize:(NSNotification *)notification {
    gonotify_windowResizeStart();
}

-(void) windowDidEndLiveResize:(NSNotification *)notification {
    gonotify_windowResizeComplete();
}

-(NSSize) windowWillResize:(NSWindow *) window toSize:(NSSize)frameSize {
    gonotify_windowResizeProgress(frameSize.width, frameSize.height);
    return frameSize;
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

uint32_t getUTF8Rune(NSEvent *event) {
    char u8[4];
    [event.characters getCString:u8 maxLength:4 encoding:NSUTF8StringEncoding];
    return *(uint32_t*)(u8);
}

uint32_t getModifiers(NSEvent *event) {
    // Modifers have to be mapped to our internal representation
    uint32_t mods = 0;

    if (event.modifierFlags&NSEventModifierFlagShift) {
        mods |= KeyModAnyShift;
    }
    if (event.modifierFlags&NSEventModifierFlagControl) {
        mods |= KeyModAnyCtrl;
    }
    if (event.modifierFlags&NSEventModifierFlagOption) {
        mods |= KeyModAnyAlt;
    }
    if (event.modifierFlags&NSEventModifierFlagCommand) {
        mods |= KeyModAnyMeta;
    }

    return mods;
}

-(void) keyDown:(NSEvent *)event {
    gonotify_keyDown(event.keyCode, getUTF8Rune(event), getModifiers(event));
}
-(void) keyUp:(NSEvent *)event {
    gonotify_keyUp(event.keyCode, getUTF8Rune(event), getModifiers(event));
}

-(void) mouseDown:(NSEvent *)event {
    gonotify_mouseDown(0, event.locationInWindow.x, event.locationInWindow.y, getModifiers(event));
}
-(void) mouseUp:(NSEvent *)event {
    gonotify_mouseUp(0, event.locationInWindow.x, event.locationInWindow.y, getModifiers(event));
}
-(void) rightMouseUp:(NSEvent *)event {
    gonotify_mouseDown(1, event.locationInWindow.x, event.locationInWindow.y, getModifiers(event));
}
-(void) rightMouseDown:(NSEvent *)event {
    gonotify_mouseUp(1, event.locationInWindow.x, event.locationInWindow.y, getModifiers(event));
}


@end
