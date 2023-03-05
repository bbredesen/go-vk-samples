package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

//go:generate glslc.exe shaders/shader.vert -o shaders/vert.spv
//go:generate glslc.exe shaders/shader.frag -o shaders/frag.spv

// This demonstration app is based on the introductory app at https://vulkan-tutorial.com and lines up with the "Drawing
// a triangle" chapter

func main() {
	fmt.Printf("Win32 Vulkan - Vertex Buffers\n")
	fmt.Printf("This program demonstrates the use of vertex buffers in Vulkan.\n\n")

	if r, ver := vk.EnumerateInstanceVersion(); r != vk.SUCCESS {
		fmt.Printf("ERROR: Could not get installed Vulkan version. Result code was %s\n", r.String())
		// ERROR_OUT_OF_HOST_MEMORY is the only error code possible, per the spec.
	} else {
		fmt.Printf("Vulkan Library API version %s\n", versionToString(ver))
	}

	app := NewApp()

	// Initialize the app and open the window
	app.Initialize("02_vertex_buffers")

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
