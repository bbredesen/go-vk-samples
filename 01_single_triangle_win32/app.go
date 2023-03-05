package main

import (
	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/go-vk-samples/shared"
)

type App_01 struct {
	*shared.Win32App
	messages                                                          <-chan shared.WindowMessage
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

func NewApp() App_01 {
	c := make(chan shared.WindowMessage, 32)
	winapp := shared.NewWin32App(c)

	return App_01{
		Win32App:                 winapp,
		messages:                 c,
		enableInstanceExtensions: winapp.GetRequiredInstanceExtensions(),
		enableDeviceExtensions:   []string{"VK_KHR_swapchain"},
		enableApiLayers:          []string{},
	}
}

func (app *App_01) MainLoop() {

	// Read any system messages...input, resize, window close, etc.
	for {
	innerLoop:
		for {
			select {
			case msg := <-app.messages:
				// fmt.Println(msg.Text)
				switch msg.Text {
				case "DESTROY":
					// Break out of the loop
					return

				}
			default:
				// Pull everything off the queue, then continue the outer loop
				break innerLoop // "break" will break out of the select statement, not the loop, so we have to use a break label
			}

		}

		// Rendering goes here
		app.drawFrame()

	}
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
			panic("Could not acquire next image! " + r.String())
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
		panic("Could not submit to graphics queue! " + r.String())
	}

	// Present the drawn image
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{app.renderFinishedSemaphore},
		PSwapchains:     []vk.SwapchainKHR{app.swapchain},
		PImageIndices:   []uint32{imageIndex},
	}

	if r := vk.QueuePresentKHR(app.presentQueue, &presentInfo); r != vk.SUCCESS && r != vk.SUBOPTIMAL_KHR && r != vk.ERROR_OUT_OF_DATE_KHR {
		panic("Could not submit to presentation queue! " + r.String())
	}

}
