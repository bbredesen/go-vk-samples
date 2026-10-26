package main

import (
	"os"
	"unsafe"

	"github.com/bbredesen/go-vk"
)

func (app *App_10) createGraphicsPipeline() {
	app.vertShaderModule = app.createShaderModule("shaders/vert.spv")
	app.fragShaderModule = app.createShaderModule("shaders/frag.spv")

	app.quadVertShaderModule = app.createShaderModule("shaders/quad_vert.spv")
	app.quadFragShaderModule = app.createShaderModule("shaders/quad_frag.spv")

	p0_vertShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_VERTEX_BIT,
		Module:              app.vertShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	p0_fragShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_FRAGMENT_BIT,
		Module:              app.fragShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	p0_shaderStages := []vk.PipelineShaderStageCreateInfo{
		p0_vertShaderStageCreateInfo, p0_fragShaderStageCreateInfo,
	}

	vertexInputCreateInfo := vk.PipelineVertexInputStateCreateInfo{
		PNext:                        nil,
		Flags:                        0,
		PVertexBindingDescriptions:   []vk.VertexInputBindingDescription{getVertexBindingDescription()},
		PVertexAttributeDescriptions: getVertexAttributeDescriptions(),
	}

	inputAssemblyCreateInfo := vk.PipelineInputAssemblyStateCreateInfo{
		Topology:               vk.PRIMITIVE_TOPOLOGY_TRIANGLE_FAN,
		PrimitiveRestartEnable: true,
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
		CullMode:                vk.CULL_MODE_NONE,
		FrontFace:               vk.FRONT_FACE_COUNTER_CLOCKWISE,
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

	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{
		PSetLayouts: []vk.DescriptorSetLayout{}, //app.descriptorSetLayout
	}

	p, err := vk.CreatePipelineLayout(app.device, &pipelineLayoutCreateInfo, nil)
	if err != nil {
		panic(err)
	}
	app.pipelineLayout = p

	depthStencilStateCreateInfo := vk.PipelineDepthStencilStateCreateInfo{
		StencilTestEnable: true,
		DepthTestEnable:   false,
		Front: vk.StencilOpState{
			PassOp:    vk.STENCIL_OP_INCREMENT_AND_WRAP,
			CompareOp: vk.COMPARE_OP_ALWAYS,
			WriteMask: 0xFF,
			Reference: 1,
		},
		Back: vk.StencilOpState{
			PassOp:    vk.STENCIL_OP_DECREMENT_AND_WRAP,
			CompareOp: vk.COMPARE_OP_ALWAYS,
			WriteMask: 0xFF,
			Reference: 1,
		},
	}

	pipelineCreateInfo := vk.GraphicsPipelineCreateInfo{
		PStages: p0_shaderStages,
		// Fixed function stage information
		PVertexInputState:   &vertexInputCreateInfo,
		PInputAssemblyState: &inputAssemblyCreateInfo,
		PViewportState:      &viewportStateCreateInfo,
		PRasterizationState: &rasterizerCreateInfo,
		PMultisampleState:   &multisampleCreateInfo,
		PColorBlendState:    &colorBlendStateCreateInfo,
		PDepthStencilState:  &depthStencilStateCreateInfo,

		PTessellationState: &vk.PipelineTessellationStateCreateInfo{},
		PDynamicState:      &vk.PipelineDynamicStateCreateInfo{}, // dynamicStateCreateInfo,

		Layout:     app.pipelineLayout,
		RenderPass: app.renderPass,
		Subpass:    0,
	}

	// Modify a copy of the initial pipeline for the second pass
	p1CreateInfo := pipelineCreateInfo
	ia := inputAssemblyCreateInfo
	p1CreateInfo.PInputAssemblyState = &ia
	ia.PrimitiveRestartEnable = false
	ia.Topology = vk.PRIMITIVE_TOPOLOGY_TRIANGLE_LIST

	p1CreateInfo.PStages = make([]vk.PipelineShaderStageCreateInfo, 2)
	copy(p1CreateInfo.PStages, p0_shaderStages)

	p1CreateInfo.PStages[0].Module = app.quadVertShaderModule
	p1CreateInfo.PStages[1].Module = app.quadFragShaderModule

	gp, err := vk.CreateGraphicsPipelines(
		app.device,
		vk.PipelineCache(vk.NULL_HANDLE),
		[]vk.GraphicsPipelineCreateInfo{pipelineCreateInfo, p1CreateInfo},
		nil,
	)
	if err != nil {
		panic(err)
	}
	app.graphicsPipelines = append(app.graphicsPipelines, gp...)

	depthStencilStateCreateInfo.DepthTestEnable = false
	// depthStencilStateCreateInfo.StencilTestEnable = false
	depthStencilStateCreateInfo.Front = vk.StencilOpState{
		// FailOp: vk.STENCIL_OP_REPLACE,
		// DepthFailOp: vk.STENCIL_OP_KEEP,
		// PassOp: vk.STENCIL_OP_KEEP,
		CompareOp:   vk.COMPARE_OP_NOT_EQUAL,
		CompareMask: 0xFF,
		// WriteMask:   0xFF,
		Reference: 0,
	}
	depthStencilStateCreateInfo.Back = depthStencilStateCreateInfo.Front
	pipelineCreateInfo.Subpass = 1

	if gp, err = vk.CreateGraphicsPipelines(
		app.device,
		vk.PipelineCache(vk.NULL_HANDLE),
		[]vk.GraphicsPipelineCreateInfo{pipelineCreateInfo},
		nil,
	); err != vk.SUCCESS {
		panic(err)
	}

	app.graphicsPipelines = append(app.graphicsPipelines, gp[0])

}

func (app *App_10) createRenderPass() {

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

	stencilAttachmentDescription := vk.AttachmentDescription{
		Format:  vk.FORMAT_S8_UINT,
		Samples: vk.SAMPLE_COUNT_1_BIT,

		// Applies to depth component
		LoadOp:  vk.ATTACHMENT_LOAD_OP_CLEAR,
		StoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,

		StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_CLEAR,
		StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,

		InitialLayout: vk.IMAGE_LAYOUT_UNDEFINED,
		FinalLayout:   vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL,
	}

	stencilAttachmentRef := vk.AttachmentReference{
		Attachment: 1,
		Layout:     vk.IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL,
	}

	stencilSubpassDescription := vk.SubpassDescription{
		PipelineBindPoint:       vk.PIPELINE_BIND_POINT_GRAPHICS,
		PColorAttachments:       []vk.AttachmentReference{},
		PDepthStencilAttachment: &stencilAttachmentRef,
	}

	colorSubpassDescription := vk.SubpassDescription{
		PipelineBindPoint:       vk.PIPELINE_BIND_POINT_GRAPHICS,
		PColorAttachments:       []vk.AttachmentReference{colorAttachmentRef},
		PDepthStencilAttachment: &stencilAttachmentRef,
	}

	// See
	// https://vulkan-tutorial.com/en/Drawing_a_triangle/Drawing/Rendering_and_presentation
	// https://registry.khronos.org/vulkan/specs/1.3-extensions/html/vkspec.html#VkSubpassDependency
	// This creates an execution/timing dependency between this render pass and the "implied" subpass (the prior renderpass) before this
	// renderpass begins. It specifiesd that the the color attachment output and depth testing stages in the prior pass
	// need to be completed before we attempt to write the color and depth attachments in this pass.

	dependencyToStencil := vk.SubpassDependency{
		SrcSubpass:    vk.SUBPASS_EXTERNAL,
		DstSubpass:    0,
		SrcStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT | vk.PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT,
		SrcAccessMask: 0,
		DstStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT | vk.PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT,
		DstAccessMask: vk.ACCESS_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT,
	}

	dependencyToColor := vk.SubpassDependency{
		SrcSubpass:    0,
		DstSubpass:    1,
		SrcStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT | vk.PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT,
		SrcAccessMask: 0,
		DstStageMask:  vk.PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT | vk.PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT,
		DstAccessMask: vk.ACCESS_COLOR_ATTACHMENT_WRITE_BIT | vk.ACCESS_DEPTH_STENCIL_ATTACHMENT_READ_BIT,
	}

	renderPassCreateInfo := vk.RenderPassCreateInfo{
		PAttachments:  []vk.AttachmentDescription{colorAttachmentDescription, stencilAttachmentDescription},
		PSubpasses:    []vk.SubpassDescription{stencilSubpassDescription, colorSubpassDescription},
		PDependencies: []vk.SubpassDependency{dependencyToStencil, dependencyToColor},
	}

	var err error
	if app.renderPass, err = vk.CreateRenderPass(app.device, &renderPassCreateInfo, nil); err != nil {
		panic(err)
	}
}

func (app *App_10) createFramebuffers() {
	app.swapChainFramebuffers = make([]vk.Framebuffer, len(app.swapchainImageViews))

	for i, iv := range app.swapchainImageViews {
		framebufferCreateInfo := vk.FramebufferCreateInfo{
			RenderPass:   app.renderPass,
			PAttachments: []vk.ImageView{iv, app.stencilImageView},
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

func (app *App_10) destroyFramebuffers() {
	for _, fb := range app.swapChainFramebuffers {
		vk.DestroyFramebuffer(app.device, fb, nil)
	}
	app.swapChainFramebuffers = nil
}

func (app *App_10) cleanupGraphicsPipeline() {

	for _, gp := range app.graphicsPipelines {
		vk.DestroyPipeline(app.device, gp, nil)
	}
	app.graphicsPipelines = nil

	vk.DestroyPipelineLayout(app.device, app.pipelineLayout, nil)

	vk.DestroyShaderModule(app.device, app.fragShaderModule, nil)
	app.fragShaderModule = vk.ShaderModule(vk.NULL_HANDLE)
	vk.DestroyShaderModule(app.device, app.vertShaderModule, nil)
	app.vertShaderModule = vk.ShaderModule(vk.NULL_HANDLE)
	vk.DestroyShaderModule(app.device, app.quadFragShaderModule, nil)
	app.quadFragShaderModule = vk.ShaderModule(vk.NULL_HANDLE)
	vk.DestroyShaderModule(app.device, app.quadVertShaderModule, nil)
	app.quadVertShaderModule = vk.ShaderModule(vk.NULL_HANDLE)

	vk.DestroyRenderPass(app.device, app.renderPass, nil)
}

func (app *App_10) createShaderModule(filename string) vk.ShaderModule {
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
