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

echo "Retrieve cc-runtime repository"
go get -d github.com/clearcontainers/runtime
pushd $GOPATH/src/github.com/clearcontainers/runtime

echo "Test cc-runtime"
make check
check_err

echo "Build cc-runtime"
make
check_err

echo "Install cc-runtime"
sudo make install
check_err

echo "Add cc-runtime as new Docker runtime"
sudo mkdir -p /etc/default
cat << EOF | sudo tee /etc/default/docker
DOCKER_OPTS="-D --add-runtime cor=/usr/libexec/clearcontainers/cc-runtime --default-runtime=cor"
EOF
check_err

echo "Restart docker service"
sudo service docker stop
sudo service docker start
check_err

echo "Install cc-runtime configuration.toml to /etc/clear-containers"
sudo mkdir -p /etc/clear-containers
sudo cp config/configuration_w_logs.toml /etc/clear-containers/configuration.toml
check_err

echo "Install cc-proxy service (/etc/init/cc-proxy.conf)"
sudo cp .ci/upstart-services/cc-proxy.conf /etc/init/
check_err

bash .ci/install_cc_env_ubuntu.sh
check_err

bash .ci/install_cc_proxy.sh
check_err

bash .ci/install_cc_shim.sh
check_err

bash .ci/install_virtcontainers.sh
check_err

popd

echo "Start cc-proxy service"
sudo service cc-proxy start
check_err
