package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/vkm"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const buffer_size = 4096

func (app *App) loadBuffers(segments sfnt.Segments, bounds fixed.Rectangle26_6) {
	verts, inds, quadVerts, quadInds := convertSegmentsToVerts(segments, bounds)

	stagingBuffer, stagingMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, buffer_size, vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)

	app.vertexBuffer, app.vertexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_VERTEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, buffer_size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	app.indexBuffer, app.indexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_INDEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, buffer_size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.quadVertStart = len(verts)
	app.quadIndsStart = len(inds)

	verts = append(verts, quadVerts...)
	inds = append(inds, quadInds...)

	app.indexCount = len(inds)

	r, ptr := vk.MapMemory(app.Device, stagingMemory, 0, buffer_size, 0)
	if r != vk.SUCCESS {
		panic(r)
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), inds)
	app.copyBuffer(stagingBuffer, app.indexBuffer, buffer_size)

	vk.MemCopySlice(unsafe.Pointer(ptr), verts)
	app.copyBuffer(stagingBuffer, app.vertexBuffer, buffer_size)

	vk.UnmapMemory(app.Device, stagingMemory)
	vk.DestroyBuffer(app.Device, stagingBuffer, nil)
	vk.FreeMemory(app.Device, stagingMemory, nil)

}

func int26_6_to_float32(x fixed.Int26_6) float32 {
	return float32(x>>6) + (float32(x&0x7f) / float32(0x7f))
}

type vertexFormat struct {
	position   vkm.Pt2
	baryCoords vkm.Pt3
}

func convertSegmentsToVerts(segments sfnt.Segments, bounds fixed.Rectangle26_6) (verts []vertexFormat, inds []uint16, quadVerts []vertexFormat, quadInds []uint16) {
	// verts = append(verts, vkm.Origin2())

	barySign := 0

	getBaryCoord := func() vkm.Pt3 {
		switch barySign {
		case -1:
			barySign *= -1
			return vkm.Pt3{0, 0, 1}
		case 0:
			barySign = 1
			return vkm.Pt3{0, 1, 0}
		case 1:
			barySign *= -1
			return vkm.Pt3{1, 0, 0}
		}
		panic("unexpected barySign ")
	}

	pt2FromFixed := func(fp fixed.Point26_6) vkm.Pt2 {
		return vkm.Pt2{int26_6_to_float32(fp.X), int26_6_to_float32(fp.Y)}
	}
	pushRestart := func() {
		inds = append(inds, 0xFFFF)
		barySign = 0
	}

	pushVertex := func(fp fixed.Point26_6) {
		nextIdx := uint16(len(verts))

		pt := pt2FromFixed(fp)
		verts = append(verts, vertexFormat{pt, getBaryCoord()})
		inds = append(inds, nextIdx)
	}

	// Segments is a list of movement instructions
	// OpCode MoveTo - Restart primitive and use arg[0] as the first point
	// OpCode QuadTO - Quadratic curve to arg[1], arg[0] is the control point
	// OpCode CubeTo - Cubic curve - Not supported

	for _, segment := range segments {
		switch segment.Op {
		case sfnt.SegmentOpMoveTo:
			pushRestart()
			pushVertex(segment.Args[0])

		case sfnt.SegmentOpLineTo:
			pushVertex(segment.Args[0])

		case sfnt.SegmentOpQuadTo:
			pushVertex(segment.Args[1]) // push for rough rendering triangle fans

			vlen := len(verts)
			// for each quad, need to push last point, this point, control point, with bary coords
			v0, v1 := verts[vlen-2], verts[vlen-1]
			v0.baryCoords = vkm.Pt3{1, 0, 0}
			v1.baryCoords = vkm.Pt3{0, 0, 1}

			qvIdxStart := uint16(len(quadVerts))

			quadVerts = append(quadVerts, v0,
				vertexFormat{
					position:   pt2FromFixed(segment.Args[0]),
					baryCoords: vkm.Pt3{0, 1, 0},
				},
				v1,
			)
			quadInds = append(quadInds, qvIdxStart, qvIdxStart+1, qvIdxStart+2)
		}
	}

	// pushRestart()
	sidx := uint16(len(quadVerts))

	minX, maxX := int26_6_to_float32(bounds.Min.X), int26_6_to_float32(bounds.Max.X)
	minY, maxY := int26_6_to_float32(bounds.Min.Y), int26_6_to_float32(bounds.Max.Y)

	logrus.WithFields(logrus.Fields{
		"minX": minX,
		"minY": minY,
		"maxX": maxX,
		"maxY": maxY,
	}).Infof("Bounds")

	// Set uniform buffer
	quadVerts = append(quadVerts, //vkm.Pt2{20, 0}, vkm.Pt2{20, -20}, vkm.Pt2{0, -20})
		vertexFormat{vkm.Pt2{minX, minY}, vkm.Origin3()},
		vertexFormat{vkm.Pt2{minX, maxY}, vkm.Origin3()},
		vertexFormat{vkm.Pt2{maxX, maxY}, vkm.Origin3()},
		vertexFormat{vkm.Pt2{maxX, minY}, vkm.Origin3()},
	)
	quadInds = append(quadInds, sidx, sidx+1, sidx+2, sidx+3)

	return
}

func (app *App) destroyBuffers() {
	vk.DestroyBuffer(app.Device, app.indexBuffer, nil)
	vk.FreeMemory(app.Device, app.indexBufferMemory, nil)

	vk.DestroyBuffer(app.Device, app.vertexBuffer, nil)
	vk.FreeMemory(app.Device, app.vertexBufferMemory, nil)

}

func (app *App) copyBuffer(srcBuffer, dstBuffer vk.Buffer, size vk.DeviceSize) {
	cbuf := app.BeginOneTimeCommands()

	region := vk.BufferCopy{
		SrcOffset: 0,
		DstOffset: 0,
		Size:      size,
	}

	vk.CmdCopyBuffer(cbuf, srcBuffer, dstBuffer, []vk.BufferCopy{region})

	app.EndOneTimeCommands(cbuf)
}

func (app *App) createBuffer(usage vk.BufferUsageFlags, size vk.DeviceSize, memProps vk.MemoryPropertyFlags) (buffer vk.Buffer, memory vk.DeviceMemory) {

	bufferCI := vk.BufferCreateInfo{
		Size:        size,
		Usage:       usage,
		SharingMode: vk.SHARING_MODE_EXCLUSIVE,
	}

	var r vk.Result

	if r, buffer = vk.CreateBuffer(app.Device, &bufferCI, nil); r != vk.SUCCESS {
		panic("Could not create buffer: " + r.String())
	}

	memReq := vk.GetBufferMemoryRequirements(app.Device, buffer)

	memAllocInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: uint32(app.FindMemoryType(memReq.MemoryTypeBits, memProps)), //vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)),
	}

	if r, memory = vk.AllocateMemory(app.Device, &memAllocInfo, nil); r != vk.SUCCESS {
		panic("Could not allocate memory for buffer: " + r.String())
	}
	if r := vk.BindBufferMemory(app.Device, buffer, memory, 0); r != vk.SUCCESS {
		panic("Could not bind memory for buffer: " + r.String())
	}

	return
}
