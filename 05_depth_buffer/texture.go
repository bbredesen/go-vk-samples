package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"unsafe"

	"github.com/bbredesen/go-vk"
)

// loadImage opens the hard coded jpeg file, and reads into a 4-component RGB byte slice, equivalent to
// vk.FORMAT_R8G8B8A8_UINT
//
// JPEGs don't support alpha, of course, so you could safely use vk.FORMAT_R8G8B8_UINT, but we'll follow the steps from
// the Tutorial as closely as possible.
func loadImage() ([]byte, uint32, uint32, vk.Format) {

	fname := "textures/texture.jpg"

	f, err := os.Open(fname)
	defer f.Close()
	if err != nil {
		panic("Could not open texture image file " + fname)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic("Could not decode the image from " + fname)
	}

	rval_len := 4 * img.Bounds().Size().X * img.Bounds().Size().Y
	rval := make([]byte, rval_len)

	offset := func(x, y int) int {
		return 4 * (y*img.Bounds().Size().X + x)
	}

	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			// RGBA returns a uint32, but color values are in a 16 bit range (up to 65535)
			// Bitshifting by 8 just truncates the extra color information that we don't care about
			rval[offset(x, y)], rval[offset(x, y)+1], rval[offset(x, y)+2], rval[offset(x, y)+3] = byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8)
		}
	}

	return rval, uint32(img.Bounds().Size().X), uint32(img.Bounds().Size().Y), vk.FORMAT_R8G8B8A8_SRGB
}

func (app *App_05) createTextureImage() {
	imgData, width, height, format := loadImage()
	size := vk.DeviceSize(len(imgData))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)

	if ptr, err := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0); err != nil {
		panic("Could not map memory for texture staging buffer: " + err.Error())
	} else {
		vk.MemCopySlice(unsafe.Pointer(ptr), imgData)
		vk.UnmapMemory(app.device, stagingBufferMemory)
	}

	app.textureImage, app.textureImageMemory = app.createImage(width, height, format, vk.IMAGE_TILING_OPTIMAL, vk.IMAGE_USAGE_TRANSFER_DST_BIT|vk.IMAGE_USAGE_SAMPLED_BIT, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.transitionImageLayout(app.textureImage, format, vk.IMAGE_LAYOUT_UNDEFINED, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL)
	app.copyBufferToImage(stagingBuffer, app.textureImage, width, height)
	app.transitionImageLayout(app.textureImage, format, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)
}

func (app *App_05) destroyTextureImage() {
	vk.DestroyImageView(app.device, app.textureImageView, nil)
	vk.DestroyImage(app.device, app.textureImage, nil)
	vk.FreeMemory(app.device, app.textureImageMemory, nil)
}

func (app *App_05) createImage(width, height uint32, format vk.Format, tiling vk.ImageTiling, usage vk.ImageUsageFlags, memProps vk.MemoryPropertyFlags) (image vk.Image, imageMemory vk.DeviceMemory) {

	imageCI := vk.ImageCreateInfo{
		ImageType: vk.IMAGE_TYPE_2D,
		Format:    format,
		Extent: vk.Extent3D{
			Width:  width,
			Height: height,
			Depth:  1,
		},
		MipLevels:           1,
		ArrayLayers:         1,
		Tiling:              vk.IMAGE_TILING_OPTIMAL,
		Usage:               usage,
		SharingMode:         vk.SHARING_MODE_EXCLUSIVE,
		PQueueFamilyIndices: []uint32{},
		InitialLayout:       vk.IMAGE_LAYOUT_UNDEFINED,
		Samples:             vk.SAMPLE_COUNT_1_BIT,
	}

	var err error

	if image, err = vk.CreateImage(app.device, &imageCI, nil); err != nil {
		panic("Could not create image: " + err.Error())
	}

	memReq := vk.GetImageMemoryRequirements(app.device, image)
	memAlloc := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: app.findMemoryType(memReq.MemoryTypeBits, memProps),
	}

	if imageMemory, err = vk.AllocateMemory(app.device, &memAlloc, nil); err != nil {
		panic("Could not allocate memory for texture image: " + err.Error())
	}

	if err := vk.BindImageMemory(app.device, image, imageMemory, 0); err != nil {
		panic("Could not bind texture image memory: " + err.Error())
	}

	return
}

