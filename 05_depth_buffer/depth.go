package main

import (
	"github.com/bbredesen/go-vk"
)

func (app *App_05) createDepthResources() {
	format := app.findDepthFormat()

	app.depthImage, app.depthImageMemory = app.createImage(
		app.swapchainExtent.Width, app.swapchainExtent.Height,
		format,
		vk.IMAGE_TILING_OPTIMAL,
		vk.IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT,
		vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT,
	)

	app.depthImageView = app.createImageView(app.depthImage, format, vk.IMAGE_ASPECT_DEPTH_BIT)

	// app.transitionImageLayout(app.depthImage, format, vk.IMAGE_LAYOUT_UNDEFINED, vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL)
}

func (app *App_05) findDepthFormat() vk.Format {
	return vk.FORMAT_D32_SFLOAT

	// depthFormats := []vk.Format{vk.FORMAT_D32_SFLOAT, vk.FORMAT_D32_SFLOAT_S8_UINT, vk.FORMAT_D24_UNORM_S8_UINT}

	// return app.findSupportedFormat(depthFormats,
	// 	vk.IMAGE_TILING_OPTIMAL,
	// 	vk.FORMAT_FEATURE_DEPTH_STENCIL_ATTACHMENT_BIT,
	// )
}

func (app *App_05) findSupportedFormat(candidates []vk.Format, tiling vk.ImageTiling, features vk.FormatFeatureFlags) vk.Format {
	for _, f := range candidates {
		pdfp := vk.GetPhysicalDeviceFormatProperties(app.physicalDevice, f)

		if tiling == vk.IMAGE_TILING_OPTIMAL && features&pdfp.OptimalTilingFeatures == features {
			return f
		}

		if tiling == vk.IMAGE_TILING_LINEAR && features&pdfp.LinearTilingFeatures == features {
			return f
		}
	}

	panic("Could not find a fully supported format and tiling combination.")
}

func formatHasStencilComponent(f vk.Format) bool {
	return f == vk.FORMAT_D32_SFLOAT_S8_UINT || f == vk.FORMAT_D24_UNORM_S8_UINT || f == vk.FORMAT_D16_UNORM_S8_UINT
}

func (app *App_05) cleanupDepthResources() {
	vk.DestroyImageView(app.device, app.depthImageView, nil)
	vk.DestroyImage(app.device, app.depthImage, nil)
	vk.FreeMemory(app.device, app.depthImageMemory, nil)
}
