package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

func (app *App_01) createSwapchain() {
	/*
		This is a big function:
		1) Determine general surface capabilities, like extents and swapchain limits
		2) Determine the image formats supported by the surface
		3) Determine the present modes supported by the surface...does the GPU update at
		   screen refresh, or immediately, etc.
		4) Make a decision about number of images in the swapchain. This will almost
		   always be two, but check the hardware to be sure.
		5) Select a surface format to use (color space, etc.)
		6) Select the extent to use for the swapchain images, which will be the window
		   size in this case.
		7) Select the surface transform to use in presentation (rotation, mirroring,
		   etc.) Prefer no transform ("identity")
		8) Select the present mode to use, preferrring mailbox, or falling back to FIFO.
		9) Take all those decisions and create the swapchain in Vulkan
		10) Get references to the images in the swapchain, and store them.
	*/

	// 1) General capabilities
	var err error
	surfaceCapabilities, err := vk.GetPhysicalDeviceSurfaceCapabilitiesKHR(app.physicalDevice, app.surface)
	if err != nil {
		panic("Could not get device surface capabilities: ")
	}

	// 2) Image formats
	surfaceFormats, err := vk.GetPhysicalDeviceSurfaceFormatsKHR(
		app.physicalDevice, app.surface,
	)
	if err != nil {
		panic("Could not get device surface formats: ")
	}

	// 3) Supported present modes
	presentModes, err := vk.GetPhysicalDeviceSurfacePresentModesKHR(app.physicalDevice, app.surface)
	if err != nil {
		panic("Could not get device surface present modes: ")
	}

	// 4) Decide on swapchain size
	imageCount := uint32(surfaceCapabilities.MinImageCount) + 1
	// Ensure we aren't over the maximum image count
	// MaxImageCount == 0 means there is no limit.
	if surfaceCapabilities.MaxImageCount != 0 && surfaceCapabilities.MaxImageCount < imageCount {
		imageCount = surfaceCapabilities.MaxImageCount
	}

	// 5) Select surface format to use
	// Prefer FORMAT_B8G8R8A8_SRGB, or fallback to the first option
	fmtIdx := -1
	var selectedSurfaceFormat vk.SurfaceFormatKHR
	for i, fmt := range surfaceFormats {
		if fmt.Format == vk.FORMAT_B8G8R8A8_SRGB && fmt.ColorSpace == vk.COLOR_SPACE_SRGB_NONLINEAR_KHR {
			fmtIdx = i
			break
		}
	}
	if fmtIdx < 0 {
		fmtIdx = 0
	}
	selectedSurfaceFormat = surfaceFormats[fmtIdx]

	// 6) Select swap chain extent
	// If current extent has values equal to max uint32, then we need to set the
	// extent somewhere between the min and max extent values. 640x480 is chosen
	// as an arbitrary target. Otherwise, just return the current surface
	// extent. The surface extent SHOULD equal the area in the window that is already open.
	var selectedExtent vk.Extent2D
	if surfaceCapabilities.CurrentExtent.Width == ^uint32(0) {
		min := func(a, b uint32) uint32 {
			if a < b {
				return a
			} else {
				return b
			}
		}
		max := func(a, b uint32) uint32 {
			if a > b {
				return a
			} else {
				return b
			}
		}

		selectedExtent.Width =
			min(
				max(640, surfaceCapabilities.MinImageExtent.Width),
				surfaceCapabilities.MaxImageExtent.Width,
			)
		selectedExtent.Height =
			min(
				max(480, surfaceCapabilities.MinImageExtent.Height),
				surfaceCapabilities.MaxImageExtent.Height,
			)
	} else {
		selectedExtent = surfaceCapabilities.CurrentExtent
	}

	// 7) Select the surface transformation, from the final GPU image to the
	// surface display
	var selectedTransform vk.SurfaceTransformFlagsKHR = surfaceCapabilities.CurrentTransform
	if (surfaceCapabilities.SupportedTransforms & vk.SURFACE_TRANSFORM_IDENTITY_BIT_KHR) != 0x0 {
		selectedTransform = vk.SURFACE_TRANSFORM_IDENTITY_BIT_KHR
	}

	// 8) Choose present mode: How the swapchain images are pushed to the
	// surface
	var presentMode vk.PresentModeKHR = vk.PRESENT_MODE_FIFO_KHR
	for _, mode := range presentModes {
		if mode == vk.PRESENT_MODE_MAILBOX_KHR {
			presentMode = vk.PRESENT_MODE_MAILBOX_KHR
			break
		}
	}

	// 9) Combine all that info into the swapchain create info struct
	createInfo := vk.SwapchainCreateInfoKHR{
		Flags:            vk.SwapchainCreateFlagsKHR(0),
		Surface:          app.surface,
		MinImageCount:    imageCount,
		ImageFormat:      selectedSurfaceFormat.Format,
		ImageColorSpace:  selectedSurfaceFormat.ColorSpace,
		ImageExtent:      selectedExtent,
		ImageArrayLayers: 1, // 2+ for stereoscopic rendering, i.e. VR
		ImageUsage:       vk.IMAGE_USAGE_COLOR_ATTACHMENT_BIT,
		ImageSharingMode: vk.SHARING_MODE_EXCLUSIVE,
		PreTransform:     selectedTransform,
		CompositeAlpha:   vk.COMPOSITE_ALPHA_OPAQUE_BIT_KHR,
		PresentMode:      presentMode,
		Clipped:          true,
	}

	swapchain, err := vk.CreateSwapchainKHR(app.device, &createInfo, nil)
	if err != nil {
		fmt.Printf("Failed to create swapchain!")
		panic(err)
	}
	app.swapchain = swapchain
	app.swapchainImageFormat = selectedSurfaceFormat.Format
	app.swapchainExtent = selectedExtent

	// 10) Finally, get the swapchain images and save them
	images, err := vk.GetSwapchainImagesKHR(app.device, app.swapchain)
	if err != nil {
		fmt.Printf("Could not get images after createing swapchain!\n")
		panic(err)
	}
	app.swapchainImages = images

}

