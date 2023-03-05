# go-vk-samples

This repository contains runnable sample apps to demonstrate use of [go-vk](https://github.com/bbredesen/go-vk), a
Go-language binding for the Vulkan graphics and compute API. Samples currently only run under Windows.

The "0n_" series of samples follow the project found at [Vulkan-Tutorial.com](https://vulkan-tutorial.com).

* 01_single_triangle - General Vulkan setup and rendering a single static triangle
* 02_vertex_buffers - Shader input bindings, indexed rendering, staging copying data from the CPU to the GPU
* 03_uniform_buffers - Descriptor sets, uniform buffers, MVP projections
* 04_texture_mapping
* 05_depth_buffering
* 06_loading_models - Loading Wavefront OBJ models from a file
* 07_generating_mipmaps - Generating, binding and using mipmap textures - TODO - Not yet implemented
* 08_multisampling - TODO - Not yet implemented
* instanced_rendering - TODO - Not yet implemented - Instanced model, repeated N times at different positions, scales,
  etc. (4k instances?)

The "shared" directory contains a basic Win32 app structure, basically just enough to get a window open and ready to
attach a Vulkan surface to. The original intent was to allow for a common event passing structure, so that these samples
could be set up on Darwin, XWindows, etc. Doing that is still on the plan, but will require significant rework.

## Other samples and pending projects:

- Glyph rendering from TTF using stencils (https://github.com/bbredesen/ttf-renderer)
- glTF model viewer - (https://github.com/bbredesen/gltf-viewer) - WIP
- Performance test - Rendering 10K, 100K, 1M, 10M poly scenes?
- Extension demos and testing
- Compute shaders - Possibly preculling a large scene, or post-processing graphics pipeline output.
