package main

import (
	"math"
	"time"
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/go-vk-samples/shared"
	"github.com/bbredesen/vkm"
)

type App_03 struct {
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
}

func NewApp(windowTitle string) App_03 {
	sharedApp, _ := shared.NewApp(windowTitle)

	return App_03{
		App:           sharedApp,
		eventsChannel: sharedApp.GetEventChannel(),

		enableInstanceExtensions: sharedApp.GetRequiredInstanceExtensions(),
		enableDeviceExtensions:   []string{"VK_KHR_swapchain"},
		enableApiLayers:          []string{},

		startTime: time.Now(),
	}
}

func (app *App_03) MainLoop(ch <-chan shared.EventMessage) {

	m := <-ch // Block on the channel until the window has been created
	if m.Type != shared.ET_Sys_Created {
		panic("expected ET_Sys_Create to start main loop")
	}
	app.InitVulkan()

	for {
	messageLoop:
		for {
			select {
			case m = <-ch:
				switch m.Type {
				case shared.ET_Sys_Closed:
					app.CleanupVulkan()
					app.OkToClose(m.SystemEvent.HandleForSurface)
					return

				}
			default: // Channel is empty
				break messageLoop

			}
		}

		app.drawFrame()
	}
}

func (app *App_03) Run() {
	go app.MainLoop(app.App.GetEventChannel())
	app.App.Run()
}

func (app *App_03) InitVulkan() {

	app.createInstance()
	app.createSurface()

	app.selectPhysicalDevice()
	app.createLogicalDevice()

	app.createSwapchain()
	app.createImageViews()

	app.createDescriptorSetLayout()

	app.createUniformBuffers()
	app.createDescriptorPool()
	app.createDescriptorSets()

	app.createRenderPass()
	app.createGraphicsPipeline()
	app.createFramebuffers()
	app.createCommandPool()
	app.createSyncObjects()

	app.createVertexBuffer()
	app.createIndexBuffer()

}

func (app *App_03) CleanupVulkan() {
	vk.QueueWaitIdle(app.presentQueue)

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

func (app *App_03) drawFrame() {
	vk.WaitForFences(app.device, []vk.Fence{app.inFlightFence}, true, ^uint64(0))

	var err error
	if app.currentImage, err = vk.AcquireNextImageKHR(app.device, app.swapchain, ^uint64(0), app.imageAvailableSemaphore, vk.Fence(vk.NULL_HANDLE)); err != nil {
		if err == vk.SUBOPTIMAL_KHR || err == vk.ERROR_OUT_OF_DATE_KHR {
			app.recreateSwapchain()
			return
		} else {
			panic("Could not acquire next image! " + err.Error())
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

	if err := vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, app.inFlightFence); err != nil {
		panic("Could not submit to graphics queue! " + err.Error())
	}

	// Present the drawn image
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
		PSwapchains:     []vk.SwapchainKHR{app.swapchain},
		PImageIndices:   []uint32{app.currentImage},
	}

	if err := vk.QueuePresentKHR(app.presentQueue, &presentInfo); err != nil && err != vk.SUBOPTIMAL_KHR && err != vk.ERROR_OUT_OF_DATE_KHR {
		panic("Could not submit to presentation queue! " + err.Error())
	}

}

func (app *App_03) updateUniformBuffer(currentImage uint32) {
	seconds := time.Since(app.startTime).Seconds()
	_ = seconds
	app.uboObjs[currentImage].model = vkm.NewMatRotate(vkm.UnitVecZ(), float32(seconds*math.Pi/2.0))

	// View and projection are initialized in createUniformBuffers(), there is no need to recalculate them every frame;
	// see ubo.go

	vk.MemCopyObj(unsafe.Pointer(app.uniformBufferMapped[currentImage]), &app.uboObjs[currentImage])
}
