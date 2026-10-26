package main

import (
	"flag"
	"os"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/ttf-renderer/shared"
	"github.com/bbredesen/ttf-renderer/vkctx"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/sys/windows"
)

func init() {
	flag.StringVar(&fontFilename, "font", `C:\WINDOWS\FONTS\ELEPHNT.TTF`, "filename to render")
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
	fontBytes, err := os.ReadFile(fontFilename)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": fontFilename,
			"error":    err,
		}).Error("Failed to open font file")
		os.Exit(1)
	}

	fontData, err := sfnt.Parse(fontBytes)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"filename": fontFilename,
			"error":    err,
		}).Error("Failed to parse font data")
		os.Exit(1)
	}

	var b sfnt.Buffer
	idx, err := fontData.GlyphIndex(&b, rune(renderString[0]))
	if err != nil {
		panic(err)
	}

	segments, err := fontData.LoadGlyph(&b, idx, fixed.I(ppem), nil)
	if err != nil {
		panic(err)
	}

	logrus.Infof("glyph loaded; %d segments for rune %+v\n", len(segments), rune(renderString[0]))

	bounds, _, err := fontData.GlyphBounds(&b, idx, fixed.I(ppem), font.HintingFull)
	if err != nil {
		panic(err)
	}

	app := NewApp()
	app.Initialize()

	app.loadBuffers(segments, bounds)

	app.winapp.DefaultMainLoop(shared.DefaultIgnoreInput, shared.DefaultIgnoreTick, app.drawFrame)

	app.Teardown()
	// Safe exit
}

type App struct {
	winapp   *shared.Win32App
	messages chan shared.WindowMessage

	vkctx.Context
	VulkanPipeline

	currentImage uint32

	vertexBuffer, indexBuffer             vk.Buffer
	vertexBufferMemory, indexBufferMemory vk.DeviceMemory

	indexCount int

	quadVertStart, quadIndsStart int
}

func NewApp() *App {
	c := make(chan shared.WindowMessage, 32)

	return &App{
		winapp:   shared.NewWin32App(c),
		messages: c,
	}
}

func (app *App) Initialize() {
	app.winapp.ClassName = "ttf-renderer"
	app.winapp.Width, app.winapp.Height = 800, 800
	app.winapp.Initialize("ttf-renderer")

	app.EnableApiLayers = append(app.EnableApiLayers, "VK_LAYER_KHRONOS_validation")
	app.EnableInstanceExtensions = app.winapp.GetRequiredInstanceExtensions()
	app.EnableDeviceExtensions = append(app.EnableDeviceExtensions, vk.KHR_SWAPCHAIN_EXTENSION_NAME)

	app.Context.Initialize(windows.Handle(app.winapp.HInstance), windows.HWND(app.winapp.HWnd))

	app.VulkanPipeline.Initialize(&app.Context)
}

func (app *App) Teardown() {
	vk.DeviceWaitIdle(app.ctx.Device)

	app.destroyBuffers()

	app.VulkanPipeline.Teardown()
	app.Context.Teardown()
}

