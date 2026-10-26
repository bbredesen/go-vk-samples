package vkctx

import (
	"github.com/bbredesen/go-vk"
	"golang.org/x/sys/windows"
)

type Context struct {
	EnableApiLayers, EnableInstanceExtensions, EnableDeviceExtensions []string

	Instance vk.Instance
	Surface  vk.SurfaceKHR

	PhysicalDevice vk.PhysicalDevice
	Device         vk.Device

	GraphicsQueueFamilyIndex, PresentQueueFamilyIndex uint32
	GraphicsQueue, PresentQueue                       vk.Queue

	// Command pool and buffers
	CommandPool    vk.CommandPool
	CommandBuffers []vk.CommandBuffer // Primary buffers

	// Swapchain handles
	Swapchain             vk.SwapchainKHR
	SwapchainExtent       vk.Extent2D
	SwapchainImageFormat  vk.Format
	SwapchainImages       []vk.Image
	SwapchainImageViews   []vk.ImageView
	SwapChainFramebuffers []vk.Framebuffer

	// Sync objects
	ImageAvailableSemaphore, RenderFinishedSemaphore vk.Semaphore
	InFlightFence                                    vk.Fence
}

func (ctx *Context) Initialize(hInstance windows.Handle, hWnd windows.HWND) {

	ctx.createInstance()
	ctx.createSurface(hInstance, hWnd)

	ctx.selectPhysicalDevice()
	ctx.createLogicalDevice()

	ctx.createSwapchain()
	ctx.createSwapchainImageViews()

	ctx.createCommandPool()
	ctx.createSyncObjects()

}

func (ctx *Context) Teardown() {
	vk.QueueWaitIdle(ctx.PresentQueue)

	ctx.destroySyncObjects()
	ctx.destroyCommandPool()

	ctx.cleanupSwapchain()

	vk.DestroyDevice(ctx.Device, nil)
	vk.DestroySurfaceKHR(ctx.Instance, ctx.Surface, nil)
	vk.DestroyInstance(ctx.Instance, nil)
}

func (ctx *Context) CreateImage(extent vk.Extent2D, format vk.Format, tiling vk.ImageTiling, usage vk.ImageUsageFlags, memProps vk.MemoryPropertyFlags) (image vk.Image, imageMemory vk.DeviceMemory) {

	imageCI := vk.ImageCreateInfo{
		ImageType: vk.IMAGE_TYPE_2D,
		Format:    format,
		Extent: vk.Extent3D{
			Width:  extent.Width,
			Height: extent.Height,
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

	var r vk.Result

	if r, image = vk.CreateImage(ctx.Device, &imageCI, nil); r != vk.SUCCESS {
		panic("Could not create image: " + r.String())
	}

	memReq := vk.GetImageMemoryRequirements(ctx.Device, image)
	memAlloc := vk.MemoryAllocateInfo{
		AllocationSize:  memReq.Size,
		MemoryTypeIndex: ctx.FindMemoryType(memReq.MemoryTypeBits, memProps),
	}

	if r, imageMemory = vk.AllocateMemory(ctx.Device, &memAlloc, nil); r != vk.SUCCESS {
		panic("Could not allocate memory for texture image: " + r.String())
	}

	if r := vk.BindImageMemory(ctx.Device, image, imageMemory, 0); r != vk.SUCCESS {
		panic("Could not bind texture image memory: " + r.String())
	}

	return
}

func (ctx *Context) CreateImageView(image vk.Image, format vk.Format, aspectMask vk.ImageAspectFlags) vk.ImageView {
	ivCI := vk.ImageViewCreateInfo{
		Image:    image,
		ViewType: vk.IMAGE_VIEW_TYPE_2D,
		Format:   format,
		SubresourceRange: vk.ImageSubresourceRange{
			AspectMask:     aspectMask,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}

	if r, iv := vk.CreateImageView(ctx.Device, &ivCI, nil); r != vk.SUCCESS {
		panic("Could not create image view: " + r.String())
	} else {
		return iv
	}

}

func (ctx *Context) FindMemoryType(typeFilter uint32, flags vk.MemoryPropertyFlags) uint32 {
	memProps := vk.GetPhysicalDeviceMemoryProperties(ctx.PhysicalDevice)
	var i uint32
	for i = 0; i < memProps.MemoryTypeCount; i++ {
		if typeFilter&1<<i != 0 && memProps.MemoryTypes[i].PropertyFlags&flags == flags {
			return i
		}
	}
	panic("Could not find appropriate memory type.")
}
