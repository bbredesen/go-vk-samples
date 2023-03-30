package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

//go:generate glslc shaders/shader.vert -o shaders/vert.spv
//go:generate glslc shaders/shader.frag -o shaders/frag.spv

// This demonstration app is based on the introductory app at https://vulkan-tutorial.com and lines up with the "Drawing
// a triangle" chapter

func main() {
	fmt.Printf("go-vk - Drawing a Triangle\n")
	fmt.Printf("This program demonstrates opening a window, setting up the basics for Vulkan rendering, and draws a single triangle.\n\n")

	if ver, err := vk.EnumerateInstanceVersion(); err != nil {
		fmt.Println("ERROR: Could not get installed Vulkan version.")
		// ERROR_OUT_OF_HOST_MEMORY is the only error code possible, per the spec.
	} else {
		fmt.Printf("Vulkan Library API version %s\n", versionToString(ver))
	}

	app := NewApp()

	if props, err := vk.EnumerateInstanceLayerProperties(); err != nil {
		panic("Could not enumerate available layers: " + err.Error())
	} else {
		foundValidationLayers := false
		for _, p := range props {
			if p.LayerName == "VK_LAYER_KHRONOS_validation" {
				app.enableApiLayers = append(app.enableApiLayers, p.LayerName)
				foundValidationLayers = true
				break
			}
		}

		if !foundValidationLayers {
			fmt.Println("NOTE: Khronos validation layer was not found!") // This is normal if you don't have the LunarG SDK or when using MoltenVK
		}
	}

	// Initialize the app and open the window
	app.Run("01_single_triangle")

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
