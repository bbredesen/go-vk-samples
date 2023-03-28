package main

import (
	"github.com/bbredesen/go-vk"
)

func (app *App_04) copyBuffer(srcBuffer, dstBuffer vk.Buffer, size vk.DeviceSize) {
	cbuf := app.beginSingleTimeCommands()

	region := vk.BufferCopy{
		SrcOffset: 0,
		DstOffset: 0,
		Size:      size,
	}

	vk.CmdCopyBuffer(cbuf, srcBuffer, dstBuffer, []vk.BufferCopy{region})

	app.endSingleTimeCommands(cbuf)
}

func (app *App_04) createBuffer(usage vk.BufferUsageFlags, size vk.DeviceSize, memProps vk.MemoryPropertyFlags) (buffer vk.Buffer, memory vk.DeviceMemory) {

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
