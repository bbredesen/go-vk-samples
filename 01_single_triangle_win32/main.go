package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/bbredesen/go-vk"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// //go:generate glslc.exe shaders/shader.vert -o shaders/vert.spv
// //go:generate glslc.exe shaders/shader.frag -o shaders/frag.spv

// This demonstration app is based on the introductory app at https://vulkan-tutorial.com and lines up with the "Drawing
// a triangle" chapter

func main() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	log.Println("glfw initialized")

	fmt.Printf("Win32 Vulkan - Drawing a Triangle\n")
	fmt.Printf("This program demonstrates opening a window, setting up the basics for Vulkan rendering, and draws a single triangle.\n\n")

	if r, ver := vk.EnumerateInstanceVersion(); r != vk.SUCCESS {
		fmt.Println("ERROR: Could not get installed Vulkan version.")
		// ERROR_OUT_OF_HOST_MEMORY is the only error code possible, per the spec.
	} else {
		fmt.Printf("Vulkan Library API version %s\n", versionToString(ver))
	}

	app := NewApp()

	// Initialize the app and open the window
	app.Initialize("01_single_triangle_win32")

	app.InitVulkan()
	app.MainLoop()
	app.CleanupVulkan()

	fmt.Println()
	fmt.Println("Clean shutdown, exiting...")
}

// Helper to extract parts of the Vulkan version and convert to a string
func versionToString(version uint32) string {
	return fmt.Sprintf("%d.%d.%d",
		vk.API_VERSION_MAJOR(version),
		vk.API_VERSION_MINOR(version),
		vk.API_VERSION_PATCH(version),
	)
}
