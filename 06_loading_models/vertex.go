package main

import (
	"fmt"
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/udhos/gwob"
)

// The Vertex struct goes away, since gwob outputs a single slice with position, texture, and normal coordinates. The
// underlying memory gets used directly and we just have to update the vertex binding and attribute descriptions.
// type Vertex struct {
// 	Pos      vkm.Pt3
// 	TexCoord vkm.Pt2
// 	Normal   vkm.Vec3
// }

var obj *gwob.Obj
var indices []uint32

func getVertexBindingDescription() vk.VertexInputBindingDescription {
	return vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    uint32(obj.StrideSize),
		InputRate: vk.VERTEX_INPUT_RATE_VERTEX,
	}
}

func getVertexAttributeDescriptions() []vk.VertexInputAttributeDescription {
	rval := make([]vk.VertexInputAttributeDescription, 3)

	rval[0] = vk.VertexInputAttributeDescription{
		Location: 0,
		Binding:  0,
		Format:   vk.FORMAT_R32G32B32_SFLOAT,
		Offset:   uint32(obj.StrideOffsetPosition),
	}
	rval[1] = vk.VertexInputAttributeDescription{
		Location: 1,
		Binding:  0,
		Format:   vk.FORMAT_R32G32B32_SFLOAT,
		Offset:   uint32(obj.StrideOffsetNormal),
	}
	rval[2] = vk.VertexInputAttributeDescription{
		Location: 2,
		Binding:  0,
		Format:   vk.FORMAT_R32G32_SFLOAT,
		Offset:   uint32(obj.StrideOffsetTexture),
	}

	return rval
}

func (app *App_06) createVertexBuffer() {

	size := vk.DeviceSize(8*unsafe.Sizeof(float32(0))) * vk.DeviceSize(len(obj.Coord))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	ptr, err := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if err != nil {
		panic("Could not map memory for vertex buffer: " + err.Error())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), obj.Coord)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.vertexBuffer, app.vertexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_VERTEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.vertexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_06) createIndexBuffer() {

	size := vk.DeviceSize(unsafe.Sizeof(indices[0])) * vk.DeviceSize(len(indices))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	ptr, err := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if err != nil {
		panic("Could not map memory for vertex buffer: " + err.Error())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), indices)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.indexBuffer, app.indexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_INDEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.indexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_06) findMemoryType(typeFilter uint32, fl vk.MemoryPropertyFlags) uint32 {
	// fmt.Println(fl)
	// _ = fl.String()
	// fmt.Printf("typeFilter was %x, props were %x\n", typeFilter, fl)
	// var memProps vk.PhysicalDeviceMemoryProperties = app.pdmp
	var memProps vk.PhysicalDeviceMemoryProperties
	memProps = vk.GetPhysicalDeviceMemoryProperties(app.physicalDevice)
	memProps = vk.GetPhysicalDeviceMemoryProperties(app.physicalDevice)

	var i uint32
	for i = 0; i < memProps.MemoryTypeCount; i++ {
		if typeFilter&(1<<i) != 0 && (memProps.MemoryTypes[i].PropertyFlags&fl) == fl {
			return i
		}
	}

	// for i = 0; i < app.pdmp.MemoryTypeCount; i++ {
	// 	if typeFilter&(1<<i) != 0 && (app.pdmp.MemoryTypes[i].PropertyFlags&fl) == fl {
	// 		return i
	// 	}
	// }

	// fmt.Printf("memprops: %+v\n", app.pdmp)
	panic(fmt.Sprintf("Could not find appropriate memory type; typeFilter was %x, props were %x, ", typeFilter, fl))
}

func (app *App_06) loadModel() {
	var err error
	if obj, err = gwob.NewObjFromFile("models/viking_room.obj", &gwob.ObjParserOptions{
		LogStats: true,
	}); err != nil {
		panic("Could not load model from file: " + err.Error())
	}

	// Vulkan accepts index sizes only up to 32 bits, but gwob returns them as int, which will be 64 bits on most machines.
	indices = make([]uint32, len(obj.Indices))
	for i, v := range obj.Indices {
		indices[i] = uint32(v)
	}
}
