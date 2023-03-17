package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/vkm"
)

type UniformBufferObject struct {
	model, view, proj vkm.Mat
}

func (app *App_04) createDescriptorSetLayout() {
	uboLayoutBinding := vk.DescriptorSetLayoutBinding{
		Binding:         0,
		DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
		StageFlags:      vk.SHADER_STAGE_VERTEX_BIT,
		DescriptorCount: 1,
	}

	samplerLayoutBinding := vk.DescriptorSetLayoutBinding{
		Binding:         1,
		DescriptorType:  vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
		StageFlags:      vk.SHADER_STAGE_FRAGMENT_BIT,
		DescriptorCount: 1,
	}

	uboLayoutCI := vk.DescriptorSetLayoutCreateInfo{
		PBindings: []vk.DescriptorSetLayoutBinding{uboLayoutBinding, samplerLayoutBinding},
	}

	if r, layout := vk.CreateDescriptorSetLayout(app.device, &uboLayoutCI, nil); r != vk.SUCCESS {
		panic("Could not create descriptor set layout: " + r.String())
	} else {
		app.descriptorSetLayout = layout
	}
}

func (app *App_04) createUniformBuffers() {
	bufferSize := vk.DeviceSize(unsafe.Sizeof(UniformBufferObject{}))

	app.uboObjs = make([]UniformBufferObject, len(app.swapchainImages))
	app.uniformBuffers = make([]vk.Buffer, len(app.swapchainImages))
	app.uniformBufferMemories = make([]vk.DeviceMemory, len(app.swapchainImages))
	app.uniformBufferMapped = make([]*byte, len(app.swapchainImages))

	for i := range app.uniformBuffers {
		app.uboObjs[i] = UniformBufferObject{
			view: vkm.LookAt(vkm.NewPt(2, 2, 2), vkm.Origin(), vkm.UnitVecZ()),
			proj: vkm.PerspectiveDeg(45.0, float32(app.swapchainExtent.Width)/float32(app.swapchainExtent.Height), 0.1, 10.0),
		}

		app.uniformBuffers[i], app.uniformBufferMemories[i] = app.createBuffer(vk.BUFFER_USAGE_UNIFORM_BUFFER_BIT, bufferSize, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)

		var r vk.Result
		if r, app.uniformBufferMapped[i] = vk.MapMemory(app.device, app.uniformBufferMemories[i], 0, bufferSize, 0); r != vk.SUCCESS {
			panic("Could not map memory for uniform buffer: " + r.String())
		}

	}
}

func (app *App_04) cleanupUniformBuffers() {
	for i := range app.uniformBuffers {
		vk.UnmapMemory(app.device, app.uniformBufferMemories[i])
		vk.FreeMemory(app.device, app.uniformBufferMemories[i], nil)
		vk.DestroyBuffer(app.device, app.uniformBuffers[i], nil)
	}

	app.uniformBufferMapped = nil
	app.uniformBufferMemories = nil
	app.uniformBuffers = nil
}

func (app *App_04) createDescriptorPool() {
	uboPoolSize := vk.DescriptorPoolSize{
		Typ:             vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
		DescriptorCount: uint32(len(app.swapchainImages)),
	}
	samplerPoolSize := vk.DescriptorPoolSize{
		Typ:             vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
		DescriptorCount: uint32(len(app.swapchainImages)),
	}

	poolCI := vk.DescriptorPoolCreateInfo{
		MaxSets:    uint32(len(app.swapchainImages)),
		PPoolSizes: []vk.DescriptorPoolSize{uboPoolSize, samplerPoolSize},
	}

	var r vk.Result
	if r, app.descriptorPool = vk.CreateDescriptorPool(app.device, &poolCI, nil); r != vk.SUCCESS {
		panic("Could not create descriptor pool: " + r.String())
	}
}

func (app *App_04) createDescriptorSets() {
	allocInfo := vk.DescriptorSetAllocateInfo{
		DescriptorPool: app.descriptorPool,

		PSetLayouts: make([]vk.DescriptorSetLayout, len(app.swapchainImages)),
	}
	for i := range allocInfo.PSetLayouts {
		allocInfo.PSetLayouts[i] = app.descriptorSetLayout
	}

	var r vk.Result
	if r, app.descriptorSets = vk.AllocateDescriptorSets(app.device, &allocInfo); r != vk.SUCCESS {
		panic("Could not allcoate descriptor sets: " + r.String())
	}

	for i := range app.descriptorSets {
		bufInfo := vk.DescriptorBufferInfo{
			Buffer: app.uniformBuffers[i],
			Offset: 0,
			Rang:   vk.DeviceSize(unsafe.Sizeof(app.uboObjs[i])),
		}

		imageInfo := vk.DescriptorImageInfo{
			Sampler:     app.textureSampler,
			ImageView:   app.textureImageView,
			ImageLayout: vk.IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL,
		}

		descriptorWrites := []vk.WriteDescriptorSet{
			{
				DstSet:          app.descriptorSets[i],
				DstBinding:      0,
				DstArrayElement: 0,
				DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
				PBufferInfo:     []vk.DescriptorBufferInfo{bufInfo},
			},
			{
				DstSet:          app.descriptorSets[i],
				DstBinding:      1,
				DstArrayElement: 0,
				DescriptorType:  vk.DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER,
				PImageInfo:      []vk.DescriptorImageInfo{imageInfo},
			},
		}

		vk.UpdateDescriptorSets(app.device, descriptorWrites, nil)
	}

}
