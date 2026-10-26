package main

//go:generate glslc.exe shaders/quad_shader.vert -o shaders/quad_vert.spv
//go:generate glslc.exe shaders/shader.vert -o shaders/vert.spv
//go:generate glslc.exe shaders/quad_shader.frag -o shaders/quad_frag.spv
//go:generate glslc.exe shaders/shader.frag -o shaders/frag.spv

import (
	"os"
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/ttf-renderer/vkctx"
)

type VulkanPipeline struct {
	ctx               *vkctx.Context
	pipelineLayout    vk.PipelineLayout
	graphicsPipelines []vk.Pipeline

	// Renderpass
	renderPass vk.RenderPass

	stencilSubpass, colorSubpass     vk.SubpassDescription
	stencilImage, colorImage         vk.Image
	stencilMemory, colorMemory       vk.DeviceMemory
	stencilImageView, colorImageView vk.ImageView

	vertShaderModule, quadVertShaderModule, fragShaderModule, quadFragShaderModule vk.ShaderModule
}

func (vp *VulkanPipeline) Initialize(ctx *vkctx.Context) {
	vp.ctx = ctx
	vp.stencilImage, vp.stencilMemory = ctx.CreateImage(ctx.SwapchainExtent, vk.FORMAT_S8_UINT, vk.IMAGE_TILING_OPTIMAL, vk.IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	vp.stencilImageView = ctx.CreateImageView(vp.stencilImage, vk.FORMAT_S8_UINT, vk.IMAGE_ASPECT_STENCIL_BIT)

	vp.colorImage, vp.colorMemory = ctx.CreateImage(ctx.SwapchainExtent, vk.FORMAT_R32G32B32A32_SFLOAT, vk.IMAGE_TILING_OPTIMAL, vk.IMAGE_USAGE_COLOR_ATTACHMENT_BIT, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)
	vp.colorImageView = ctx.CreateImageView(vp.colorImage, vk.FORMAT_R32G32B32A32_SFLOAT, vk.IMAGE_ASPECT_COLOR_BIT)

	vp.CreateRenderPass()

	vp.CreateFramebuffers()

	vp.CreateGraphicsPipelines()
}

func (vp *VulkanPipeline) standardViewport() *vk.PipelineViewportStateCreateInfo {

	viewport := vk.Viewport{
		X:        0.0,
		Y:        0.0,
		Width:    float32(vp.ctx.SwapchainExtent.Width),
		Height:   float32(vp.ctx.SwapchainExtent.Height),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: vp.ctx.SwapchainExtent,
	}

	return &vk.PipelineViewportStateCreateInfo{
		PViewports: []vk.Viewport{viewport},
		PScissors:  []vk.Rect2D{scissor},
	}

}

func (vp *VulkanPipeline) CreateGraphicsPipelines() {
	// Two pipelines to build, three bindings
	// 1) Triangle fans for rough outline of shapes in stencil
	// 2) Triangles for quad curves in stencil
	// 3) Re-bind tri fans for color render

	vp.vertShaderModule = vp.createShaderModule("shaders/vert.spv")
	vp.fragShaderModule = vp.createShaderModule("shaders/frag.spv")

	vp.quadVertShaderModule = vp.createShaderModule("shaders/quad_vert.spv")
	vp.quadFragShaderModule = vp.createShaderModule("shaders/quad_frag.spv")

	p0_vertShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_VERTEX_BIT,
		Module:              vp.vertShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	p0_fragShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
		Stage:               vk.SHADER_STAGE_FRAGMENT_BIT,
		Module:              vp.fragShaderModule,
		PName:               "main",
		PSpecializationInfo: &vk.SpecializationInfo{},
	}

	p0_shaderStages := []vk.PipelineShaderStageCreateInfo{
		p0_vertShaderStageCreateInfo, p0_fragShaderStageCreateInfo,
	}

	// TODO
	bindings := []vk.VertexInputBindingDescription{
		{
			Binding: 0,
			Stride:  uint32(5 * unsafe.Sizeof(float32(0))),
		},
	}
	attrs := []vk.VertexInputAttributeDescription{
		{
			Location: 0,
			Binding:  0,
			Format:   vk.FORMAT_R32G32_SFLOAT,
			Offset:   0,
		},
		{
			Location: 1,
			Binding:  0,
			Format:   vk.FORMAT_R32G32B32_SFLOAT,
			Offset:   uint32(2 * unsafe.Sizeof(float32(0))),
		},
	}

	vertexInputCreateInfo := vk.PipelineVertexInputStateCreateInfo{
		PVertexBindingDescriptions:   bindings,
		PVertexAttributeDescriptions: attrs,
	}

	inputAssemblyCreateInfo := vk.PipelineInputAssemblyStateCreateInfo{
		Topology:               vk.PRIMITIVE_TOPOLOGY_TRIANGLE_FAN,
		PrimitiveRestartEnable: true,
	}

	rasterizerCreateInfo := vk.PipelineRasterizationStateCreateInfo{
		DepthClampEnable:        false,
		RasterizerDiscardEnable: false,
		PolygonMode:             vk.POLYGON_MODE_FILL,
		LineWidth:               1.0,
		CullMode:                vk.CULL_MODE_NONE,
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

	// depthStencilStateCreateInfo := vk.PipelineDepthStencilStateCreateInfo{
	// 	StencilTestEnable: true,
	// 	Front: vk.StencilOpState{
	// 		FailOp:      vk.STENCIL_OP_INCREMENT_AND_WRAP,
	// 		PassOp:      vk.STENCIL_OP_INCREMENT_AND_WRAP,
	// 		DepthFailOp: 0,
	// 		CompareOp:   vk.COMPARE_OP_ALWAYS,
	// 		CompareMask: 0x01,
	// 		WriteMask:   0xFF,
	// 		Reference:   1,
	// 	},
	// }
	// depthStencilStateCreateInfo.Back = depthStencilStateCreateInfo.Front

	viewport := vk.Viewport{
		X:        0.0,
		Y:        0.0,
		Width:    float32(vp.ctx.SwapchainExtent.Width),
		Height:   float32(vp.ctx.SwapchainExtent.Height),
		MinDepth: 0.0,
		MaxDepth: 1.0,
	}

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: vp.ctx.SwapchainExtent,
	}

	viewportStateCreateInfo := vk.PipelineViewportStateCreateInfo{
		PViewports: []vk.Viewport{viewport},
		PScissors:  []vk.Rect2D{scissor},
	}

	pipelineLayoutCreateInfo := vk.PipelineLayoutCreateInfo{
		PSetLayouts:         []vk.DescriptorSetLayout{},
		PPushConstantRanges: []vk.PushConstantRange{},
	}

	var r vk.Result
	if r, vp.pipelineLayout = vk.CreatePipelineLayout(vp.ctx.Device, &pipelineLayoutCreateInfo, nil); r != vk.SUCCESS {
		panic(r)
	}

	/* Rendering method:
	Each outline is drawn as a triangle fan from the start point of that outline. TTF fonts specify a clockwise winding rule
	for determining the "interior" of a glyph. The stencil operation below increments the stencil value on for front faces
	(i.e. drawn clockwise) and decrements the value for back faces. For the second subpass, the stencil test rule is changed
	to pass any fragments with a non-zero stencil value.
	*/

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

		PDepthStencilState: &depthStencilStateCreateInfo,

		Layout:     vp.pipelineLayout,
		RenderPass: vp.renderPass,
		Subpass:    0,
	}

	p1CreateInfo := pipelineCreateInfo
	ia := inputAssemblyCreateInfo
	p1CreateInfo.PInputAssemblyState = &ia
	ia.PrimitiveRestartEnable = false
	ia.Topology = vk.PRIMITIVE_TOPOLOGY_TRIANGLE_LIST

	p1CreateInfo.PStages = make([]vk.PipelineShaderStageCreateInfo, 2)
	copy(p1CreateInfo.PStages, p0_shaderStages)

	p1CreateInfo.PStages[0].Module = vp.quadVertShaderModule
	p1CreateInfo.PStages[1].Module = vp.quadFragShaderModule

	var tmp []vk.Pipeline
	if r, tmp = vk.CreateGraphicsPipelines(
		vp.ctx.Device,
		vk.PipelineCache(vk.NULL_HANDLE),
		[]vk.GraphicsPipelineCreateInfo{pipelineCreateInfo, p1CreateInfo},
		nil,
	); r != vk.SUCCESS {
		panic(r)
	}

	vp.graphicsPipelines = append(vp.graphicsPipelines, tmp...)

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

	if r, tmp = vk.CreateGraphicsPipelines(
		vp.ctx.Device,
		vk.PipelineCache(vk.NULL_HANDLE),
		[]vk.GraphicsPipelineCreateInfo{pipelineCreateInfo},
		nil,
	); r != vk.SUCCESS {
		panic(r)
	}

	vp.graphicsPipelines = append(vp.graphicsPipelines, tmp[0])

}

