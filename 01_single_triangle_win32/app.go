package main

import (
	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/go-vk-samples/shared"
)

type App_01 struct {
	shared.App

	eventsChannel <-chan shared.EventMessage
	windowHandle  uintptr

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
}

func NewApp() *App_01 {
	app := &App_01{
		enableInstanceExtensions: []string{
			"VK_KHR_surface", "VK_EXT_metal_surface",
			"VK_EXT_debug_report",
			"VK_EXT_debug_utils",
		},
		enableDeviceExtensions: []string{"VK_KHR_swapchain"},
		enableApiLayers:        []string{},
	}

	app.App, _ = shared.NewApp()

	return app
}

func (app *App_01) MainLoop(ch <-chan shared.EventMessage) {

	// Read any system messages...input, resize, window close, etc.
	for m, open := <-ch; open; m, open = <-ch {
		switch m.Type {
		case shared.ET_Sys_Create:
			app.windowHandle = m.HandleForSurface
			app.InitVulkan()
		}
		// Rendering goes here
		app.drawFrame()

	}
}

func (app *App_01) Run(windowTitle string) {

	go app.MainLoop(app.App.GetEventChannel())

	app.App.Run()

}

func (app *App_01) InitVulkan() {
	app.createInstance()
	app.createSurface()

	app.selectPhysicalDevice()
	app.createLogicalDevice()

	app.createSwapchain()
	app.createImageViews()

	app.createRenderPass()
	app.createGraphicsPipeline()
	app.createFramebuffers()

	app.createCommandPool()
	app.createSyncObjects()

}

func (app *App_01) CleanupVulkan() {
	vk.QueueWaitIdle(app.presentQueue)

	app.destroySyncObjects()
	app.destroyCommandPool()

	// handles pipeline, shader modules, and renderpass
	app.cleanupGraphicsPipeline()

	app.cleanupSwapchain()

	vk.DestroyDevice(app.device, nil)
	vk.DestroySurfaceKHR(app.instance, app.surface, nil)
	vk.DestroyInstance(app.instance, nil)
}

func (app *App_01) drawFrame() {
	vk.WaitForFences(app.device, []vk.Fence{app.inFlightFence}, true, ^uint64(0))

	r, imageIndex := vk.AcquireNextImageKHR(app.device, app.swapchain, ^uint64(0), app.imageAvailableSemaphore, vk.Fence(vk.NULL_HANDLE))
	if r != vk.SUCCESS {
		if r == vk.SUBOPTIMAL_KHR || r == vk.ERROR_OUT_OF_DATE_KHR {
			app.recreateSwapchain()
			return
		} else {
			panic("Could not acquire next image!")
		}
	}

	vk.ResetFences(app.device, []vk.Fence{app.inFlightFence})

	vk.ResetCommandBuffer(app.commandBuffers[imageIndex], 0)
	app.recordCommandBuffer(app.commandBuffers[imageIndex], imageIndex)

	submitInfo := vk.SubmitInfo{
		PWaitSemaphores:   []vk.Semaphore{app.imageAvailableSemaphore},
		PWaitDstStageMask: []vk.PipelineStageFlags{vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT},
		PCommandBuffers:   []vk.CommandBuffer{app.commandBuffers[imageIndex]},
		PSignalSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
	}

	if r := vk.QueueSubmit(app.graphicsQueue, []vk.SubmitInfo{submitInfo}, app.inFlightFence); r != vk.SUCCESS {
		panic("Could not submit to graphics queue! ")
	}

	// Present the drawn image
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
		PSwapchains:     []vk.SwapchainKHR{app.swapchain},
		PImageIndices:   []uint32{imageIndex},
	}

	if r := vk.QueuePresentKHR(app.presentQueue, &presentInfo); r != vk.SUCCESS && r != vk.SUBOPTIMAL_KHR && r != vk.ERROR_OUT_OF_DATE_KHR {
		panic("Could not submit to presentation queue! ")
	}

}
