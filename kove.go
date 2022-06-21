package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/net/context"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	deviceCheckInterval = 5 * time.Second
)

// DeviceSpec defines a device that should be discovered and scheduled.
// DeviceSpec allows multiple host devices to be selected and scheduled fungibly under the same name.
// Furthermore, host devices can be composed into groups of device nodes that should be scheduled
// as an atomic unit.
type DeviceSpec struct {
	// Name is a unique string representing the kind of device this specification describes.
	Name string
	//// Groups is a list of groups of devices that should be scheduled under the same name.
	// Groups []*Group `json:"groups"`
	Count uint `json:"count,omitempty"`
}

// // Group represents a set of devices that should be grouped and mounted into a container together as one single meta-device.
// type Group struct {
// 	// Paths is the list of devices of which the device group consists.
// 	// Paths can be globs, in which case each device matched by the path will be schedulable `Count` times.
// 	// When the paths have differing cardinalities, that is, the globs match different numbers of devices,
// 	// the cardinality of each path is capped at the lowest cardinality.
// 	Paths []*Path `json:"paths"`
// 	// Count specifies how many times this group can be mounted concurrently.
// 	// When unspecified, Count defaults to 1.
// 	// Count uint `json:"count,omitempty"`
// }

// // Path represents a file path that should be discovered.
// type Path struct {
// 	// Path is the file path of a device in the host.
// 	Path string `json:"path"`
// 	// MountPath is the file path at which the host device should be mounted within the container.
// 	// When unspecified, MountPath defaults to the Path.
// 	MountPath string `json:"mountPath,omitempty"`
// 	// Permissions is the file-system permissions given to the mounted device.
// 	// Permissions applies only to mounts of type `Device`.
// 	// This can be one or more of:
// 	// * r - allows the container to read from the specified device.
// 	// * w - allows the container to write to the specified device.
// 	// * m - allows the container to create device files that do not yet exist.
// 	// When unspecified, Permissions defaults to mrw.
// 	Permissions string `json:"permissions,omitempty"`
// 	// ReadOnly specifies whether the path should be mounted read-only.
// 	// ReadOnly applies only to mounts of type `Mount`.
// 	ReadOnly bool `json:"readOnly,omitempty"`
// 	// Type describes what type of file-system node this Path represents and thus how it should be mounted.
// 	// When unspecified, Type defaults to Device.
// 	Type PathType `json:"type"`
// }

// // PathType represents the kinds of file-system nodes that can be scheduled.
// type PathType string

type device struct {
	v1beta1.Device
	deviceSpecs []*v1beta1.DeviceSpec
}

// GenericPlugin is a plugin for generic devices that can:
// * be found using a file path; and
// * mounted and used without special logic.
type GenericPlugin struct {
	ds      *DeviceSpec
	devices map[string]device
	logger  log.Logger
	mu      sync.Mutex
}

// Allocate assigns generic devices to a Pod.
func (gp *GenericPlugin) Allocate(_ context.Context, req *v1beta1.AllocateRequest) (*v1beta1.AllocateResponse, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	res := &v1beta1.AllocateResponse{
		ContainerResponses: make([]*v1beta1.ContainerAllocateResponse, 0, len(req.ContainerRequests)),
	}
	for _, r := range req.ContainerRequests {
		resp := new(v1beta1.ContainerAllocateResponse)
		// Add all requested devices to to response.
		for _, id := range r.DevicesIDs {
			d, ok := gp.devices[id]
			if !ok {
				return nil, fmt.Errorf("requested device does not exist %q", id)
			}
			if d.Health != v1beta1.Healthy {
				return nil, fmt.Errorf("requested device is not healthy %q", id)
			}
			resp.Devices = append(resp.Devices, d.deviceSpecs...)
			// resp.Mounts = append(resp.Mounts, d.mounts...)
		}
		res.ContainerResponses = append(res.ContainerResponses, resp)
	}
	// gp.allocationsCounter.Add(float64(len(res.ContainerResponses)))
	return res, nil
}

// GetDevicePluginOptions always returns an empty response.
func (gp *GenericPlugin) GetDevicePluginOptions(_ context.Context, _ *v1beta1.Empty) (*v1beta1.DevicePluginOptions, error) {
	level.Info(gp.logger).Log("msg", "starting GetDevicePluginOptions")
	return &v1beta1.DevicePluginOptions{}, nil
}

// ListAndWatch lists all devices and then refreshes every deviceCheckInterval.
func (gp *GenericPlugin) ListAndWatch(_ *v1beta1.Empty, stream v1beta1.DevicePlugin_ListAndWatchServer) error {
	<-time.After(deviceCheckInterval)
	rand.Seed(time.Now().UnixNano())
	gp.ds.Count = uint(rand.Intn(9999))
	return nil
}

// PreStartContainer always returns an empty response.
func (gp *GenericPlugin) PreStartContainer(_ context.Context, _ *v1beta1.PreStartContainerRequest) (*v1beta1.PreStartContainerResponse, error) {
	level.Info(gp.logger).Log("msg", "starting PreStartContainer")
	return &v1beta1.PreStartContainerResponse{}, nil
}

// GetPreferredAllocation always returns an empty response.
func (gp *GenericPlugin) GetPreferredAllocation(context.Context, *v1beta1.PreferredAllocationRequest) (*v1beta1.PreferredAllocationResponse, error) {
	level.Info(gp.logger).Log("msg", "starting GetPreferredAllocation")
	return &v1beta1.PreferredAllocationResponse{}, nil
}
