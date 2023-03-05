package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
	"golang.org/x/sys/windows"
)

func (app *App_03) createInstance() {
	appInfo := vk.ApplicationInfo{
		PApplicationName:   "02_vertex_buffers",
		ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
		EngineVersion:      vk.MAKE_VERSION(1, 0, 0),
		ApiVersion:         vk.MAKE_VERSION(1, 2, 0),
	}

	app.enableApiLayers = append(app.enableApiLayers, "VK_LAYER_KHRONOS_validation")

	icInfo := vk.InstanceCreateInfo{
		PApplicationInfo:        &appInfo,
		PpEnabledExtensionNames: app.enableInstanceExtensions,
		PpEnabledLayerNames:     app.enableApiLayers,
	}

	var r vk.Result
	r, app.instance = vk.CreateInstance(&icInfo, nil)

	if r != vk.SUCCESS {
		fmt.Println("Could not create instance!")
		panic(r.String())
	}
}

func (app *App_03) createSurface() {
	ci := vk.Win32SurfaceCreateInfoKHR{
		Hinstance: windows.Handle(app.Win32App.HInstance),
		Hwnd:      windows.HWND(app.Win32App.HWnd),
	}

	var r vk.Result
	r, app.surface = vk.CreateWin32SurfaceKHR(app.instance, &ci, nil)
	if r != vk.SUCCESS {
		fmt.Printf("Could not create surface!\n")
		panic(r)
	}
}

func (app *App_03) selectPhysicalDevice() {
	r, devices := vk.EnumeratePhysicalDevices(app.instance)
	if r != vk.SUCCESS {
		panic("Could not enumerate physical devices: " + r.String())
	}

	for _, dev := range devices {
		if app.isDeviceSuitable(dev) {
			app.physicalDevice = dev
			return
		}
	}

	panic("Could not find a suitable physical device!")
}

func (app *App_03) isDeviceSuitable(device vk.PhysicalDevice) bool {
	props := vk.GetPhysicalDeviceProperties(device)

	fmt.Printf("Found Physical Device:\n")
	fmt.Printf("  Device Name:\t\t%s\n", props.DeviceName)
	fmt.Printf("  Vendor/Device ID:\t0x%x / 0x%x \n", props.VendorID, props.DeviceID)
	fmt.Printf("  Device Type:\t\t%v\n", props.DeviceType)
	fmt.Printf("  Device API Version:\t%s\n", versionToString(props.ApiVersion))
	fmt.Printf("  Driver Version:\t%s\n", versionToString(props.DriverVersion))

	/* Suitability is:
	1) Support for the queue families we want to use (graphics)
	2) Support for the surface presentation extensions we want to use
	3) Support for swap chains // TODO
	*/

	extSupport := app.checkDeviceExtensionSupport(device)
	if !extSupport {
		panic("extensions not supported!")
	}
	inds := app.analyzeQueueFamilies(device)

	return extSupport && inds.graphicsIndex.HasValue() && inds.presentIndex.HasValue()
}

func (app *App_03) checkDeviceExtensionSupport(device vk.PhysicalDevice) bool {
	r, devExtensions := vk.EnumerateDeviceExtensionProperties(device, "")
	if r != vk.SUCCESS {
		panic(r.String() + ": Could not enumerate device extension properties!")
	}

	foundProps := make(map[string]bool, len(app.enableDeviceExtensions))
	// fmt.Println("Searching for device extensions:")
	for _, name := range app.enableDeviceExtensions {
		// init the found map with a false entry for each required extension
		foundProps[name] = false
	}

	for _, exProp := range devExtensions {
		if _, ok := foundProps[exProp.ExtensionName]; ok {
			foundProps[exProp.ExtensionName] = true
		}
	}

	haveAllExtensions := true
	for _, v := range foundProps {
		haveAllExtensions = haveAllExtensions && v
	}
	return haveAllExtensions
}

func (app *App_03) analyzeQueueFamilies(device vk.PhysicalDevice) queueFamIndices {
	qfp := vk.GetPhysicalDeviceQueueFamilyProperties(device)

	var inds queueFamIndices

	for i, p := range qfp {
		if (p.QueueFlags & vk.QUEUE_GRAPHICS_BIT) != 0 {
			inds.graphicsIndex.Set(uint32(i))
		}
		r, surf := vk.GetPhysicalDeviceSurfaceSupportKHR(device, uint32(i), app.surface)
		if r != vk.SUCCESS {
			panic(r)
		}

		if surf {
			inds.presentIndex.Set(uint32(i))
		}

		if inds.isComplete() {
			break
		}
	}
	return inds
}

func (app *App_03) createLogicalDevice() {
	// Re-analyze for the selected device
	qfInds := app.analyzeQueueFamilies(app.physicalDevice)

	// creates one or two entries, depending on how many queue families are needed
	uniqueQueueFams := make(map[uint32]bool)
	if !qfInds.graphicsIndex.HasValue() {
		panic("no graphics index found!")
	}
	uniqueQueueFams[qfInds.graphicsIndex.Value()] = true

	if !qfInds.presentIndex.HasValue() {
		panic("no presentation index found!")
	}
	uniqueQueueFams[qfInds.presentIndex.Value()] = true

	var dqCreateInfos []vk.DeviceQueueCreateInfo
	for k, v := range uniqueQueueFams {
		if v {
			// This family is selected as (one of possibly many) needed queues
			dqCreateInfos = append(dqCreateInfos,
				vk.DeviceQueueCreateInfo{
					QueueFamilyIndex: k,
					PQueuePriorities: []float32{1.0},
				})
		}
	}

	deviceFeatures := vk.PhysicalDeviceFeatures{
		// SamplerAnisotropy: true,
	}

	createInfo := vk.DeviceCreateInfo{
		PQueueCreateInfos:       dqCreateInfos,
		PEnabledFeatures:        &deviceFeatures,
		PpEnabledExtensionNames: app.enableDeviceExtensions,
		// EnabledLayerNames:     (deprecated)
	}

	r, device := vk.CreateDevice(app.physicalDevice, &createInfo, nil)

	if r != vk.SUCCESS {
		fmt.Printf("Logical device creation failed! (%s)\n", r.String())
		panic(r)
	}
	app.device = device

	app.graphicsQueueFamilyIndex = qfInds.graphicsIndex.Value()
	app.graphicsQueue = vk.GetDeviceQueue(app.device, qfInds.graphicsIndex.Value(), 0)
	app.presentQueueFamilyIndex = qfInds.presentIndex.Value()
	app.presentQueue = vk.GetDeviceQueue(app.device, qfInds.presentIndex.Value(), 0)
}

// ---------------------
type optUint32 struct {
	hasValue bool
	value    uint32
}

func (s *optUint32) Value() uint32 {
	if !s.hasValue {
		panic("Value not set on optUint32!")
	}
	return s.value
}

func (s *optUint32) Set(val uint32) {
	s.value = val
	s.hasValue = true
}

func (s *optUint32) HasValue() bool {
	return s.hasValue
}

type queueFamIndices struct {
	graphicsIndex optUint32
	presentIndex  optUint32
}

func (q *queueFamIndices) isComplete() bool {
	return q.graphicsIndex.HasValue() && q.presentIndex.HasValue()
}
