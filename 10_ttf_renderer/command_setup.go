package main

import (
	"github.com/bbredesen/go-vk"
)

// Create command pool, associated command buffers, and record commands to clear
// the screen.
func (app *App_10) createCommandPool() {
	// 1) Create the command pool
	poolCreateInfo := vk.CommandPoolCreateInfo{
		Flags:            vk.COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
		QueueFamilyIndex: app.presentQueueFamilyIndex,
	}
	commandPool, err := vk.CreateCommandPool(app.device, &poolCreateInfo, nil)
	if err != nil {
		panic("Could not create command pool! " + err.Error())
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
		panic("Could not allocate command buffers! " + err.Error())
	}
	app.commandBuffers = commandBuffers
}

func (app *App_10) destroyCommandPool() {
	vk.FreeCommandBuffers(app.device, app.commandPool, app.commandBuffers)
	vk.DestroyCommandPool(app.device, app.commandPool, nil)
}

func (app *App_10) createSyncObjects() {
	createInfo := vk.SemaphoreCreateInfo{}

	imgSem, err := vk.CreateSemaphore(app.device, &createInfo, nil)
	if err != nil {
		panic("Could not create semaphore! " + err.Error())
	}
	app.imageAvailableSemaphore = imgSem

	renSem, err := vk.CreateSemaphore(app.device, &createInfo, nil)
	if err != nil {
		panic("Could not create semaphore! " + err.Error())
	}
	app.renderFinishedSemaphore = renSem

	fenceCreateInfo := vk.FenceCreateInfo{
		Flags: vk.FENCE_CREATE_SIGNALED_BIT,
	}
	if app.inFlightFence, err = vk.CreateFence(app.device, &fenceCreateInfo, nil); err != nil {
		panic("Could not create fence! " + err.Error())
	}
}

func (app *App_10) destroySyncObjects() {
	vk.DestroyFence(app.device, app.inFlightFence, nil)

	vk.DestroySemaphore(app.device, app.imageAvailableSemaphore, nil)
	vk.DestroySemaphore(app.device, app.renderFinishedSemaphore, nil)
}

func (app *App_10) recordCommandBuffer(cb vk.CommandBuffer, imageIndex uint32) {
	cbBeginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	rpBeginInfo := vk.RenderPassBeginInfo{
		RenderPass:  app.renderPass,
		Framebuffer: app.swapChainFramebuffers[imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: app.swapchainExtent,
		},
	}

	colorCV, stencilCV := vk.ClearValue{}, vk.ClearValue{}
	ccv := vk.ClearColorValue{}
	ccv.AsTypeFloat32([4]float32{0, 0, 0, 1})
	colorCV.AsColor(ccv)

	stencilCV.AsDepthStencil(vk.ClearDepthStencilValue{
		Depth:   0,
		Stencil: 0,
	})

	rpBeginInfo.PClearValues = append(rpBeginInfo.PClearValues, colorCV, stencilCV)

	if err := vk.BeginCommandBuffer(cb, &cbBeginInfo); err != nil {
		panic("Could not begin command buffer recording! " + err.Error())
	}

	vk.CmdBeginRenderPass(cb, &rpBeginInfo, vk.SUBPASS_CONTENTS_INLINE)

	// bind vert, index bufs
	vk.CmdBindVertexBuffers(cb, 0, []vk.Buffer{app.vertexBuffer}, []vk.DeviceSize{0})
	vk.CmdBindIndexBuffer(cb, app.indexBuffer, 0, vk.INDEX_TYPE_UINT16)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[0]) // stencil pipeline
	vk.CmdDrawIndexed(cb, uint32(app.quadIndsStart), 1, 0, 0, 0)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[1]) // stencil quad portion pipeline
	vk.CmdDrawIndexed(cb, uint32(app.indexCount-app.quadIndsStart)-4, 1, uint32(app.quadIndsStart), int32(app.quadVertStart), 0)

	vk.CmdNextSubpass(cb, vk.SUBPASS_CONTENTS_INLINE)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[2]) // Color pass

	// vk.CmdDrawIndexed(cb, 5, 1, uint32(app.indexCount)-5, 0, 0)
	vk.CmdDrawIndexed(cb, 4, 1, uint32(app.indexCount)-4, int32(app.quadVertStart), 0)

	// draw

	vk.CmdEndRenderPass(cb)

	if err := vk.EndCommandBuffer(cb); err != nil {
		panic("Could not end command buffer recording! " + err.Error())
	}
}

func (app *App_10) beginSingleTimeCommands() vk.CommandBuffer {
	bufferAlloc := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}

	var err error
	var bufs []vk.CommandBuffer

	if bufs, err = vk.AllocateCommandBuffers(app.device, &bufferAlloc); err != nil {
		panic("Could not allocate one-time command buffer: " + err.Error())
	}

	cbbInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	if err := vk.BeginCommandBuffer(bufs[0], &cbbInfo); err != nil {
		panic("Could not begin recording one-time command buffer: " + err.Error())
	}

	return bufs[0]
}

func (app *App_10) endSingleTimeCommands(buf vk.CommandBuffer) {
	if err := vk.EndCommandBuffer(buf); err != nil {
		panic("Could not end one-time command buffer: " + err.Error())
	}

	submitInfo := vk.SubmitInfo{
		PCommandBuffers: []vk.CommandBuffer{buf},
	}

	if err := vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, vk.Fence(vk.NULL_HANDLE)); err != nil {
		panic("Could not submit one-time command buffer: " + err.Error())
	}
	if err := vk.QueueWaitIdle(app.graphicsQueue); err != nil {
		panic("QueueWaitIdle failed: " + err.Error())
	}

	vk.FreeCommandBuffers(app.device, app.commandPool, []vk.CommandBuffer{buf})
}
