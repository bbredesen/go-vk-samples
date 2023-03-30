package main

import (
	"fmt"

	"github.com/bbredesen/go-vk"
)

func (app *App_06) createInstance() {
	appInfo := vk.ApplicationInfo{
		PApplicationName:   "06_loading_models",
		ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
		EngineVersion:      vk.MAKE_VERSION(1, 0, 0),
		ApiVersion:         vk.MAKE_VERSION(1, 2, 0),
	}

	// app.enableApiLayers = append(app.enableApiLayers, "VK_LAYER_KHRONOS_validation")

	icInfo := vk.InstanceCreateInfo{
		PApplicationInfo:        &appInfo,
		PpEnabledExtensionNames: app.enableInstanceExtensions,
		PpEnabledLayerNames:     app.enableApiLayers,
	}

	var err error
	app.instance, err = vk.CreateInstance(&icInfo, nil)

	if err != nil {
		fmt.Println("Could not create instance!")
		panic(err.Error())
	}
}

func (app *App_06) createSurface() {
	var err error
	app.surface, err = app.DelegateCreateSurface(app.instance)
	if err != nil {
		panic("Could not create surface: " + err.Error())
	}
}

func (app *App_06) selectPhysicalDevice() {
	devices, err := vk.EnumeratePhysicalDevices(app.instance)
	if err != nil {
		panic("Could not enumerate physical devices: " + err.Error())
	}

	for _, dev := range devices {
		if app.isDeviceSuitable(dev) {
			app.physicalDevice = dev
			// app.pdmp = vk.GetPhysicalDeviceMemoryProperties(app.physicalDevice)
			return
		}
	}

	panic("Could not find a suitable physical device!")
}

func (app *App_06) isDeviceSuitable(device vk.PhysicalDevice) bool {
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
	4) Support for sampler anisotropy
	*/

	features := vk.GetPhysicalDeviceFeatures(device)

	extSupport := app.checkDeviceExtensionSupport(device)
	if !extSupport {
		panic("extensions not supported!")
	}
	inds := app.analyzeQueueFamilies(device)

	return extSupport && inds.graphicsIndex.HasValue() && inds.presentIndex.HasValue() && features.SamplerAnisotropy
}

func (app *App_06) checkDeviceExtensionSupport(device vk.PhysicalDevice) bool {
	devExtensions, err := vk.EnumerateDeviceExtensionProperties(device, "")
	if err != nil {
		panic(err.Error() + ": Could not enumerate device extension properties!")
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

func (app *App_06) analyzeQueueFamilies(device vk.PhysicalDevice) queueFamIndices {
	qfp := vk.GetPhysicalDeviceQueueFamilyProperties(device)

	var inds queueFamIndices

	for i, p := range qfp {
		if (p.QueueFlags & vk.QUEUE_GRAPHICS_BIT) != 0 {
			inds.graphicsIndex.Set(uint32(i))
		}
		surf, err := vk.GetPhysicalDeviceSurfaceSupportKHR(device, uint32(i), app.surface)
		if err != nil {
			panic(err)
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

func (app *App_06) createLogicalDevice() {
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
		SamplerAnisotropy: true,
	}

	createInfo := vk.DeviceCreateInfo{
		PQueueCreateInfos:       dqCreateInfos,
		PEnabledFeatures:        &deviceFeatures,
		PpEnabledExtensionNames: app.enableDeviceExtensions,
		// EnabledLayerNames:     (deprecated)
	}

	device, err := vk.CreateDevice(app.physicalDevice, &createInfo, nil)

	if err != nil {
		fmt.Printf("Logical device creation failed! (%s)\n", err.Error())
		panic(err)
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
