package main

import (
	"github.com/bbredesen/go-vk"
)

func (app *App_03) copyBuffer(srcBuffer, dstBuffer vk.Buffer, size vk.DeviceSize) {
	cbAllocInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}
	r, bufs := vk.AllocateCommandBuffers(app.device, &cbAllocInfo)
	if r != vk.SUCCESS {
		panic("Could not allocate command buffer for copy operations: " + r.String())
	}

	beginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}
	if r := vk.BeginCommandBuffer(bufs[0], &beginInfo); r != vk.SUCCESS {
		panic("Could not begin command buffer for copy operations: " + r.String())
	}

	region := vk.BufferCopy{
		SrcOffset: 0,
		DstOffset: 0,
		Size:      size,
	}

	vk.CmdCopyBuffer(bufs[0], srcBuffer, dstBuffer, []vk.BufferCopy{region})

	vk.EndCommandBuffer(bufs[0])

	submitInfo := vk.SubmitInfo{
		PCommandBuffers: bufs,
	}
	vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, vk.Fence(vk.NULL_HANDLE))
	vk.QueueWaitIdle(app.graphicsQueue)

	vk.FreeCommandBuffers(app.device, app.commandPool, bufs)
}

func (app *App_03) createBuffer(usage vk.BufferUsageFlags, size vk.DeviceSize, memProps vk.MemoryPropertyFlags) (buffer vk.Buffer, memory vk.DeviceMemory) {

	bufferCI := vk.BufferCreateInfo{
		Size:        size,
		Usage:       usage,
		SharingMode: vk.SHARING_MODE_EXCLUSIVE,
	}

	var r vk.Result

	if r, buffer = vk.CreateBuffer(app.device, &bufferCI, nil); r != vk.SUCCESS {
		panic("Could not create buffer: " + r.String())
	}

	memReq := vk.GetBufferMemoryRequirements(app.device, buffer)

	memAllocInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: uint32(app.findMemoryType(memReq.MemoryTypeBits, memProps)), //vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)),
	}

	if r, memory = vk.AllocateMemory(app.device, &memAllocInfo, nil); r != vk.SUCCESS {
		panic("Could not allocate memory for buffer: " + r.String())
	}
	if r := vk.BindBufferMemory(app.device, buffer, memory, 0); r != vk.SUCCESS {
		panic("Could not bind memory for buffer: " + r.String())
	}

	return
}
