/*
Copyright 2016 The Kubernetes Authors.

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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const errorLogLevel = 2

var (
	kubeconfig = flag.String("kubeconfig", "./kubeconfig", "absolute path to the kubeconfig file")
)

func main() {

	flag.Parse()

	if len(os.Args) == 1 { // With no arguments, print default usage
		flag.PrintDefaults()
		os.Exit(0)
	}

	// glog.V(errorLogLevel).Infof("The given kubeconfig is: %s ", *kubeconfig)
	glog.Infof("The given kubeconfig is: %s ", *kubeconfig)

	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Errorf("Error parsing kubeconfig. Error: %s", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Error creating Kubernetes client. Error: %s", err)
	}

	for {
		pods, err := clientset.Core().Pods("").List(v1.ListOptions{})

		fmt.Printf(" -----> There are %d pods in the cluster\n", len(pods.Items))

		for nr, pod := range (*pods).Items {
			spec := pod.Spec
			meta := pod.ObjectMeta

			// JSON pretty-print the ObjectMeta
			jsonpodmeta, _ := json.MarshalIndent(meta, "", " ")

			// JSON pretty-print the PodSpec
			jsonpodspec, _ := json.MarshalIndent(spec, "", " ")

			fmt.Printf("====> Pod nr %d <====\n ######## Pod's ObjectMetadata ########\n%s\n ######## Pod's Spec ########\n%s\n\n ", nr, string(jsonpodmeta), string(jsonpodspec))
		}

		if err != nil {
			glog.Errorf("Error listing pods. Error: %s", err)
		}

		time.Sleep(3 * time.Second)
	}
}
