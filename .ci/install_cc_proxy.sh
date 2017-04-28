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

echo "Retrieve cc-proxy repository"
go get -d github.com/clearcontainers/proxy
pushd $GOPATH/src/github.com/clearcontainers/proxy

echo "Build cc-proxy"
make
check_err

echo "Install cc-proxy"
sudo make install
check_err

popd