func (app *App) drawFrame() {
	vk.WaitForFences(app.ctx.Device, []vk.Fence{app.ctx.InFlightFence}, true, ^uint64(0))

	var r vk.Result
	if r, app.currentImage = vk.AcquireNextImageKHR(app.ctx.Device, app.ctx.Swapchain, ^uint64(0), app.ctx.ImageAvailableSemaphore, vk.Fence(vk.NULL_HANDLE)); r != vk.SUCCESS {
		if r == vk.SUBOPTIMAL_KHR || r == vk.ERROR_OUT_OF_DATE_KHR {
			// app.vp.recreateSwapchain()
			return
		} else {
			panic("Could not acquire next image! " + r.String())
		}
	}

	vk.ResetFences(app.ctx.Device, []vk.Fence{app.ctx.InFlightFence})

	vk.ResetCommandBuffer(app.ctx.CommandBuffers[app.currentImage], 0)
	app.recordRenderingCommands(app.ctx.CommandBuffers[app.currentImage])

	// app.updateUniformBuffer(app.currentImage)

	submitInfo := vk.SubmitInfo{
		PWaitSemaphores:   []vk.Semaphore{app.ctx.ImageAvailableSemaphore},
		PWaitDstStageMask: []vk.PipelineStageFlags{vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT},
		PCommandBuffers:   []vk.CommandBuffer{app.ctx.CommandBuffers[app.currentImage]},
		PSignalSemaphores: []vk.Semaphore{app.ctx.RenderFinishedSemaphore},
	}

	if r := vk.QueueSubmit(app.ctx.GraphicsQueue, []vk.SubmitInfo{submitInfo}, app.ctx.InFlightFence); r != vk.SUCCESS {
		panic("Could not submit to graphics queue! " + r.String())
	}

	// Present the drawn image
	presentInfo := vk.PresentInfoKHR{
		PWaitSemaphores: []vk.Semaphore{app.ctx.RenderFinishedSemaphore},
		PSwapchains:     []vk.SwapchainKHR{app.ctx.Swapchain},
		PImageIndices:   []uint32{app.currentImage},
	}

	if r := vk.QueuePresentKHR(app.ctx.PresentQueue, &presentInfo); r != vk.SUCCESS && r != vk.SUBOPTIMAL_KHR && r != vk.ERROR_OUT_OF_DATE_KHR {
		panic("Could not submit to presentation queue! " + r.String())
	}

}

func (app *App) recordRenderingCommands(cb vk.CommandBuffer) {
	cbBeginInfo := vk.CommandBufferBeginInfo{
		Flags: vk.COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT,
	}

	colorCV, stencilCV := vk.ClearValue{}, vk.ClearValue{}
	ccv := vk.ClearColorValue{}
	ccv.AsTypeFloat32([4]float32{0, 0, 0, 1})
	colorCV.AsColor(ccv)

	stencilCV.AsDepthStencil(vk.ClearDepthStencilValue{
		Depth:   0,
		Stencil: 0,
	})

	rpBeginInfo := vk.RenderPassBeginInfo{
		RenderPass:  app.renderPass,
		Framebuffer: app.SwapChainFramebuffers[app.currentImage],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: app.SwapchainExtent,
		},
		PClearValues: []vk.ClearValue{colorCV, stencilCV},
	}

	vk.BeginCommandBuffer(cb, &cbBeginInfo)

	vk.CmdBeginRenderPass(cb, &rpBeginInfo, vk.SUBPASS_CONTENTS_INLINE)

	// bind vert, index bufs
	vk.CmdBindVertexBuffers(cb, 0, []vk.Buffer{app.vertexBuffer}, []vk.DeviceSize{0})
	vk.CmdBindIndexBuffer(cb, app.indexBuffer, 0, vk.INDEX_TYPE_UINT16)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[0]) // stencil pipeline
	vk.CmdDrawIndexed(cb, uint32(app.quadIndsStart), 1, 0, 0, 0)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[1]) // stencil quad portion pipeline
	vk.CmdDrawIndexed(cb, uint32(app.indexCount-app.quadIndsStart)-4, 1, uint32(app.quadIndsStart), int32(app.quadVertStart), 0)

	vk.CmdNextSubpass(cb, vk.SUBPASS_CONTENTS_INLINE)

	vk.CmdBindPipeline(cb, vk.PIPELINE_BIND_POINT_GRAPHICS, app.graphicsPipelines[2]) // Color pass

	// vk.CmdDrawIndexed(cb, 5, 1, uint32(app.indexCount)-5, 0, 0)
	vk.CmdDrawIndexed(cb, 4, 1, uint32(app.indexCount)-4, int32(app.quadVertStart), 0)

	// draw

	vk.CmdEndRenderPass(cb)
	vk.EndCommandBuffer(cb)

}
