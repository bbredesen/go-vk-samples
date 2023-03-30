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
	fmt.Printf("go-vk - Uniform Buffers\n")
	fmt.Printf("This program demonstrates the use of uniform buffers in Vulkan.\n\n")

	if ver, err := vk.EnumerateInstanceVersion(); err != nil {
		fmt.Printf("ERROR: Could not get installed Vulkan version. Result code was %s\n", err.Error())
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

	app.Run("05_depth_buffer")

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
