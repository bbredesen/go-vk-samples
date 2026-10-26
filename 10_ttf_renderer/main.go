package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bbredesen/go-vk"
	// "github.com/bbredesen/ttf-renderer/shared"
	// "github.com/bbredesen/ttf-renderer/vkctx"
	// "github.com/sirupsen/logrus"
)

//go:generate glslc shaders/shader.vert -o shaders/vert.spv
//go:generate glslc shaders/shader.frag -o shaders/frag.spv
//go:generate glslc shaders/quad_shader.vert -o shaders/quad_vert.spv
//go:generate glslc shaders/quad_shader.frag -o shaders/quad_frag.spv

func init() {
	// /System/Library/Fonts/Supplemental/Georgia.ttf
	// C:\WINDOWS\FONTS\ELEPHNT.TTF
	flag.StringVar(&fontFilename, "font", `/System/Library/Fonts/Supplemental/Georgia.ttf`, "filename to render")
	// flag.StringVar(&fontFilename, "font", `C:\WINDOWS\FONTS\BKANT.TTF`BAHNSCHRIFT, "filename to render")
	flag.StringVar(&renderString, "char", "R", "single character to render")

	flag.Parse()
}

var (
	fontFilename, renderString string
)

const (
	ppem = 640
)

func main() {
	fmt.Printf("go-vk - Glyph Rendering with Stencils\n")
	fmt.Printf("This program demonstrates loading and displaying a font glyph in Go + Vulkan.\n\n")

	if override := os.Getenv("GOVK_SAMPLES_VULKAN_LIB"); override != "" {
		vk.OverrideDefaultVulkanLibrary(override)
	}

	app := NewApp("go-vk Sample: 10_ttf_renderer")
	if props, err := vk.EnumerateInstanceLayerProperties(); err != nil {
		panic("Could not enumerate available layers: " + err.Error())
	} else {
		foundValidationLayers := false
		for _, p := range props {
			if p.LayerName == "VK_LAYER_KHRONOS_validation" {
				app.enableApiLayers = append(app.enableApiLayers, p.LayerName)
				foundValidationLayers = true
				break
			}
		}

		if !foundValidationLayers {
			fmt.Println("NOTE: Khronos validation layer was not found!") // This is normal if you don't have the LunarG SDK or when using MoltenVK
		}
	}

	// Initialize the app and open the window
	app.Run()

	fmt.Println()
	fmt.Println("Clean shutdown, exiting...")

	// app.Initialize()

	// app.loadBuffers(segments, bounds)

	// app.winapp.DefaultMainLoop(shared.DefaultIgnoreInput, shared.DefaultIgnoreTick, app.drawFrame)

	// app.Teardown()
	// Safe exit
}

// type App struct {
// 	winapp   *shared.Win32App
// 	messages chan shared.WindowMessage

// 	vkctx.Context
// 	VulkanPipeline

// 	currentImage uint32

// 	vertexBuffer, indexBuffer             vk.Buffer
// 	vertexBufferMemory, indexBufferMemory vk.DeviceMemory

// 	indexCount int

// 	quadVertStart, quadIndsStart int
// }

// func NewApp() *App {
// 	c := make(chan shared.WindowMessage, 32)

// 	return &App{
// 		winapp:   shared.NewWin32App(c),
// 		messages: c,
// 	}
// }

// func (app *App) Initialize() {
// 	app.winapp.ClassName = "ttf-renderer"
// 	app.winapp.Width, app.winapp.Height = 800, 800
// 	app.winapp.Initialize("ttf-renderer")

// 	app.EnableApiLayers = append(app.EnableApiLayers, "VK_LAYER_KHRONOS_validation")
// 	app.EnableInstanceExtensions = app.winapp.GetRequiredInstanceExtensions()
// 	app.EnableDeviceExtensions = append(app.EnableDeviceExtensions, vk.KHR_SWAPCHAIN_EXTENSION_NAME)

// 	app.Context.Initialize(windows.Handle(app.winapp.HInstance), windows.HWND(app.winapp.HWnd))

// 	app.VulkanPipeline.Initialize(&app.Context)
// }

// func (app *App) Teardown() {
// 	vk.DeviceWaitIdle(app.ctx.Device)

// 	app.destroyBuffers()

// 	app.VulkanPipeline.Teardown()
// 	app.Context.Teardown()
// }

// func (app *App) recordRenderingCommands(cb vk.CommandBuffer) {
// 	cbBeginInfo := vk.CommandBufferBeginInfo{
// 		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
// 	}

// 	colorCV, stencilCV := vk.ClearValue{}, vk.ClearValue{}
// 	ccv := vk.ClearColorValue{}
// 	ccv.AsTypeFloat32([4]float32{0, 0, 0, 1})
// 	colorCV.AsColor(ccv)

// 	stencilCV.AsDepthStencil(vk.ClearDepthStencilValue{
// 		Depth:   0,
// 		Stencil: 0,
// 	})

// 	rpBeginInfo := vk.RenderPassBeginInfo{
// 		RenderPass:  app.renderPass,
// 		Framebuffer: app.SwapChainFramebuffers[app.currentImage],
// 		RenderArea: vk.Rect2D{
// 			Offset: vk.Offset2D{X: 0, Y: 0},
// 			Extent: app.SwapchainExtent,
// 		},
// 		PClearValues: []vk.ClearValue{colorCV, stencilCV},
// 	}

// 	vk.BeginCommandBuffer(cb, &cbBeginInfo)

// 	vk.CmdBeginRenderPass(cb, &rpBeginInfo, vk.SUBPASS_CONTENTS_INLINE)

// 	// bind vert, index bufs
// 	vk.CmdBindVertexBuffers(cb, 0, []vk.Buffer{app.vertexBuffer}, []vk.DeviceSize{0})
// 	vk.CmdBindIndexBuffer(cb, app.indexBuffer, 0, vk.INDEX_TYPE_UINT16)

// 	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[0]) // stencil pipeline
// 	vk.CmdDrawIndexed(cb, uint32(app.quadIndsStart), 1, 0, 0, 0)

// 	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[1]) // stencil quad portion pipeline
// 	vk.CmdDrawIndexed(cb, uint32(app.indexCount-app.quadIndsStart)-4, 1, uint32(app.quadIndsStart), int32(app.quadVertStart), 0)

// 	vk.CmdNextSubpass(cb, vk.SUBPASS_CONTENTS_INLINE)

// 	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[2]) // Color pass

// 	// vk.CmdDrawIndexed(cb, 5, 1, uint32(app.indexCount)-5, 0, 0)
// 	vk.CmdDrawIndexed(cb, 4, 1, uint32(app.indexCount)-4, int32(app.quadVertStart), 0)

// 	// draw

// 	vk.CmdEndRenderPass(cb)
// 	vk.EndCommandBuffer(cb)

// }

func versionToString(version uint32) string {
	return fmt.Sprintf("%d.%d.%d",
		vk.API_VERSION_MAJOR(version),
		vk.API_VERSION_MINOR(version),
		vk.API_VERSION_PATCH(version),
	)
}
