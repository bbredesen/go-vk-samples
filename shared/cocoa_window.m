#import "cocoa_window.h"
#import "go_bridge.h"

#include <stdio.h>

void* initCocoaWindow() {
    [NSAutoreleasePool new];
    [NSApplication sharedApplication];

    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    
    id appName = [[NSProcessInfo processInfo] processName];

    id menubar = [[NSMenu new] autorelease];
    [menubar initWithTitle:appName];

    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];

    id appMenu = [[NSMenu new] autorelease];
    id quitTitle = [@"Quit " stringByAppendingString:appName];
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"] autorelease];
    [appMenu addItem:quitMenuItem];
    
    [appMenuItem setSubmenu:appMenu];

    NSRect contentRect = NSMakeRect(0, 0, 640, 480);

    id window = [[[NSWindow alloc] 
            initWithContentRect:contentRect
            styleMask: NSWindowStyleMaskTitled|NSWindowStyleMaskClosable|NSWindowStyleMaskMiniaturizable
            backing:NSBackingStoreBuffered 
            defer:NO]
        autorelease];
    
    id del = [[GVKApplicationDelegate alloc] init];
    [NSApp setDelegate:del];
    [window setDelegate:del];

    GVKView *view = [[GVKView alloc] initWithFrame:contentRect];
    [window setContentView:view];

    [window cascadeTopLeftFromPoint:NSMakePoint(20,20)];
    [window setTitle:appName];
    [window makeKeyAndOrderFront:nil];
    [NSApp activateIgnoringOtherApps:YES];

    gonotify_windowCreated((uintptr_t)view.layer);

    return view.layer;
}

void runCocoaWindow() {
    [NSApp run];
}