package main

import (
	"math"
	"time"
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/go-vk-samples/shared"
	"github.com/bbredesen/vkm"
)

type App_04 struct {
	shared.App
	eventsChannel <-chan shared.EventMessage

	enableInstanceExtensions, enableDeviceExtensions, enableApiLayers []string

	instance vk.Instance
	surface  vk.SurfaceKHR

	physicalDevice vk.PhysicalDevice
	device         vk.Device

	// Device queue information
	graphicsQueueFamilyIndex, presentQueueFamilyIndex uint32
	graphicsQueue, presentQueue                       vk.Queue

	// Command pool and buffers
	commandPool    vk.CommandPool
	commandBuffers []vk.CommandBuffer // Primary buffers

	// Swapchain handles
	swapchain             vk.SwapchainKHR
	swapchainExtent       vk.Extent2D
	swapchainImageFormat  vk.Format
	swapchainImages       []vk.Image
	swapchainImageViews   []vk.ImageView
	swapChainFramebuffers []vk.Framebuffer

	// Pipeline
	pipelineLayout   vk.PipelineLayout
	graphicsPipeline vk.Pipeline

	// Renderpass
	renderPass vk.RenderPass

	depthImage       vk.Image
	depthImageMemory vk.DeviceMemory
	depthImageView   vk.ImageView

	// Sync objects
	imageAvailableSemaphore, renderFinishedSemaphore vk.Semaphore
	inFlightFence                                    vk.Fence

	vertShaderModule, fragShaderModule vk.ShaderModule

	vertexBuffer, indexBuffer             vk.Buffer
	vertexBufferMemory, indexBufferMemory vk.DeviceMemory

	// Descriptors
	descriptorSetLayout vk.DescriptorSetLayout
	descriptorPool      vk.DescriptorPool
	descriptorSets      []vk.DescriptorSet

	uniformBuffers        []vk.Buffer
	uniformBufferMemories []vk.DeviceMemory
	uniformBufferMapped   []*byte
	uboObjs               []UniformBufferObject

	startTime    time.Time
	currentImage uint32

	// Textures
	textureImage       vk.Image
	textureImageView   vk.ImageView
	textureImageMemory vk.DeviceMemory
	textureSampler     vk.Sampler
}

func NewApp() App_04 {
	sharedApp, _ := shared.NewApp()

	return App_04{
		App:                      sharedApp,
		eventsChannel:            sharedApp.GetEventChannel(),
		enableInstanceExtensions: sharedApp.GetRequiredInstanceExtensions(),
		enableDeviceExtensions:   []string{"VK_KHR_swapchain"},
		enableApiLayers:          []string{},

		startTime: time.Now(),
	}
}

func (app *App_04) MainLoop(ch <-chan shared.EventMessage) {

	m := <-ch
	if m.Type != shared.ET_Sys_Create {
		panic("expected ET_Sys_Create to start mainloop")
	}
	app.InitVulkan()

	for {
		// Read any system messages...input, resize, window close, etc.
		// for m, open := <-ch; open; m, open = <-ch {
	innerLoop:
		for {
			select {
			case m := <-ch:
				switch m.Type {
				// case shared.ET_Sys_Create:
				// 	app.InitVulkan()
				}
			default:
				break innerLoop

			}
		}
		// Rendering goes here
		app.drawFrame()
	}
}

func (app *App_04) Run(windowTitle string) {

	go app.MainLoop(app.App.GetEventChannel())

	app.App.Run()
}

func (app *App_04) InitVulkan() {

	app.createInstance()
	app.createSurface()

	app.selectPhysicalDevice()
	app.createLogicalDevice()

	app.createSwapchain()
	app.createSwapchainImageViews()

	app.createCommandPool()

	app.createDescriptorSetLayout()

	app.createUniformBuffers()

	app.createTextureImage()
	app.createTextureImageView()
	app.createTextureSampler()

	app.createDescriptorPool()
	app.createDescriptorSets()

	app.createRenderPass()
	app.createGraphicsPipeline()
	app.createFramebuffers()

	app.createSyncObjects()

	app.createVertexBuffer()
	app.createIndexBuffer()

}

