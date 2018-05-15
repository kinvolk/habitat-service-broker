// Copyright (c) 2018 Chef Software Inc. and/or applicable contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"flag"
	"fmt"
	"os"
	"testing"

	brokerFramework "github.com/kinvolk/habitat-service-broker/test/e2e/framework"
)

var framework *brokerFramework.Framework

func TestMain(m *testing.M) {
	var (
		err  error
		code int
	)

	image := flag.String("image", "", "habitat service broker image, 'kinvolk/habitat-service-broker'")
	kubeconfig := flag.String("kubeconfig", "", "path to kube config file")
	externalIP := flag.String("ip", "", "external ip, eg. minikube ip")
	flag.Parse()

	if framework, err = brokerFramework.Setup(*image, *kubeconfig, *externalIP); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code = m.Run()

	_ = framework.TearDown()

	os.Exit(code)
}
