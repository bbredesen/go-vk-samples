package vkctx

import (
	"github.com/bbredesen/go-vk"
)

// Create command pool, associated command buffers, and record commands to clear
// the screen.
func (ctx *Context) createCommandPool() {
	// 1) Create the command pool
	poolCreateInfo := vk.CommandPoolCreateInfo{
		Flags:            vk.COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
		QueueFamilyIndex: ctx.PresentQueueFamilyIndex,
	}
	r, commandPool := vk.CreateCommandPool(ctx.Device, &poolCreateInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create command pool! " + r.String())
	}
	ctx.CommandPool = commandPool

	// 2) Allocate primary command buffers, one for each swapchain image, from the pool
	allocInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        ctx.CommandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: uint32(len(ctx.SwapchainImages)),
	}
	r, commandBuffers := vk.AllocateCommandBuffers(ctx.Device, &allocInfo)
	if r != vk.SUCCESS {
		panic("Could not allocate command buffers! " + r.String())
	}
	ctx.CommandBuffers = commandBuffers
}

func (ctx *Context) destroyCommandPool() {
	vk.FreeCommandBuffers(ctx.Device, ctx.CommandPool, ctx.CommandBuffers)
	vk.DestroyCommandPool(ctx.Device, ctx.CommandPool, nil)
}

func (ctx *Context) createSyncObjects() {
	createInfo := vk.SemaphoreCreateInfo{}

	r, imgSem := vk.CreateSemaphore(ctx.Device, &createInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create semaphore! " + r.String())
	}
	ctx.ImageAvailableSemaphore = imgSem

	r, renSem := vk.CreateSemaphore(ctx.Device, &createInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create semaphore! " + r.String())
	}
	ctx.RenderFinishedSemaphore = renSem

	fenceCreateInfo := vk.FenceCreateInfo{
		Flags: vk.FENCE_CREATE_SIGNALED_BIT,
	}
	if r, ctx.InFlightFence = vk.CreateFence(ctx.Device, &fenceCreateInfo, nil); r != vk.SUCCESS {
		panic("Could not create fence! " + r.String())
	}
}

func (app *Context) destroySyncObjects() {
	vk.DestroyFence(app.Device, app.InFlightFence, nil)

	vk.DestroySemaphore(app.Device, app.ImageAvailableSemaphore, nil)
	vk.DestroySemaphore(app.Device, app.RenderFinishedSemaphore, nil)
}

func (app *Context) BeginOneTimeCommands() vk.CommandBuffer {
	bufferAlloc := vk.CommandBufferAllocateInfo{
		CommandPool:        app.CommandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}

	var r vk.Result
	var bufs []vk.CommandBuffer

	if r, bufs = vk.AllocateCommandBuffers(app.Device, &bufferAlloc); r != vk.SUCCESS {
		panic("Could not allocate one-time command buffer: " + r.String())
	}

	cbbInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	if r := vk.BeginCommandBuffer(bufs[0], &cbbInfo); r != vk.SUCCESS {
		panic("Could not begin recording one-time command buffer: " + r.String())
	}

	return bufs[0]
}

func (app *Context) EndOneTimeCommands(buf vk.CommandBuffer) {
	if r := vk.EndCommandBuffer(buf); r != vk.SUCCESS {
		panic("Could not end one-time command buffer: " + r.String())
	}

	submitInfo := vk.SubmitInfo{
		PCommandBuffers: []vk.CommandBuffer{buf},
	}

	if r := vk.QueueSubmit(app.GraphicsQueue, []vk.SubmitInfo{submitInfo}, vk.Fence(vk.NULL_HANDLE)); r != vk.SUCCESS {
		panic("Could not submit one-time command buffer: " + r.String())
	}
	if r := vk.QueueWaitIdle(app.GraphicsQueue); r != vk.SUCCESS {
		panic("QueueWaitIdle failed: " + r.String())
	}

	vk.FreeCommandBuffers(app.Device, app.CommandPool, []vk.CommandBuffer{buf})
}
