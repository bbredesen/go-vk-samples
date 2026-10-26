package main

import (
	"github.com/bbredesen/go-vk"
)

// Must match the stencil attachment format in createRenderPass()
const stencilFormat = vk.FORMAT_S8_UINT

// Create the stencil image, memory, and view used as attachment 1 in the framebuffers. Call before
// createFramebuffers(), and again after the swapchain is recreated, since the image must match the swapchain extent.
func (app *App_10) createStencilResources() {
	props := vk.GetPhysicalDeviceFormatProperties(app.physicalDevice, stencilFormat)
	if props.OptimalTilingFeatures&vk.FORMAT_FEATURE_DEPTH_STENCIL_ATTACHMENT_BIT == 0 {
		panic("Stencil format " + stencilFormat.String() + " is not supported as an attachment on this device!")
	}

	app.stencilImage, app.stencilImageMemory = app.createImage(
		app.swapchainExtent,
		stencilFormat,
		vk.IMAGE_TILING_OPTIMAL,
		vk.IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT,
		vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT,
	)

	app.stencilImageView = app.createImageView(app.stencilImage, stencilFormat, vk.IMAGE_ASPECT_STENCIL_BIT)
}

func (app *App_10) cleanupStencilResources() {
	vk.DestroyImageView(app.device, app.stencilImageView, nil)
	vk.DestroyImage(app.device, app.stencilImage, nil)
	vk.FreeMemory(app.device, app.stencilImageMemory, nil)

	app.stencilImageView = vk.ImageView(vk.NULL_HANDLE)
	app.stencilImage = vk.Image(vk.NULL_HANDLE)
	app.stencilImageMemory = vk.DeviceMemory(vk.NULL_HANDLE)
}

func (app *App_10) createImage(extent vk.Extent2D, format vk.Format, tiling vk.ImageTiling, usage vk.ImageUsageFlags, fl vk.MemoryPropertyFlags) (image vk.Image, memory vk.DeviceMemory) {
	imageCI := vk.ImageCreateInfo{
		ImageType: vk.IMAGE_TYPE_2D,
		Format:    format,
		Extent: vk.Extent3D{
			Width:  extent.Width,
			Height: extent.Height,
			Depth:  1,
		},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       vk.SAMPLE_COUNT_1_BIT,
		Tiling:        tiling,
		Usage:         usage,
		SharingMode:   vk.SHARING_MODE_EXCLUSIVE,
		InitialLayout: vk.IMAGE_LAYOUT_UNDEFINED,
	}

	var err error
	if image, err = vk.CreateImage(app.device, &imageCI, nil); err != nil {
		panic("Could not create image: " + err.Error())
	}

	memReq := vk.GetImageMemoryRequirements(app.device, image)

	memAllocInfo := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: app.findMemoryType(memReq.MemoryTypeBits, fl),
	}

	if memory, err = vk.AllocateMemory(app.device, &memAllocInfo, nil); err != nil {
		panic("Could not allocate memory for image: " + err.Error())
	}
	if err := vk.BindImageMemory(app.device, image, memory, 0); err != nil {
		panic("Could not bind memory for image: " + err.Error())
	}

	return
}
