# go-vk-samples

This repository contains runnable sample apps to demonstrate use of [go-vk](https://github.com/bbredesen/go-vk), a
Go-language binding for the Vulkan graphics and compute API. Samples currently only run under Windows.

Before running, change to the shared/ folder and run `make` to build the OS delegation libary, which the shared package links to as part of `go build` or `go run`.

Run the 00_minimal example first, which will confirm that Vulkan is installed and minimally working on your system, and
that go-vk is linking to it.

The "0n_" series of samples follow the project found at [Vulkan-Tutorial.com](https://vulkan-tutorial.com).

* 01_single_triangle - General Vulkan setup and rendering a single static triangle
* 02_vertex_buffers - Shader input bindings, indexed rendering, staging copying data from the CPU to the GPU
* 03_uniform_buffers - Descriptor sets, uniform buffers, MVP projections
* 04_texture_mapping - Loading a texture from a file and sampling it in the fragment shader.
* 05_depth_buffering - Depth buffering on a pair of textured squares.
* 06_loading_models - Loading Wavefront OBJ models from a file
* 07_generating_mipmaps - Generating, binding and using mipmap textures - TODO - Not yet implemented
* 08_multisampling - TODO - Not yet implemented

The "shared" directory contains a very basic OS-independent app structure, with implementations for Win32 and MacOS/Cocoa.
Operating system events are passed into a channel that the sample applications should empty before drawing each frame. Reading that channel should happen in a separate goroutine, as the windowing code is locked to the main thread, and blocks in App.Run() loop as long as the window is open.

## Other samples and pending projects:

- Glyph rendering from TTF using stencils (https://github.com/bbredesen/ttf-renderer) - Could be migrated to this project
- glTF model viewer - (https://github.com/bbredesen/gltf-viewer) - WIP
