package main

import (
	"os"
	"unsafe"

	"github.com/bbredesen/go-vk"
)

func (app *App_01) createGraphicsPipeline() {
	app.vertShaderModule = app.createShaderModule("shaders/vert.spv")
	app.fragShaderModule = app.createShaderModule("shaders/frag.spv")

	vertShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_VERTEX_BIT,
		Module:              app.vertShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	fragShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_FRAGMENT_BIT,
		Module:              app.fragShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	shaderStages := []vk.PipelineShaderStageCreateInfo{
		vertShaderStageCreateInfo, fragShaderStageCreateInfo,
	}

	vertexInputCreateInfo := vk.PipelineVertexInputStateCreateInfo{}

	inputAssemblyCreateInfo := vk.PipelineInputAssemblyStateCreateInfo{
		Topology:               vk.PRIMITIVE_TOPOLOGY_TRIANGLE_LIST,
		PrimitiveRestartEnable: false,
	}

	viewport := vk.Viewport{
		X:        0.0,
		Y:        0.0,
		Width:    float32(app.swapchainExtent.Width),
		Height:   float32(app.swapchainExtent.Height),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: app.swapchainExtent,
	}

	viewportStateCreateInfo := vk.PipelineViewportStateCreateInfo{
		PViewports: []vk.Viewport{viewport},
		PScissors:  []vk.Rect2D{scissor},
	}

	rasterizerCreateInfo := vk.PipelineRasterizationStateCreateInfo{
		DepthClampEnable:        false,
		RasterizerDiscardEnable: false,
		PolygonMode:             vk.POLYGON_MODE_FILL,
		LineWidth:               1.0,
		CullMode:                vk.CULL_MODE_BACK_BIT,
		FrontFace:               vk.FRONT_FACE_CLOCKWISE,
		DepthBiasEnable:         false,
	}

	multisampleCreateInfo := vk.PipelineMultisampleStateCreateInfo{
		SampleShadingEnable:  false,
		RasterizationSamples: vk.SAMPLE_COUNT_1_BIT,
		MinSampleShading:     1.0,
	}

	writeMask := vk.COLOR_COMPONENT_R_BIT |
		vk.COLOR_COMPONENT_G_BIT |
		vk.COLOR_COMPONENT_B_BIT |
		vk.COLOR_COMPONENT_A_BIT

	colorBlendAttachment := vk.PipelineColorBlendAttachmentState{
		ColorWriteMask: writeMask,
		BlendEnable:    false,

		// All ignored, b/c blend enable is false above
		// SrcColorBlendFactor: vk.BLEND_FACTOR_ONE,
		// DstColorBlendFactor: vk.BLEND_FACTOR_ZERO,
		// ColorBlendOp:        vk.BLEND_OP_ADD,
		// SrcAlphaBlendFactor: vk.BLEND_FACTOR_ONE,
		// DstAlphaBlendFactor: vk.BLEND_FACTOR_ZERO,
		// AlphaBlendOp:        vk.BLEND_OP_ADD,
	}

	colorBlendStateCreateInfo := vk.PipelineColorBlendStateCreateInfo{
		PAttachments: []vk.PipelineColorBlendAttachmentState{colorBlendAttachment},
	}

	// dynamicStateCreateInfo := vk.PipelineDynamicStateCreateInfo{
	// 	PDynamicStates: []vk.DynamicState{vk.DYNAMIC_STATE_VIEWPORT, vk.DYNAMIC_STATE_SCISSOR},
	// }

	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{}

	p, err := vk.CreatePipelineLayout(app.device, &pipelineLayoutCreateInfo, nil)
	if err != vk.SUCCESS {
		panic(err)
	}
	app.pipelineLayout = p

	pipelineCreateInfo := vk.GraphicsPipelineCreateInfo{
		PStages: shaderStages,
		// Fixed function stage information
		PVertexInputState:   &vertexInputCreateInfo,
		PInputAssemblyState: &inputAssemblyCreateInfo,
		PViewportState:      &viewportStateCreateInfo,
		PRasterizationState: &rasterizerCreateInfo,
		PMultisampleState:   &multisampleCreateInfo,
		PColorBlendState:    &colorBlendStateCreateInfo,
		PDepthStencilState:  &vk.PipelineDepthStencilStateCreateInfo{},

		PTessellationState: &vk.PipelineTessellationStateCreateInfo{},
		PDynamicState:      &vk.PipelineDynamicStateCreateInfo{}, // dynamicStateCreateInfo,

		Layout:     app.pipelineLayout,
		RenderPass: app.renderPass,
		Subpass:    0,
	}

	gp, err := vk.CreateGraphicsPipelines(
		app.device,
		0, // vk.NULL_HANDLE missing
		[]vk.GraphicsPipelineCreateInfo{pipelineCreateInfo},
		nil,
	)
	if err != nil {
		panic(err)
	}
	app.graphicsPipeline = gp[0]
}

