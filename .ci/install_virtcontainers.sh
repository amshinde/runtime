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

echo "Retrieve virtcontainers repository"
go get -d github.com/containers/virtcontainers
pushd $GOPATH/src/github.com/containers/virtcontainers

echo "Build pause binary"
make pause
check_err

echo "Create /var/lib/clearcontainers/runtime/pause_bundle/bin"
sudo mkdir -p /var/lib/clearcontainers/runtime/bundles/pause_bundle/bin

echo "Install pause binary"
sudo cp pause/pause /var/lib/clearcontainers/runtime/bundles/pause_bundle/bin/
check_err

popd
