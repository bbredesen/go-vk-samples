package main

import (
	"github.com/bbredesen/go-vk"
)

func (app *App_02) copyBuffer(srcBuffer, dstBuffer vk.Buffer, size vk.DeviceSize) {
	cbAllocInfo := vk.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vk.COMMAND_BUFFER_LEVEL_PRIMARY,
		CommandBufferCount: 1,
	}
	bufs, err := vk.AllocateCommandBuffers(app.device, &cbAllocInfo)
	if err != nil {
		panic("Could not allocate command buffer for copy operations: " + err.Error())
	}

	beginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}
	if err := vk.BeginCommandBuffer(bufs[0], &beginInfo); err != nil {
		panic("Could not begin command buffer for copy operations: " + err.Error())
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

func (app *App_02) createBuffer(usage vk.BufferUsageFlags, size vk.DeviceSize, memProps vk.MemoryPropertyFlags) (buffer vk.Buffer, memory vk.DeviceMemory) {

	bufferCI := vk.BufferCreateInfo{
		Size:        size,
		Usage:       usage,
		SharingMode: vk.SHARING_MODE_EXCLUSIVE,
	}

	var err error

	if buffer, err = vk.CreateBuffer(app.device, &bufferCI, nil); err != nil {
		panic("Could not create buffer: " + err.Error())
	}

	memReq := vk.GetBufferMemoryRequirements(app.device, buffer)

	memAllocInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: uint32(app.findMemoryType(memReq.MemoryTypeBits, memProps)), //vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)),
	}

	if memory, err = vk.AllocateMemory(app.device, &memAllocInfo, nil); err != nil {
		panic("Could not allocate memory for buffer: " + err.Error())
	}
	if err := vk.BindBufferMemory(app.device, buffer, memory, 0); err != nil {
		panic("Could not bind memory for buffer: " + err.Error())
	}

	return
}
