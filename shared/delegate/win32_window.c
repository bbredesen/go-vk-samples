
#include "go_bridge.h"

#include <Windows.h>
#include <stdio.h>

#define UNICODE

LRESULT wndProc(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam);

HWND initWin32Window(int width, int height, int left, int top) {
    HINSTANCE hInstance = GetModuleHandle(NULL);
    HCURSOR cursor = LoadCursor(NULL, IDC_ARROW);

    LPCWSTR classname = L"go-vk-samples-shared";
    LPCWSTR wintitle = L"TODO SHARED WINDOW";

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

LRESULT wndProc(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam)  {

    switch (msg) {
    case WM_CLOSE:
		DestroyWindow(hWnd);
        gonotify_windowWillClose((uintptr_t) hWnd);
		return 0;

	case WM_DESTROY:
		PostQuitMessage(0);
		return 0;

// TODO: Mouse clicks, movement, drag, key down, key up

    default:
        return DefWindowProcW(hWnd, msg, wParam, lParam);
    }

    return 0;
}