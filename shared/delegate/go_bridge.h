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

extern void gonotify_windowCreated(uintptr_t handle);
extern void gonotify_windowWillClose(uintptr_t handle);
extern void gonotify_keyDown(uint16_t keyCode);
extern void gonotify_mouseDown(uint8_t button, uint16_t locationX, uint16_t locationY);
extern void gonotify_mouseUp(uint8_t button, uint16_t locationX, uint16_t locationY);

#endif