# Copyright (c) 2017 Intel Corporation

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#!/bin/bash

function check_err {
	if [ $? -ne 0 ]
	then
		exit 1
	fi
}

echo "Add clear containers sources to apt list"
sudo sh -c "echo 'deb http://download.opensuse.org/repositories/home:/clearlinux:/preview:/clear-containers-2.1/xUbuntu_16.10/ /' >> /etc/apt/sources.list.d/cc-oci-runtime.list"
check_err

echo "Update apt repositories"
sudo apt-get update
check_err

echo "Install systemd"
sudo apt-get install -y systemd-services
check_err

echo "Install linux-container kernel"
sudo apt-get install -y --force-yes linux-container
check_err

echo "Install qemu-lite binary"
sudo apt-get install -y --force-yes qemu-lite
check_err

echo "Download clear containers image"
latest_version=$(curl -sL https://download.clearlinux.org/latest)
curl -LO "https://download.clearlinux.org/current/clear-${latest_version}-containers.img.xz"
check_err

echo "Extract clear containers image"
unxz clear-${latest_version}-containers.img.xz
check_err

sudo mkdir -p /usr/share/clear-containers/
echo "Install clear containers image"
sudo install --owner root --group root --mode 0755 clear-${latest_version}-containers.img /usr/share/clear-containers/
check_err

echo "Create symbolic link /usr/share/clear-containers/clear-containers.img"
sudo ln -fs /usr/share/clear-containers/clear-${latest_version}-containers.img /usr/share/clear-containers/clear-containers.img
check_err

echo "Install clear containers kernel 4.5-50"
sudo install --owner root --group root --mode 0755 .ci/kernel/vmlinux-4.5-50.container /usr/share/clear-containers/
check_err

echo "Create symbolic link /usr/share/clear-containers/vmlinux.container"
sudo ln -fs /usr/share/clear-containers/vmlinux-4.5-50.container /usr/share/clear-containers/vmlinux.container
check_err
