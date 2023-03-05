# 04_texture_mapping_win32

This basic Vulkan application follows the fifth chapter (Depth Buffers) of the excellent [Vulkan
Tutorial](https://vulkan-tutorial.com), building on the application at `04_texture_mapping`. It demonstrates creation of
the depth resources and testing against the depth buffer in the pipeline.

This program (and all of the samples) were built both as a guide to using Vulkan with Go, and as a basic validation and test of the [go-vk
binding](https://github.com/bbredesen/go-vk). Barring any design updates in go-vk as we move towards a 1.0 release
(e.g. struct member naming conventions), this project should not have any material changes.

## Building

Before running, be sure to `go generate` in the project root. This will compile the shaders with glslc.exe. That
compiler is bundled with the Vulkan SDK. The go:generate annotations can be found in main.go.
