# 01_single_triangle_win32

![](screenshot.png?raw=true)

This basic Vulkan application follows the first chapter (Drawing a Triangle) of the excellent [Vulkan
Tutorial](https://vulkan-tutorial.com). It demonstrates:

* device selection and setup
* setup of the swap chain and surface presentation
* creation of a graphics pipeline and renderpass
* creation and recording of command buffers
* creation and use of synchronization constructs
* clean shutdown and destruction of allocated Vulkan resources

This program was built both as a guide to using Vulkan with Go, and as a basic validation and test of the [go-vk
binding](https://github.com/bbredesen/go-vk). Barring any design updates in go-vk as we move towards a 1.0 release
(e.g. struct member naming conventions), this project should not have any material changes.

## Building

Before running, be sure to `go generate` in the project root. This will compile the shaders with glslc.exe. That
compiler is bundled with the Vulkan SDK. The go:generate annotations can be found in main.go.
