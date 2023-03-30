
#include "go_bridge.h"

#include <Windows.h>
#include <windowsx.h>
#include <wingdi.h>
#include <stdio.h>

#define UNICODE

LRESULT wndProc(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam);
KeyModBitFlags getModifierKeys();
MouseBtnBitFlags getButtons(WPARAM wParam);

HWND initWin32Window(uint16_t *title, int width, int height, int left, int top) {
    HINSTANCE hInstance = GetModuleHandle(NULL);
    HCURSOR cursor = LoadCursor(NULL, IDC_ARROW);

    LPCWSTR classname = L"go-vk-samples-shared";
    LPCWSTR wintitle = title;

    WNDCLASSW wndClass = { 0 };
    wndClass.style = CS_OWNDC|CS_HREDRAW|CS_VREDRAW;
    wndClass.lpfnWndProc = wndProc;
    wndClass.hInstance = hInstance;

    wndClass.hIcon = LoadIcon(NULL, IDI_APPLICATION); 
    wndClass.hCursor = LoadCursor(NULL, IDC_ARROW); 

    wndClass.hbrBackground = (HBRUSH)(COLOR_WINDOW+1);
    wndClass.lpszClassName = classname;

    ATOM c = RegisterClassW(&wndClass);

    HWND hWnd = CreateWindowW(
        classname, wintitle, WS_VISIBLE|WS_OVERLAPPEDWINDOW, left, top, width, height, NULL, NULL, hInstance, NULL);

    return hWnd;
}

void runWin32Window(HWND hWnd) {
    MSG msg;
    BOOL result;
    while ( (result = GetMessage(&msg, NULL, 0,0)) != 0) {
        if (result == -1) {
            printf("-1 error result in get message loop!\n");
            DWORD err = GetLastError();
            printf("Err was %d message was %d\n", err, msg.message);

            return;
        }

        TranslateMessage(&msg);
        DispatchMessageW(&msg);
    }
}

BOOL okToClose = FALSE;

// Called from a different thread than the main message loop, so we can't directly call
// DestroyWindow here.
void wmnotify_okToClose(uintptr_t hWnd) {
    okToClose = TRUE;
}

LRESULT wndProc(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam)  {

    switch (msg) {
    case WM_CLOSE:
        gonotify_windowWillClose((uintptr_t) hWnd);
        while (!okToClose) { }
		DestroyWindow(hWnd);
		break;

	case WM_DESTROY:
		PostQuitMessage(0);
		break;

    case WM_CREATE:
        gonotify_windowCreated((uintptr_t) hWnd);
        break;

    case WM_PAINT:
        {
            PAINTSTRUCT ps;
            HDC hdc = BeginPaint(hWnd, &ps);
        
            FillRect(hdc, &ps.rcPaint, (HBRUSH) (COLOR_WINDOW+1));
            EndPaint(hWnd, &ps);
            break;
        }

    case WM_LBUTTONDOWN:
        gonotify_mouseDown(MouseBtnLeft, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;
    case WM_LBUTTONUP:
        gonotify_mouseUp(MouseBtnLeft, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;

    case WM_RBUTTONDOWN:
        gonotify_mouseDown(MouseBtnRight, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;
    case WM_RBUTTONUP:
        gonotify_mouseUp(MouseBtnRight, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;

    case WM_MBUTTONDOWN:
        gonotify_mouseDown(MouseBtnMiddle, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;
    case WM_MBUTTONUP:
        gonotify_mouseUp(MouseBtnMiddle, getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;

    case WM_MOUSEMOVE:
        gonotify_mouseMove(getButtons(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam), getModifierKeys());
        break;

// TODO: Mouse clicks, movement, drag, key down, key up

    default:
        return DefWindowProcW(hWnd, msg, wParam, lParam);
    }

    return 0;
}

KeyModBitFlags getModifierKeys() {
    KeyModBitFlags modifiers = 0;
    if (GetKeyState(VK_LSHIFT)&0x8000) {
        modifiers &= KeyModLeftShift;
    }
    if (GetKeyState(VK_RSHIFT)&0x8000) {
        modifiers &= KeyModRightShift;
    }
    if (GetKeyState(VK_LCONTROL)&0x8000) {
        modifiers &= KeyModLeftCtrl;
    }
    if (GetKeyState(VK_RCONTROL)&0x8000) {
        modifiers &= KeyModRightCtrl;
    }
    if (GetKeyState(VK_LMENU)&0x8000) {
        modifiers &= KeyModLeftAlt;
    }
    if (GetKeyState(VK_RMENU)&0x8000) {
        modifiers &= KeyModRightAlt;
    }
    if (GetKeyState(VK_LWIN)&0x8000) {
        modifiers &= KeyModLeftMeta;
    }
    if (GetKeyState(VK_RWIN)&0x8000) {
        modifiers &= KeyModRightMeta;
    }

    return modifiers;
}

MouseBtnBitFlags getButtons(WPARAM wParam) {

    MouseBtnBitFlags buttons = MouseBtnNone;
    if (wParam & MK_LBUTTON) {
        buttons |= MouseBtnLeft;
    }
    if (wParam & MK_RBUTTON) {
        buttons |= MouseBtnRight;
    }
    if (wParam & MK_MBUTTON) {
        buttons |= MouseBtnMiddle;
    }

    return buttons;
}