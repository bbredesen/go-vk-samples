package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

/*
This sample will do (nearly) the absolute minimum amount of work to validate go-vk: print the Vukan version installed on
your system, create an instance, and then destroy it.
*/

func main() {

	fmt.Printf("Minimal Vulkan - Validate API Presence\n")
	fmt.Printf("This program does the absolute minimum amount of work to validate go-vk, which is to get and print your Vulkan version, create an instance, and then destroy it.\n\n")

	if r, ver := vk.EnumerateInstanceVersion(); r != vk.SUCCESS {
		fmt.Printf("ERROR: Could not get installed Vulkan version. Result code was %d\n", r)
		// ERROR_OUT_OF_HOST_MEMORY is the only error code possible, per the spec.
	} else {
		fmt.Printf("Vulkan Library API version %s\n", versionToString(ver))
	}

	appInfo := vk.ApplicationInfo{
		PApplicationName:   "00_minimal",
		ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
		EngineVersion:      vk.MAKE_VERSION(1, 0, 0),
		ApiVersion:         vk.MAKE_VERSION(1, 2, 0),
	}

	ci := vk.InstanceCreateInfo{
		PApplicationInfo: &appInfo,
	}

	var r vk.Result
	var instance vk.Instance

	if r, instance = vk.CreateInstance(&ci, nil); r != vk.SUCCESS {
		panic(fmt.Errorf("Failed to create an instance, error code was %d", r))
	}
	fmt.Printf("Instance created, handle value is 0x%x\n", instance)

	vk.DestroyInstance(instance, nil)
	fmt.Printf("Instance destroyed, exiting...")
}

// Helper to extract parts of the Vulkan version and convert to a string
func versionToString(version uint32) string {
	return fmt.Sprintf("%d.%d.%d",
		vk.API_VERSION_MAJOR(version),
		vk.API_VERSION_MINOR(version),
		vk.API_VERSION_PATCH(version),
	)
}
