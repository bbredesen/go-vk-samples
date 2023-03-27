package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/vkm"
)

type Vertex struct {
	Pos      vkm.Pt3
	Color    vkm.Vec3
	TexCoord vkm.Pt2
}

// Note: The texture coordinates are different here than in the Tutorial. In Vulkan, textures coordinates run from
// [0..1] going left to right and going top to bottom on the texture image. If we assume world-space Y is supposed to be up on the
// rendered image, then u,v = 0,0 should be on the -x, +y point of the square. The way the Tutorial is written, the
// image will be mirrored when rendered, because some of the texture coordinates are inverted.
var verts = []Vertex{

	{vkm.Pt3{-0.5, -0.5, +0.0}, vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.Pt3{+0.5, -0.5, +0.0}, vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{1.0, 0.0}},
	{vkm.Pt3{+0.5, +0.5, +0.0}, vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{1.0, 1.0}},
	{vkm.Pt3{-0.5, +0.5, +0.0}, vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{0.0, 1.0}},

	{vkm.Pt3{-0.5, -0.5, -0.5}, vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.Pt3{+0.5, -0.5, -0.5}, vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{1.0, 0.0}},
	{vkm.Pt3{+0.5, +0.5, -0.5}, vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{1.0, 1.0}},
	{vkm.Pt3{-0.5, +0.5, -0.5}, vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{0.0, 1.0}},

	{vkm.NewPt3(-0.5, -0.5, +0.0), vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{1.0, 0.0}},
	{vkm.NewPt3(+0.5, -0.5, +0.0), vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.NewPt3(+0.5, +0.5, +0.0), vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{0.0, 1.0}},
	{vkm.NewPt3(-0.5, +0.5, +0.0), vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{1.0, 1.0}},

	{vkm.NewPt3(-0.5, -0.5, -0.5), vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.NewPt3(+0.5, -0.5, -0.5), vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{1.5, 0.0}},
	{vkm.NewPt3(+0.5, +0.5, -0.5), vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{1.5, 1.5}},
	{vkm.NewPt3(-0.5, +0.5, -0.5), vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{0.0, 1.5}},

	{vkm.NewPt3(-0.5+0.25, -0.5, -1.5), vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.NewPt3(+0.5+0.25, -0.5, -1.5), vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{1.0, 0.0}},
	{vkm.NewPt3(+0.5+0.25, +0.5, -1.5), vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{1.0, 1.0}},
	{vkm.NewPt3(-0.5+0.25, +0.5, -1.5), vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{0.0, 1.0}},

	{vkm.NewPt3(-0.5, -0.5, -2.5), vkm.Vec3{1.0, 0.0, 0.0}, vkm.Pt2{0.0, 0.0}},
	{vkm.NewPt3(+0.5, -0.5, -2.5), vkm.Vec3{0.0, 1.0, 0.0}, vkm.Pt2{1.0, 0.0}},
	{vkm.NewPt3(+0.5, +0.5, -2.5), vkm.Vec3{0.0, 0.0, 1.0}, vkm.Pt2{1.0, 1.0}},
	{vkm.NewPt3(-0.5, +0.5, -2.5), vkm.Vec3{1.0, 1.0, 1.0}, vkm.Pt2{0.0, 1.0}},
}

var indices = []uint16{
	0, 1, 2, 2, 3, 0,
	4, 5, 6, 6, 7, 4,
	// 8, 9, 10, 10, 11, 8,
}

func getVertexBindingDescription() vk.VertexInputBindingDescription {
	return vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    uint32(unsafe.Sizeof(Vertex{})),
		InputRate: vk.VERTEX_INPUT_RATE_VERTEX,
	}
}

func getVertexAttributeDescriptions() []vk.VertexInputAttributeDescription {
	rval := make([]vk.VertexInputAttributeDescription, 3)

	rval[0] = vk.VertexInputAttributeDescription{
		Location: 0,
		Binding:  0,
		Format:   vk.FORMAT_R32G32B32_SFLOAT,
		Offset:   uint32(unsafe.Offsetof(Vertex{}.Pos)),
	}
	rval[1] = vk.VertexInputAttributeDescription{
		Location: 1,
		Binding:  0,
		Format:   vk.FORMAT_R32G32B32_SFLOAT,
		Offset:   uint32(unsafe.Offsetof(Vertex{}.Color)),
	}
	rval[2] = vk.VertexInputAttributeDescription{
		Location: 2,
		Binding:  0,
		Format:   vk.FORMAT_R32G32_SFLOAT,
		Offset:   uint32(unsafe.Offsetof(Vertex{}.TexCoord)),
	}

	return rval
}

func (app *App_05) createVertexBuffer() {

	size := vk.DeviceSize(unsafe.Sizeof(Vertex{})) * vk.DeviceSize(len(verts))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	ptr, err := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if err != nil {
		panic("Could not map memory for vertex buffer: " + err.Error())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), verts)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.vertexBuffer, app.vertexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_VERTEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.vertexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_05) createIndexBuffer() {

	size := vk.DeviceSize(unsafe.Sizeof(uint16(0))) * vk.DeviceSize(len(indices))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	ptr, err := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if err != nil {
		panic("Could not map memory for vertex buffer: " + err.Error())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), indices)

	// // MapMemory needs to return a []byte maybe...shouldn't have to do this:
	// var sl = struct {
	// 	addr uintptr
	// 	len  int
	// 	cap  int
	// }{uintptr(unsafe.Pointer(ptr)), int(size), int(size)}
	// bytes := *(*[]byte)(unsafe.Pointer(&sl))

	// vb := AnySliceToBytes(indices)
	// copy(bytes, vb)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.indexBuffer, app.indexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_INDEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.indexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_05) findMemoryType(typeFilter uint32, flags vk.MemoryPropertyFlags) uint32 {
	memProps := vk.GetPhysicalDeviceMemoryProperties(app.physicalDevice)
	var i uint32
	for i = 0; i < memProps.MemoryTypeCount; i++ {
		if typeFilter&1<<i != 0 && memProps.MemoryTypes[i].PropertyFlags&flags == flags {
			return i
		}
	}
	panic("Could not find appropriate memory type.")
}

// // VERY MUCH TODO: This or something similar should be a convenience function provided with go-vk
// func AnySliceToBytes[T any](input []T) []byte {
// 	type sl struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}

// 	if len(input) == 0 {
// 		return []byte{}
// 	}

// 	inputLen := len(input) * int(unsafe.Sizeof(input[0]))
// 	sl_input := sl{uintptr(unsafe.Pointer(&input[0])), inputLen, inputLen}

// 	return *(*[]byte)(unsafe.Pointer(&sl_input))
// }

// func AnyTypeToBytes[T any](input T) []byte {
// 	type sl struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}
// 	sl_input := sl{uintptr(unsafe.Pointer(&input)), int(unsafe.Sizeof(input)), int(unsafe.Sizeof(input))}
// 	return *(*[]byte)(unsafe.Pointer(&sl_input))

// }
