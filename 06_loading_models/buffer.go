package main

import (
	"github.com/bbredesen/go-vk"
)

func (app *App_06) copyBuffer(srcBuffer, dstBuffer vk.Buffer, size vk.DeviceSize) {
	cbuf := app.beginSingleTimeCommands()

	region := vk.BufferCopy{
		SrcOffset: 0,
		DstOffset: 0,
		Size:      size,
	}

	vk.CmdCopyBuffer(cbuf, srcBuffer, dstBuffer, []vk.BufferCopy{region})

	app.endSingleTimeCommands(cbuf)
}

func (app *App_06) createBuffer(usage vk.BufferUsageFlags, size vk.DeviceSize, fl vk.MemoryPropertyFlags) (buffer vk.Buffer, memory vk.DeviceMemory) {

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
	// time.Sleep(1 * time.Second)

	// if size > 1024 {
	// fmt.Printf("memreq: %x\n", memReq.MemoryTypeBits)
	// _ = fl
	// }

	memAllocInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: app.findMemoryType(memReq.MemoryTypeBits, fl), //vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)),
	}

	if r, memory = vk.AllocateMemory(app.device, &memAllocInfo, nil); r != vk.SUCCESS {
		panic("Could not allocate memory for buffer: " + r.String())
	}
	if r := vk.BindBufferMemory(app.device, buffer, memory, 0); r != vk.SUCCESS {
		panic("Could not bind memory for buffer: " + r.String())
	}

	return
}
