package main

import (
	"github.com/bbredesen/go-vk"
)

// Create command pool, associated command buffers, and record commands to clear
// the screen.
func (app *App_01) createCommandPool() {
	// 1) Create the command pool
	poolCreateInfo := vk.CommandPoolCreateInfo{
		Flags:            vk.COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
		QueueFamilyIndex: app.presentQueueFamilyIndex,
	}
	commandPool, err := vk.CreateCommandPool(app.device, &poolCreateInfo, nil)
	if err != nil {
		panic("Could not create command pool! ")
	}
	app.commandPool = commandPool

	// 2) Allocate command buffers from the pool
	allocInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: uint32(len(app.swapchainImages)),
	}
	commandBuffers, err := vk.AllocateCommandBuffers(app.device, &allocInfo)
	if err != nil {
		panic("Could not allocate command buffers! ")
	}
	app.commandBuffers = commandBuffers
}

func (app *App_01) destroyCommandPool() {
	vk.FreeCommandBuffers(app.device, app.commandPool, app.commandBuffers)
	vk.DestroyCommandPool(app.device, app.commandPool, nil)
}

func (app *App_01) createSyncObjects() {
	createInfo := vk.SemaphoreCreateInfo{}

	imgSem, err := vk.CreateSemaphore(app.device, &createInfo, nil)
	if err != nil {
		panic("Could not create semaphore! ")
	}
	app.imageAvailableSemaphore = imgSem

	renSem, err := vk.CreateSemaphore(app.device, &createInfo, nil)
	if err != nil {
		panic("Could not create semaphore! ")
	}
	app.renderFinishedSemaphore = renSem

	fenceCreateInfo := vk.FenceCreateInfo{
		Flags: vk.FENCE_CREATE_SIGNALED_BIT,
	}
	if app.inFlightFence, err = vk.CreateFence(app.device, &fenceCreateInfo, nil); err != nil {
		panic("Could not create fence! ")
	}
}

func (app *App_01) destroySyncObjects() {
	vk.DestroyFence(app.device, app.inFlightFence, nil)

	vk.DestroySemaphore(app.device, app.imageAvailableSemaphore, nil)
	vk.DestroySemaphore(app.device, app.renderFinishedSemaphore, nil)
}

func (app *App_01) recordCommandBuffer(buffer vk.CommandBuffer, imageIndex uint32) {
	cbBeginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	if err := vk.BeginCommandBuffer(buffer, &cbBeginInfo); err != nil {
		panic("Could not begin command buffer recording! ")
	}

	rpBeginInfo := vk.RenderPassBeginInfo{
		RenderPass:  app.renderPass,
		Framebuffer: app.swapChainFramebuffers[imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: app.swapchainExtent,
		},
	}

	cv, ccv := vk.ClearValue{}, vk.ClearColorValue{}
	ccv.AsTypeFloat32([4]float32{0.0, 0.0, 0.0, 1.0})
	cv.AsColor(ccv)

	rpBeginInfo.PClearValues = append(rpBeginInfo.PClearValues, cv)

	vk.CmdBindPipeline(buffer, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipeline)

	vk.CmdBeginRenderPass(buffer, &rpBeginInfo, vk.SUBPASS_CONTENTS_INLINE)

	vk.CmdDraw(buffer, 3, 1, 0, 0)

	vk.CmdEndRenderPass(buffer)

	if err := vk.EndCommandBuffer(buffer); err != nil {
		panic("Could not end command buffer recording! ")
	}
}
