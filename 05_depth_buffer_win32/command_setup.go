package main

import (
	"github.com/bbredesen/go-vk"
)

// Create command pool, associated command buffers, and record commands to clear
// the screen.
func (app *App_05) createCommandPool() {
	// 1) Create the command pool
	poolCreateInfo := vk.CommandPoolCreateInfo{
		Flags:            vk.COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT,
		QueueFamilyIndex: app.presentQueueFamilyIndex,
	}
	r, commandPool := vk.CreateCommandPool(app.device, &poolCreateInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create command pool! " + r.String())
	}
	app.commandPool = commandPool

	// 2) Allocate command buffers from the pool
	allocInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: uint32(len(app.swapchainImages)),
	}
	r, commandBuffers := vk.AllocateCommandBuffers(app.device, &allocInfo)
	if r != vk.SUCCESS {
		panic("Could not allocate command buffers! " + r.String())
	}
	app.commandBuffers = commandBuffers
}

func (app *App_05) destroyCommandPool() {
	vk.FreeCommandBuffers(app.device, app.commandPool, app.commandBuffers)
	vk.DestroyCommandPool(app.device, app.commandPool, nil)
}

func (app *App_05) createSyncObjects() {
	createInfo := vk.SemaphoreCreateInfo{}

	r, imgSem := vk.CreateSemaphore(app.device, &createInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create semaphore! " + r.String())
	}
	app.imageAvailableSemaphore = imgSem

	r, renSem := vk.CreateSemaphore(app.device, &createInfo, nil)
	if r != vk.SUCCESS {
		panic("Could not create semaphore! " + r.String())
	}
	app.renderFinishedSemaphore = renSem

	fenceCreateInfo := vk.FenceCreateInfo{
		Flags: vk.FENCE_CREATE_SIGNALED_BIT,
	}
	if r, app.inFlightFence = vk.CreateFence(app.device, &fenceCreateInfo, nil); r != vk.SUCCESS {
		panic("Could not create fence! " + r.String())
	}
}

func (app *App_05) destroySyncObjects() {
	vk.DestroyFence(app.device, app.inFlightFence, nil)

	vk.DestroySemaphore(app.device, app.imageAvailableSemaphore, nil)
	vk.DestroySemaphore(app.device, app.renderFinishedSemaphore, nil)
}

func (app *App_05) recordCommandBuffer(buffer vk.CommandBuffer, imageIndex uint32) {
	cbBeginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	if r := vk.BeginCommandBuffer(buffer, &cbBeginInfo); r != vk.SUCCESS {
		panic("Could not begin command buffer recording! " + r.String())
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

	dv := vk.ClearValue{}
	dv.AsDepthStencil(vk.ClearDepthStencilValue{Depth: 1.0, Stencil: 0.0})

	rpBeginInfo.PClearValues = append(rpBeginInfo.PClearValues, cv, dv)

	vk.CmdBindPipeline(buffer, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipeline)

	vk.CmdBindVertexBuffers(buffer, 0, []vk.Buffer{app.vertexBuffer}, []vk.DeviceSize{0})
	vk.CmdBindIndexBuffer(buffer, app.indexBuffer, 0, vk.INDEX_TYPE_UINT16)
	vk.CmdBindDescriptorSets(buffer, vk.PIPELINE_BIND_POINT_GRAPHICS, app.pipelineLayout, 0, []vk.DescriptorSet{app.descriptorSets[app.currentImage]}, nil)

	vk.CmdBeginRenderPass(buffer, &rpBeginInfo, vk.SUBPASS_CONTENTS_INLINE)

	// vk.CmdDraw(buffer, uint32(len(verts)), 1, 0, 0)
	vk.CmdDrawIndexed(buffer, uint32(len(indices)), 1, 0, 0, 0)

	vk.CmdEndRenderPass(buffer)

	if r := vk.EndCommandBuffer(buffer); r != vk.SUCCESS {
		panic("Could not end command buffer recording! " + r.String())
	}
}

func (app *App_05) beginSingleTimeCommands() vk.CommandBuffer {
	bufferAlloc := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}

	var r vk.Result
	var bufs []vk.CommandBuffer

	if r, bufs = vk.AllocateCommandBuffers(app.device, &bufferAlloc); r != vk.SUCCESS {
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

func (app *App_05) endSingleTimeCommands(buf vk.CommandBuffer) {
	if r := vk.EndCommandBuffer(buf); r != vk.SUCCESS {
		panic("Could not end one-time command buffer: " + r.String())
	}

	submitInfo := vk.SubmitInfo{
		PCommandBuffers: []vk.CommandBuffer{buf},
	}

	if r := vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, vk.Fence(vk.NULL_HANDLE)); r != vk.SUCCESS {
		panic("Could not submit one-time command buffer: " + r.String())
	}
	if r := vk.QueueWaitIdle(app.graphicsQueue); r != vk.SUCCESS {
		panic("QueueWaitIdle failed: " + r.String())
	}

	vk.FreeCommandBuffers(app.device, app.commandPool, []vk.CommandBuffer{buf})
}