func (vp *VulkanPipeline) CreateRenderPass() {

	colorAttachmentDescription := vk.AttachmentDescription{
		Format:  vp.ctx.SwapchainImageFormat,
		Samples: vk.SAMPLE_COUNT_1_BIT,
		LoadOp:  vk.ATTACHMENT_LOAD_OP_CLEAR,
		StoreOp: vk.ATTACHMENT_STORE_OP_STORE,

		StencilLoadOp:  vk.ATTACHMENT_LOAD_OP_DONT_CARE,
		StencilStoreOp: vk.ATTACHMENT_STORE_OP_DONT_CARE,

		InitialLayout: vk.IMAGE_LAYOUT_UNDEFINED,
		FinalLayout:   vk.IMAGE_LAYOUT_PRESENT_SRC_KHR,
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

	var r vk.Result
	if r, vp.renderPass = vk.CreateRenderPass(vp.ctx.Device, &renderPassCreateInfo, nil); r != vk.SUCCESS {
		panic(r)
	}
}

func (vp *VulkanPipeline) CreateFramebuffers() {

	vp.ctx.SwapChainFramebuffers = make([]vk.Framebuffer, len(vp.ctx.SwapchainImageViews))

	for i, iv := range vp.ctx.SwapchainImageViews {
		framebufferCreateInfo := vk.FramebufferCreateInfo{
			RenderPass:   vp.renderPass,
			PAttachments: []vk.ImageView{iv, vp.stencilImageView},
			Width:        vp.ctx.SwapchainExtent.Width,
			Height:       vp.ctx.SwapchainExtent.Height,
			Layers:       1,
		}

		r, fb := vk.CreateFramebuffer(vp.ctx.Device, &framebufferCreateInfo, nil)
		if r != vk.SUCCESS {
			panic(r)
		}
		vp.ctx.SwapChainFramebuffers[i] = fb
	}
}

func (vp *VulkanPipeline) destroyFramebuffers() {
	for _, fb := range vp.ctx.SwapChainFramebuffers {
		vk.DestroyFramebuffer(vp.ctx.Device, fb, nil)
	}
	vp.ctx.SwapChainFramebuffers = nil
}

func (vp *VulkanPipeline) Teardown() {

	vk.DestroyShaderModule(vp.ctx.Device, vp.vertShaderModule, nil)
	vk.DestroyShaderModule(vp.ctx.Device, vp.fragShaderModule, nil)
	vk.DestroyShaderModule(vp.ctx.Device, vp.quadVertShaderModule, nil)
	vk.DestroyShaderModule(vp.ctx.Device, vp.quadFragShaderModule, nil)

	vk.DestroyImageView(vp.ctx.Device, vp.colorImageView, nil)
	vk.DestroyImageView(vp.ctx.Device, vp.stencilImageView, nil)

	vk.DestroyImage(vp.ctx.Device, vp.colorImage, nil)
	vk.DestroyImage(vp.ctx.Device, vp.stencilImage, nil)

	vk.FreeMemory(vp.ctx.Device, vp.colorMemory, nil)
	vk.FreeMemory(vp.ctx.Device, vp.stencilMemory, nil)

	vp.destroyFramebuffers()

	for _, gp := range vp.graphicsPipelines {
		vk.DestroyPipeline(vp.ctx.Device, gp, nil)
	}
	vp.graphicsPipelines = nil

	vk.DestroyPipelineLayout(vp.ctx.Device, vp.pipelineLayout, nil)
	vp.pipelineLayout = vk.PipelineLayout(vk.NULL_HANDLE)

	// vk.DestroyShaderModule(app.ctx.Device, app.fragShaderModule, nil)
	// app.fragShaderModule = vk.ShaderModule(vk.NULL_HANDLE)
	// vk.DestroyShaderModule(app.device, app.vertShaderModule, nil)
	// app.vertShaderModule = vk.ShaderModule(vk.NULL_HANDLE)

	vk.DestroyRenderPass(vp.ctx.Device, vp.renderPass, nil)
	vp.renderPass = vk.RenderPass(vk.NULL_HANDLE)
}

func (vp *VulkanPipeline) createShaderModule(filename string) vk.ShaderModule {
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

	if r, mod := vk.CreateShaderModule(vp.ctx.Device, &smCI, nil); r != vk.SUCCESS {
		panic(r)
	} else {
		return mod
	}
}
