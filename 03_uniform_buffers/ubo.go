package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/vkm"
)

type UniformBufferObject struct {
	model, view, proj vkm.Mat
}

func (app *App_03) createDescriptorSetLayout() {
	uboLayoutBinding := vk.DescriptorSetLayoutBinding{
		Binding:         0,
		DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
		StageFlags:      vk.SHADER_STAGE_VERTEX_BIT,
		DescriptorCount: 1,
	}

	uboLayoutCI := vk.DescriptorSetLayoutCreateInfo{
		PBindings: []vk.DescriptorSetLayoutBinding{uboLayoutBinding},
	}

	if layout, err := vk.CreateDescriptorSetLayout(app.device, &uboLayoutCI, nil); err != nil {
		panic("Could not create descriptor set layout: " + err.Error())
	} else {
		app.descriptorSetLayout = layout
	}
}

func (app *App_03) createUniformBuffers() {
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

		var err error
		if app.uniformBufferMapped[i], err = vk.MapMemory(app.device, app.uniformBufferMemories[i], 0, bufferSize, 0); err != nil {
			panic("Could not map memory for uniform buffer: " + err.Error())
		}

	}
}

func (app *App_03) cleanupUniformBuffers() {
	for i := range app.uniformBuffers {
		vk.UnmapMemory(app.device, app.uniformBufferMemories[i])
		vk.FreeMemory(app.device, app.uniformBufferMemories[i], nil)
		vk.DestroyBuffer(app.device, app.uniformBuffers[i], nil)
	}

	app.uniformBufferMapped = nil
	app.uniformBufferMemories = nil
	app.uniformBuffers = nil
}

func (app *App_03) createDescriptorPool() {
	poolSize := vk.DescriptorPoolSize{
		Typ:             vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
		DescriptorCount: uint32(len(app.swapchainImages)),
	}

	poolCI := vk.DescriptorPoolCreateInfo{
		MaxSets:    uint32(len(app.swapchainImages)),
		PPoolSizes: []vk.DescriptorPoolSize{poolSize},
	}

	var err error
	if app.descriptorPool, err = vk.CreateDescriptorPool(app.device, &poolCI, nil); err != nil {
		panic("Could not create descriptor pool: " + err.Error())
	}
}

func (app *App_03) createDescriptorSets() {

	allocInfo := vk.DescriptorSetAllocateInfo{
		DescriptorPool: app.descriptorPool,

		PSetLayouts: make([]vk.DescriptorSetLayout, len(app.swapchainImages)),
	}
	for i := range allocInfo.PSetLayouts {
		allocInfo.PSetLayouts[i] = app.descriptorSetLayout
	}

	var err error
	if app.descriptorSets, err = vk.AllocateDescriptorSets(app.device, &allocInfo); err != nil {
		panic("Could not allcoate descriptor sets: " + err.Error())
	}

	for i := range app.descriptorSets {
		bufInfo := vk.DescriptorBufferInfo{
			Buffer: app.uniformBuffers[i],
			Offset: 0,
			Rang:   vk.DeviceSize(unsafe.Sizeof(app.uboObjs[i])),
		}
		descriptorWrite := vk.WriteDescriptorSet{
			DstSet:          app.descriptorSets[i],
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorType:  vk.DESCRIPTOR_TYPE_UNIFORM_BUFFER,
			PBufferInfo:     []vk.DescriptorBufferInfo{bufInfo},
		}

		vk.UpdateDescriptorSets(app.device, []vk.WriteDescriptorSet{descriptorWrite}, nil)
	}

}