func (app *App_05) copyBufferToImage(buffer vk.Buffer, image vk.Image, width, height uint32) {
	copyRegion := vk.BufferImageCopy{
		BufferOffset:      0,
		BufferRowLength:   0,
		BufferImageHeight: 0,
		ImageSubresource: vk.ImageSubresourceLayers{
			AspectMask:     vk.IMAGE_ASPECT_COLOR_BIT,
			MipLevel:       0,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
		ImageOffset: vk.Offset3D{0, 0, 0},
		ImageExtent: vk.Extent3D{width, height, 1},
	}

	cb := app.beginSingleTimeCommands()
	vk.CmdCopyBufferToImage(cb, buffer, image, vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL, []vk.BufferImageCopy{copyRegion})
	app.endSingleTimeCommands(cb)
}

func (app *App_05) transitionImageLayout(image vk.Image, format vk.Format, oldLayout, newLayout vk.ImageLayout) {
	cb := app.beginSingleTimeCommands()

	imgMemBarrier := vk.ImageMemoryBarrier{
		SrcAccessMask:       0,
		DstAccessMask:       0,
		OldLayout:           oldLayout,
		NewLayout:           newLayout,
		SrcQueueFamilyIndex: vk.QUEUE_FAMILY_IGNORED,
		DstQueueFamilyIndex: vk.QUEUE_FAMILY_IGNORED,
		Image:               image,
		SubresourceRange: vk.ImageSubresourceRange{
			AspectMask:     vk.IMAGE_ASPECT_COLOR_BIT,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	var srcStage, dstStage vk.PipelineStageFlags

	if newLayout == vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL {
		imgMemBarrier.SubresourceRange.AspectMask = vk.IMAGE_ASPECT_DEPTH_BIT
	}
	if formatHasStencilComponent(format) {
		imgMemBarrier.SubresourceRange.AspectMask = imgMemBarrier.SubresourceRange.AspectMask | vk.IMAGE_ASPECT_STENCIL_BIT
	}

	if oldLayout == vk.IMAGE_LAYOUT_UNDEFINED && newLayout == vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL {
		imgMemBarrier.SrcAccessMask = 0
		imgMemBarrier.DstAccessMask = vk.ACCESS_TRANSFER_WRITE_BIT
		srcStage = vk.PIPELINE_STAGE_TOP_OF_PIPE_BIT
		dstStage = vk.PIPELINE_STAGE_TRANSFER_BIT

	} else if oldLayout == vk.IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL && newLayout == vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL {
		imgMemBarrier.SrcAccessMask = vk.ACCESS_TRANSFER_WRITE_BIT
		imgMemBarrier.DstAccessMask = vk.ACCESS_SHADER_READ_BIT
		srcStage = vk.PIPELINE_STAGE_TRANSFER_BIT
		dstStage = vk.PIPELINE_STAGE_FRAGMENT_SHADER_BIT

	} else if oldLayout == vk.IMAGE_LAYOUT_UNDEFINED && newLayout == vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL {
		imgMemBarrier.SrcAccessMask = 0
		imgMemBarrier.DstAccessMask = vk.ACCESS_DEPTH_STENCIL_ATTACHMENT_READ_BIT | vk.ACCESS_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT
		srcStage = vk.PIPELINE_STAGE_TOP_OF_PIPE_BIT
		dstStage = vk.PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT

	} else {
		panic("Unhandled image layout transition from " + oldLayout.String() + " to " + newLayout.String())
	}

	vk.CmdPipelineBarrier(cb, srcStage, dstStage, 0, nil, nil, []vk.ImageMemoryBarrier{imgMemBarrier})

	app.endSingleTimeCommands(cb)
}

func (app *App_05) createTextureImageView() {
	app.textureImageView = app.createImageView(app.textureImage, vk.FORMAT_R8G8B8A8_SRGB, vk.IMAGE_ASPECT_COLOR_BIT)
}

func (app *App_05) createTextureSampler() {

	devProps := vk.GetPhysicalDeviceProperties(app.physicalDevice)

	samplerInfo := vk.SamplerCreateInfo{
		MagFilter:               vk.FILTER_LINEAR,
		MinFilter:               vk.FILTER_LINEAR,
		AddressModeU:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AddressModeV:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AddressModeW:            vk.SAMPLER_ADDRESS_MODE_REPEAT,
		AnisotropyEnable:        true,
		MaxAnisotropy:           devProps.Limits.MaxSamplerAnisotropy,
		CompareEnable:           false,
		CompareOp:               vk.COMPARE_OP_ALWAYS,
		BorderColor:             vk.BORDER_COLOR_INT_OPAQUE_BLACK,
		UnnormalizedCoordinates: false,
	}

	var err error
	if app.textureSampler, err = vk.CreateSampler(app.device, &samplerInfo, nil); err != nil {
		panic("Could not create texture image sampler: " + err.Error())
	}
}
