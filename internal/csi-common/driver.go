/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package csicommon

import (
	"github.com/ceph/ceph-csi/internal/util/log"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

// CSIDriver stores driver information.
type CSIDriver struct {
	name    string
	nodeID  string
	version string

	// instance is the instance ID that is unique to an instance of CSI, used when sharing
	// ceph clusters across CSI instances, to differentiate omap names per CSI instance.
	instance string

	// topology constraints that this nodeserver will advertise
	topology          map[string]string
	capabilities      []*csi.ControllerServiceCapability
	groupCapabilities []*csi.GroupControllerServiceCapability
	vc                []*csi.VolumeCapability_AccessMode
}

// NewCSIDriver Creates a NewCSIDriver object. Assumes vendor
// version is equal to driver version &  does not support optional
// driver plugin info manifest field. Refer to CSI spec for more details.
func NewCSIDriver(name, v, nodeID, instance string) *CSIDriver {
	if name == "" {
		klog.Errorf("Driver name missing")

		return nil
	}

	if nodeID == "" {
		klog.Errorf("NodeID missing")

		return nil
	}
	// TODO version format and validation
	if v == "" {
		klog.Errorf("Version argument missing")

		return nil
	}

	if instance == "" {
		klog.Errorf("Instance argument missing")

		return nil
	}

	driver := CSIDriver{
		name:     name,
		version:  v,
		nodeID:   nodeID,
		instance: instance,
	}

	return &driver
}

// GetInstance returns the instance identification of the CSI driver.
func (d *CSIDriver) GetInstanceID() string {
	return d.instance
}

// ValidateControllerServiceRequest validates the controller
// plugin capabilities.
func (d *CSIDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, capability := range d.capabilities {
		if c == capability.GetRpc().GetType() {
			return nil
		}
	}

	return status.Error(codes.InvalidArgument, c.String())
}

// AddControllerServiceCapabilities stores the controller capabilities
// in driver object.
func (d *CSIDriver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	csc := make([]*csi.ControllerServiceCapability, 0, len(cl))

	for _, c := range cl {
		log.DefaultLog("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}

	d.capabilities = csc
}

// AddVolumeCapabilityAccessModes stores volume access modes.
func (d *CSIDriver) AddVolumeCapabilityAccessModes(
	vc []csi.VolumeCapability_AccessMode_Mode,
) []*csi.VolumeCapability_AccessMode {
	vca := make([]*csi.VolumeCapability_AccessMode, 0, len(vc))
	for _, c := range vc {
		log.DefaultLog("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.vc = vca

	return vca
}

// GetVolumeCapabilityAccessModes returns access modes.
func (d *CSIDriver) GetVolumeCapabilityAccessModes() []*csi.VolumeCapability_AccessMode {
	return d.vc
}

// AddControllerServiceCapabilities stores the group controller capabilities
// in driver object.
func (d *CSIDriver) AddGroupControllerServiceCapabilities(cl []csi.GroupControllerServiceCapability_RPC_Type) {
	csc := make([]*csi.GroupControllerServiceCapability, 0, len(cl))

	for _, c := range cl {
		log.DefaultLog("Enabling group controller service capability: %v", c.String())
		csc = append(csc, NewGroupControllerServiceCapability(c))
	}

	d.groupCapabilities = csc
}

// ValidateGroupControllerServiceRequest validates the group controller
// plugin capabilities.
func (d *CSIDriver) ValidateGroupControllerServiceRequest(c csi.GroupControllerServiceCapability_RPC_Type) error {
	if c == csi.GroupControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, capability := range d.groupCapabilities {
		if c == capability.GetRpc().GetType() {
			return nil
		}
	}

	return status.Error(codes.InvalidArgument, c.String())
}
