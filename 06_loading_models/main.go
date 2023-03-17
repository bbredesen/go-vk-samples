package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

//go:generate glslc shaders/shader.vert -o shaders/vert.spv
//go:generate glslc shaders/shader.frag -o shaders/frag.spv

// This demonstration app is based on the introductory app at https://vulkan-tutorial.com and lines up with the "Uniform
// Buffers" chapter

func main() {
	fmt.Printf("Win32 Vulkan - Loading Models\n")
	fmt.Printf("This program demonstrates loading and displaying an OBJ model in Vulkan.\n\n")

	if r, ver := vk.EnumerateInstanceVersion(); r != vk.SUCCESS {
		fmt.Printf("ERROR: Could not get installed Vulkan version. Result code was %s\n", r.String())
		// ERROR_OUT_OF_HOST_MEMORY is the only error code possible, per the spec.
	} else {
		fmt.Printf("Vulkan Library API version %s\n", versionToString(ver))
	}

	app := NewApp()

	// Initialize the app and open the window
	app.Run("06_loading_models")

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