func (app *App_01) createRenderPass() {

	colorAttachmentDescription := vk.AttachmentDescription{
		Format:  app.swapchainImageFormat,
		Samples: vk.SAMPLE_COUNT_1_BIT,
		LoadOp:  vk.ATTACHMENT_LOAD_OP_CLEAR,
		StoreOp: vk.ATTACHMENT_STORE_OP_STORE,

		StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_DONT_CARE,
		StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,

		InitialLayout: vk.IMAGE_LAYOUT_UNDEFINED,
		// InitialLayout: vk.IMAGE_LAYOUT_PRESENT_SRC_KHR,
		FinalLayout: vk.IMAGE_LAYOUT_PRESENT_SRC_KHR,
	}

	colorAttachmentRef := vk.AttachmentReference{
		Attachment: 0,
		Layout:     vk.IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL,
	}

	depthAttachmentRef := vk.AttachmentReference{
		Attachment: vk.ATTACHMENT_UNUSED,
	}

	subpassDescription := vk.SubpassDescription{
		PipelineBindPoint:       vk.PIPELINE_BIND_POINT_GRAPHICS,
		PColorAttachments:       []vk.AttachmentReference{colorAttachmentRef},
		PDepthStencilAttachment: &depthAttachmentRef,
	}

	// See
	// https://vulkan-tutorial.com/en/Drawing_a_triangle/Drawing/Rendering_and_presentation
	// https://registry.khronos.org/vulkan/specs/1.3-extensions/html/vkspec.html#VkSubpassDependency
	// This creates a dependency between this render pass and the "implied" subpass (the prior renderpass) before this
	// renderpass begins. It specifiesd that the the color attachment stage in the prior pass needs to be completed before we
	// attempt to write the attachment in this pass.

	dependency := vk.SubpassDependency{
		SrcSubpass:    vk.SUBPASS_EXTERNAL,
		DstSubpass:    0,
		SrcStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
		SrcAccessMask: 0,
		DstStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT,
		DstAccessMask: vk.ACCESS_COLOR_ATTACHMENT_WRITE_BIT,
	}

	renderPassCreateInfo := vk.RenderPassCreateInfo{
		PAttachments:  []vk.AttachmentDescription{colorAttachmentDescription},
		PSubpasses:    []vk.SubpassDescription{subpassDescription},
		PDependencies: []vk.SubpassDependency{dependency},
	}

	var err error
	if app.renderPass, err = vk.CreateRenderPass(app.device, &renderPassCreateInfo, nil); err != vk.SUCCESS {
		panic(err)
	}
}

func (app *App_01) createFramebuffers() {
	app.swapChainFramebuffers = make([]vk.Framebuffer, len(app.swapchainImageViews))

	for i, iv := range app.swapchainImageViews {
		framebufferCreateInfo := vk.FramebufferCreateInfo{
			RenderPass:   app.renderPass,
			PAttachments: []vk.ImageView{iv},
			Width:        app.swapchainExtent.Width,
			Height:       app.swapchainExtent.Height,
			Layers:       1,
		}

		fb, err := vk.CreateFramebuffer(app.device, &framebufferCreateInfo, nil)
		if err != nil {
			panic(err)
		}
		app.swapChainFramebuffers[i] = fb
	}
}

func (app *App_01) destroyFramebuffers() {
	for _, fb := range app.swapChainFramebuffers {
		vk.DestroyFramebuffer(app.device, fb, nil)
	}
	app.swapChainFramebuffers = nil
}

func (app *App_01) cleanupGraphicsPipeline() {

	vk.DestroyPipeline(app.device, app.graphicsPipeline, nil)
	vk.DestroyPipelineLayout(app.device, app.pipelineLayout, nil)

	vk.DestroyShaderModule(app.device, app.fragShaderModule, nil)
	app.fragShaderModule = vk.ShaderModule(vk.NULL_HANDLE)
	vk.DestroyShaderModule(app.device, app.vertShaderModule, nil)
	app.vertShaderModule = vk.ShaderModule(vk.NULL_HANDLE)

	vk.DestroyRenderPass(app.device, app.renderPass, nil)
}

func (app *App_01) createShaderModule(filename string) vk.ShaderModule {
	smCI := vk.ShaderModuleCreateInfo{
		CodeSize: 0,
		PCode:    new(uint32),
	}

	if dat, err := os.ReadFile(filename); err != nil {
		panic("Failed to read shader file " + filename + ": " + err.Error())
	} else {
		smCI.CodeSize = uintptr(len(dat))
		smCI.PCode = (*uint32)(unsafe.Pointer(&dat[0]))
	}

	if mod, err := vk.CreateShaderModule(app.device, &smCI, nil); err != nil {
		panic(err)
	} else {
		return mod
	}
}
