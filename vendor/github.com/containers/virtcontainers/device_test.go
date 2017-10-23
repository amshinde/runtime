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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fileMode0640 = os.FileMode(0640)

func TestIsVFIO(t *testing.T) {
	type testData struct {
		path     string
		expected bool
	}

	data := []testData{
		{"/dev/vfio/16", true},
		{"/dev/vfio/1", true},
		{"/dev/vfio/", false},
		{"/dev/vfio", false},
		{"/dev/vf", false},
		{"/dev", false},
		{"/dev/vfio/vfio", false},
		{"/dev/vfio/vfio/12", false},
	}

	for _, d := range data {
		isVFIO := isVFIO(d.path)
		assert.Equal(t, d.expected, isVFIO)
	}
}

func TestIsBlock(t *testing.T) {
	type testData struct {
		path     string
		expected bool
	}

	data := []testData{
		{"/dev/sda", true},
		{"/dev/sdbb", true},
		{"/dev/hda", true},
		{"/dev/hdb", true},
		{"/dev/vf", false},
		{"/dev/vdj", true},
		{"/dev/vdzzz", true},
		{"/dev/ida", false},
		{"/dev/ida/", false},
		{"/dev/ida/c0d0p10", true},
	}

	for _, d := range data {
		isBlock := isBlock(d.path)
		assert.Equal(t, d.expected, isBlock)
	}
}

func TestCreateDevice(t *testing.T) {
	devInfo := DeviceInfo{
		HostPath: "/dev/vfio/8",
	}

	device := createDevice(devInfo)
	_, ok := device.(*VFIODevice)
	assert.True(t, ok)

	devInfo.HostPath = "/dev/sda"
	device = createDevice(devInfo)
	_, ok = device.(*BlockDevice)
	assert.True(t, ok)

	devInfo.HostPath = "/dev/tty"
	device = createDevice(devInfo)
	_, ok = device.(*GenericDevice)
	assert.True(t, ok)
}

func TestNewDevices(t *testing.T) {
	savedSysDevPrefix := sysDevPrefix

	major := int64(252)
	minor := int64(3)

	tmpDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	sysDevPrefix = tmpDir
	defer func() {
		os.RemoveAll(tmpDir)
		sysDevPrefix = savedSysDevPrefix
	}()

	path := "/dev/vfio/2"
	deviceInfo := DeviceInfo{
		ContainerPath: "",
		Major:         major,
		Minor:         minor,
		UID:           2,
		GID:           2,
		DevType:       "c",
	}

	_, err = newDevices([]DeviceInfo{deviceInfo})
	assert.NotNil(t, err)

	format := strconv.FormatInt(major, 10) + ":" + strconv.FormatInt(minor, 10)
	ueventPathPrefix := filepath.Join(sysDevPrefix, "char", format)
	ueventPath := filepath.Join(ueventPathPrefix, "uevent")

	// Return true for non-existent /sys/dev path.
	deviceInfo.ContainerPath = path
	_, err = newDevices([]DeviceInfo{deviceInfo})
	assert.Nil(t, err)

	err = os.MkdirAll(ueventPathPrefix, dirMode)
	assert.Nil(t, err)

	// Should return error for bad data in uevent file
	content := []byte("nonkeyvaluedata")
	err = ioutil.WriteFile(ueventPath, content, fileMode0640)
	assert.Nil(t, err)

	_, err = newDevices([]DeviceInfo{deviceInfo})
	assert.NotNil(t, err)

	content = []byte("MAJOR=252\nMINOR=3\nDEVNAME=vfio/2")
	err = ioutil.WriteFile(ueventPath, content, fileMode0640)
	assert.Nil(t, err)

	devices, err := newDevices([]DeviceInfo{deviceInfo})
	assert.Nil(t, err)

	assert.Equal(t, len(devices), 1)
	vfioDev, ok := devices[0].(*VFIODevice)
	assert.True(t, ok)
	assert.Equal(t, vfioDev.DeviceInfo.HostPath, path)
	assert.Equal(t, vfioDev.DeviceInfo.ContainerPath, path)
	assert.Equal(t, vfioDev.DeviceInfo.DevType, "c")
	assert.Equal(t, vfioDev.DeviceInfo.Major, major)
	assert.Equal(t, vfioDev.DeviceInfo.Minor, minor)
	assert.Equal(t, vfioDev.DeviceInfo.UID, uint32(2))
	assert.Equal(t, vfioDev.DeviceInfo.GID, uint32(2))
}

func TestGetBDF(t *testing.T) {
	type testData struct {
		deviceStr   string
		expectedBDF string
	}

	data := []testData{
		{"0000:02:10.0", "02:10.0"},
		{"0000:0210.0", ""},
		{"test", ""},
		{"", ""},
	}

	for _, d := range data {
		deviceBDF, err := getBDF(d.deviceStr)
		assert.Equal(t, d.expectedBDF, deviceBDF)
		if d.expectedBDF == "" {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestAttachVFIODevice(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	os.RemoveAll(tmpDir)

	testFDIOGroup := "2"
	testDeviceBDFPath := "0000:00:1c.0"

	devicesDir := filepath.Join(tmpDir, testFDIOGroup, "devices")
	err = os.MkdirAll(devicesDir, dirMode)
	assert.Nil(t, err)

	deviceFile := filepath.Join(devicesDir, testDeviceBDFPath)
	_, err = os.Create(deviceFile)
	assert.Nil(t, err)

	savedIOMMUPath := sysIOMMUPath
	sysIOMMUPath = tmpDir

	defer func() {
		sysIOMMUPath = savedIOMMUPath
	}()

	path := filepath.Join(vfioPath, testFDIOGroup)
	deviceInfo := DeviceInfo{
		HostPath:      path,
		ContainerPath: path,
		DevType:       "c",
	}

	device := createDevice(deviceInfo)
	_, ok := device.(*VFIODevice)
	assert.True(t, ok)

	hypervisor := &mockHypervisor{}
	err = device.attach(hypervisor)
	assert.Nil(t, err)

	err = device.detach(hypervisor)
	assert.Nil(t, err)
}

func TestAttachGenericDevice(t *testing.T) {
	path := "/dev/tty2"
	deviceInfo := DeviceInfo{
		HostPath:      path,
		ContainerPath: path,
		DevType:       "c",
	}

	device := createDevice(deviceInfo)
	_, ok := device.(*GenericDevice)
	assert.True(t, ok)

	hypervisor := &mockHypervisor{}
	err := device.attach(hypervisor)
	assert.Nil(t, err)

	err = device.detach(hypervisor)
	assert.Nil(t, err)
}

func TestAttachBlockDevice(t *testing.T) {
	path := "/dev/hda"
	deviceInfo := DeviceInfo{
		HostPath:      path,
		ContainerPath: path,
		DevType:       "c",
	}

	device := createDevice(deviceInfo)
	_, ok := device.(*BlockDevice)
	assert.True(t, ok)

	hypervisor := &mockHypervisor{}
	err := device.attach(hypervisor)
	assert.Nil(t, err)

	err = device.detach(hypervisor)
	assert.Nil(t, err)
}
