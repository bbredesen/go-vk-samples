# 04_texture_mapping_win32
![](screenshot.png?raw=true)

This basic Vulkan application follows the fourth chapter (Texture Mapping) of the excellent [Vulkan
Tutorial](https://vulkan-tutorial.com), building on the application at `03_uniform_buffers`. It demonstrates:

* Loading the texture file, which is significantly different in Go from the C library used by the original tutorial
* Loading and transitioning the texture data to the GPU
* Creating an image sampler for the shader to use when accessing the texture
* Modifying the pipeline and shader to use the texture

This program (and all of the samples) were built both as a guide to using Vulkan with Go, and as a basic validation and test of the [go-vk
binding](https://github.com/bbredesen/go-vk). Barring any design updates in go-vk as we move towards a 1.0 release
(e.g. struct member naming conventions), this project should not have any material changes.

## Building

Before running, be sure to `go generate` in the project root. This will compile the shaders with glslc.exe. That
compiler is bundled with the Vulkan SDK. The go:generate annotations can be found in main.go.
