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

	if r, layout := vk.CreateDescriptorSetLayout(app.device, &uboLayoutCI, nil); r != vk.SUCCESS {
		panic("Could not create descriptor set layout: " + r.String())
	} else {
		app.descriptorSetLayout = layout
	}
}

func (app *App_03) createUniformBuffers() {
	bufferSize := vk.DeviceSize(unsafe.Sizeof(UniformBufferObject{}))

	app.uboObjs = make([]UniformBufferObject, 2)
	app.uniformBuffers = make([]vk.Buffer, 2)
	app.uniformBufferMemories = make([]vk.DeviceMemory, 2)
	app.uniformBufferMapped = make([]*byte, 2)

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
		DescriptorCount: 2,
	}

	poolCI := vk.DescriptorPoolCreateInfo{
		MaxSets:    2,
		PPoolSizes: []vk.DescriptorPoolSize{poolSize},
	}

	var r vk.Result
	if r, app.descriptorPool = vk.CreateDescriptorPool(app.device, &poolCI, nil); r != vk.SUCCESS {
		panic("Could not create descriptor pool: " + r.String())
	}
}

func (app *App_03) createDescriptorSets() {

	allocInfo := vk.DescriptorSetAllocateInfo{
		DescriptorPool: app.descriptorPool,

		PSetLayouts: []vk.DescriptorSetLayout{app.descriptorSetLayout, app.descriptorSetLayout},
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
