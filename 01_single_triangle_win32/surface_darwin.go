//go:build darwin

package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
)

func (app *App_01) createSurface() {
	var r vk.Result
	ci := vk.MetalSurfaceCreateInfoEXT{
		PLayer: (*vk.CAMetalLayer)(unsafe.Pointer(app.windowHandle)),
	}

	if r, app.surface = vk.CreateMetalSurfaceEXT(app.instance, &ci, nil); r != vk.SUCCESS {
		panic(r)
	}
}