func (app *App_04) CleanupVulkan() {
	vk.QueueWaitIdle(app.presentQueue)

	app.destroyTextureImage()
	vk.DestroySampler(app.device, app.textureSampler, nil)

	app.destroySyncObjects()
	app.destroyCommandPool()

	// handles pipeline, shader modules, and renderpass
	app.cleanupGraphicsPipeline()

	app.cleanupSwapchain()

	// Also destroys the allocated descriptor sets
	vk.DestroyDescriptorPool(app.device, app.descriptorPool, nil)
	vk.DestroyDescriptorSetLayout(app.device, app.descriptorSetLayout, nil)

	app.cleanupUniformBuffers()

	vk.DestroyBuffer(app.device, app.vertexBuffer, nil)
	vk.FreeMemory(app.device, app.vertexBufferMemory, nil)
	vk.DestroyBuffer(app.device, app.indexBuffer, nil)
	vk.FreeMemory(app.device, app.indexBufferMemory, nil)

	vk.DestroyDevice(app.device, nil)
	vk.DestroySurfaceKHR(app.instance, app.surface, nil)
	vk.DestroyInstance(app.instance, nil)
}

func (app *App_04) drawFrame() {
	vk.WaitForFences(app.device, []vk.Fence{app.inFlightFence}, true, ^uint64(0))

	var r vk.Result
	if r, app.currentImage = vk.AcquireNextImageKHR(app.device, app.swapchain, ^uint64(0), app.imageAvailableSemaphore, vk.Fence(vk.NULL_HANDLE)); r != vk.SUCCESS {
		if r == vk.SUBOPTIMAL_KHR || r == vk.ERROR_OUT_OF_DATE_KHR {
			app.recreateSwapchain()
			return
		} else {
			panic("Could not acquire next image! " + r.String())
		}
	}

	vk.ResetFences(app.device, []vk.Fence{app.inFlightFence})

	vk.ResetCommandBuffer(app.commandBuffers[app.currentImage], 0)
	app.recordCommandBuffer(app.commandBuffers[app.currentImage], app.currentImage)

	app.updateUniformBuffer(app.currentImage)

	submitInfo := vk.SubmitInfo{
		PWaitSemaphores:   []vk.Semaphore{app.imageAvailableSemaphore},
		PWaitDstStageMask: []vk.PipelineStageFlags{vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT},
		PCommandBuffers:   []vk.CommandBuffer{app.commandBuffers[app.currentImage]},
		PSignalSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
	}

	if r := vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, app.inFlightFence); r != vk.SUCCESS {
		panic("Could not submit to graphics queue! " + r.String())
	}

	// Present the drawn image
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
		PSwapchains:     []vk.SwapchainKHR{app.swapchain},
		PImageIndices:   []uint32{app.currentImage},
	}

	if r := vk.QueuePresentKHR(app.presentQueue, &presentInfo); r != vk.SUCCESS && r != vk.SUBOPTIMAL_KHR && r != vk.ERROR_OUT_OF_DATE_KHR {
		panic("Could not submit to presentation queue! " + r.String())
	}

}

func (app *App_04) updateUniformBuffer(currentImage uint32) {
	seconds := time.Since(app.startTime).Seconds()
	_ = seconds
	app.uboObjs[currentImage].model = vkm.NewMatRotate(vkm.UnitVecZ(), float32(seconds*math.Pi/2.0))

	// View and projection are initialized in createUniformBuffers(), there is no need to recalculate them every frame;
	// see ubo.go

	// vk.MemCopyObj(unsafe.Pointer(app.uniformBufferMapped[currentImage]), &app.uboObjs[currentImage])

	// var sl = struct {
	// 	addr uintptr
	// 	len  int
	// 	cap  int
	// }{uintptr(unsafe.Pointer(app.uniformBufferMapped[currentImage])), 192, 192}
	// bytes := *(*[]byte)(unsafe.Pointer(&sl))

	// // bytes := *(*[]byte)(unsafe.Pointer(app.uniformBufferMapped[currentImage]))
	// copy(bytes, AnyTypeToBytes(app.uboObjs[currentImage]))

	vk.MemCopyObj(unsafe.Pointer(app.uniformBufferMapped[currentImage]), &app.uboObjs[currentImage])
}
