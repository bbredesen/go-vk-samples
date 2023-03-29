#ifndef __go_bridge_h__
#define __go_bridge_h__

#include <stdint.h>

#if defined(__APPLE__) && defined(__MACH__)
void* initCocoaWindow(int width, int height, int left, int top);
void runCocoaWindow();
#endif

#if defined(_WIN32)

#include <Windows.h>

HWND initWin32Window(int width, int height, int left, int top);
void runWin32Window(HWND hWnd);

#endif

void wmnotify_okToClose(uintptr_t handle);

extern void gonotify_windowCreated(uintptr_t handle);
extern void gonotify_windowWillClose(uintptr_t handle);
// extern void gonotify_windowDestroyed(uintptr_t handle);

extern void gonotify_windowResizeStart();
extern void gonotify_windowResizeProgress(uint32_t width, uint32_t height);
extern void gonotify_windowResizeComplete();

extern void gonotify_keyDown(uint16_t keyCode, uint32_t keyRune, uint32_t modifiers);
extern void gonotify_keyUp(uint16_t keyCode, uint32_t keyRune, uint32_t modifiers);

extern void gonotify_mouseDown(uint8_t button, uint16_t locationX, uint16_t locationY, uint32_t modifiers);
extern void gonotify_mouseUp(uint8_t button, uint16_t locationX, uint16_t locationY, uint32_t modifiers);

// These values are manually duplicated from event.go, as Cgo does not export any var or const values, only functions
typedef uint32_t KeyModBitFlags;

static const KeyModBitFlags KeyModLeftShift  = 1 << 1;
static const KeyModBitFlags KeyModRightShift = 1 << 2;
static const KeyModBitFlags KeyModLeftCtrl   = 1 << 3;
static const KeyModBitFlags KeyModRightCtrl  = 1 << 4;
static const KeyModBitFlags KeyModLeftAlt    = 1 << 5; // LeftAlt is the same as left option on Mac keyboard
static const KeyModBitFlags KeyModRightAlt   = 1 << 6; // RightAlt is the same as right option on Mac keyboard
static const KeyModBitFlags KeyModLeftMeta   = 1 << 7; // LeftMeta is the left command key on Mac and the left Windows key on win32
static const KeyModBitFlags KeyModRightMeta  = 1 << 8; // RightMeta is the right command key on Mac and the right Windows key on win32

static const KeyModBitFlags KeyModAnyShift = KeyModLeftShift | KeyModRightShift;
static const KeyModBitFlags KeyModAnyCtrl  = KeyModLeftCtrl | KeyModRightCtrl;
static const KeyModBitFlags KeyModAnyAlt   = KeyModLeftAlt | KeyModRightAlt;
static const KeyModBitFlags KeyModAnyMeta  = KeyModLeftMeta | KeyModRightMeta;

#endif