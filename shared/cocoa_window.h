#ifndef __COCOA_WINDOW_H__
#define __COCOA_WINDOW_H__

#include <Cocoa/Cocoa.h>

@interface GVKApplicationDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate>
@end

@interface GVKView : NSView <NSWindowDelegate>
@end

#endif