func (app *App_01) createImageViews() {
	// Careful...if image views already exist, this will cause a leak. Call destroyImageViews() first!
	app.swapchainImageViews = make([]vk.ImageView, len(app.swapchainImages))

	for i, img := range app.swapchainImages {
		ivCI := vk.ImageViewCreateInfo{
			Image:    img,
			ViewType: vk.IMAGE_VIEW_TYPE_2D,
			Format:   app.swapchainImageFormat,
			Components: vk.ComponentMapping{
				R: vk.COMPONENT_SWIZZLE_IDENTITY,
				G: vk.COMPONENT_SWIZZLE_IDENTITY,
				B: vk.COMPONENT_SWIZZLE_IDENTITY,
				A: vk.COMPONENT_SWIZZLE_IDENTITY,
			},
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask:     vk.IMAGE_ASPECT_COLOR_BIT,
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}

		if iv, err := vk.CreateImageView(app.device, &ivCI, nil); err != nil {
			panic(err)
		} else {
			app.swapchainImageViews[i] = iv
		}
	}
}

func (app *App_01) destroyImageViews() {
	for _, iv := range app.swapchainImageViews {
		vk.DestroyImageView(app.device, iv, nil)
	}
	app.swapchainImageViews = nil
}

func (app *App_01) cleanupSwapchain() {
	for _, fb := range app.swapChainFramebuffers {
		vk.DestroyFramebuffer(app.device, fb, nil)
	}
	for _, iv := range app.swapchainImageViews {
		vk.DestroyImageView(app.device, iv, nil)
	}

	vk.DestroySwapchainKHR(app.device, app.swapchain, nil)
}

func (app *App_01) recreateSwapchain() {
	if err := vk.DeviceWaitIdle(app.device); err != nil {
		panic(err)
	}

	app.cleanupSwapchain()

	app.createSwapchain()
	app.createImageViews()
	app.createFramebuffers()
}
