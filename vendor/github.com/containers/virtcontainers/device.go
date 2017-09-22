//
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package virtcontainers

import (
	"fmt"
	"os"
	"strings"
)

const (
	sysIOMMUPath = "/sys/kernel/iommu_groups"
)

// Device that must be available in the container.
type Device struct {
	// Full path to the device inside container.
	Path string

	// Type of device: c, b, u or p
	Type string

	// Major, minor numbers for device.
	Major int64
	Minor int64

	// FileMode permission bits for the device.
	FileMode *os.FileMode

	// id of the device owner.
	UID uint32

	// id of the device group.
	GID uint32
}

// VFIODevice is a vfio device meant to be passed to the hypervisor
// to be used by the Virtual Machine.
type VFIODevice struct {
	// Bus-Device-Function for pci device
	BDF string
}

// isVFIODeviceGroup checks if the device provided is a vfio group
func (d *Device) isVFIODeviceGroup() bool {
	vfioPath := "/dev/vfio/"

	if strings.HasPrefix(d.Path, vfioPath) && len(d.Path) > len(vfioPath) {
		return true
	}

	return false
}

// bdf returns the BDF of pci device
// Expected input strng format is [<domain>]:[<bus>][<slot>].[<func>] eg. 0000:02:10.0
func bdf(deviceSysStr string) (string, error) {
	tokens := strings.Split(deviceSysStr, ":")

	if len(tokens) != 3 {
		return "", fmt.Errorf("Incorrect number of tokens found while parsing bdf for device : %s", deviceSysStr)
	}

	tokens = strings.SplitN(deviceSysStr, ":", 2)
	return tokens[1], nil
}
