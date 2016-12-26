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
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/FlorianOtel/k8s-client/handler"

	"github.com/golang/glog"

	"github.com/FlorianOtel/client-go/kubernetes"
	"github.com/FlorianOtel/client-go/pkg/util/wait"

	"github.com/FlorianOtel/client-go/tools/clientcmd"

	// apiv1 "k8s.io/kubernetes/pkg/api/v1"
	apiv1 "github.com/FlorianOtel/client-go/pkg/api/v1"
	// "k8s.io/kubernetes/pkg/apis/extensions"
	// k8sfields "k8s.io/kubernetes/pkg/fields"
	// k8slabels "k8s.io/kubernetes/pkg/labels"
)

const errorLogLevel = 2

var (
	kubeconfig     = flag.String("kubeconfig", "./kubeconfig", "absolute path to the kubeconfig file")
	UseNetPolicies = false
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

	////////
	//////// Discover K8S API -- version, extensions
	////////

	sver, err := clientset.ServerVersion()

	glog.Infof("Kubernetes server details: %#v", *sver)

	//
	sres, err := clientset.ServerResources()

	for _, res := range sres {
		for _, apires := range res.APIResources {
			switch apires.Name {
			case "networkpolicies":
				glog.Infof(" ====> Found Kubernetes API server support for %#v. Available under / GroupVersion is: %#v . APIResource details: %#v", apires.Name, res.GroupVersion, apires)
				UseNetPolicies = true
			default:
				// glog.Infof("Kubernetes API Server discovery: API Server Resource:\n%#v\n", apires)
			}
		}
	}

	////////
	//////// Get initial list of Network Policies (if available)
	////////

	if UseNetPolicies {
		netpolicies, err := clientset.Extensions().NetworkPolicies("default").List(apiv1.ListOptions{})
		if err != nil {
			glog.Errorf("Error getting network policies for namespace: %s. Error: %s", "default", err)
		}

		for _, netpolicy := range netpolicies.Items {
			handler.JsonPrettyPrint("networkpolicy", &netpolicy)
		}
	}

	////////
	//////// Get initial list of Namespaces
	////////

	nses, err := clientset.Namespaces().List(apiv1.ListOptions{})

	//// Alternative, with selectors
	// ns, err := clientset.Namespaces().List(apiv1.ListOptions{LabelSelector: labels.Everything(), FieldSelector: fields.Everything()})

	if err != nil {
		glog.Errorf("Error listing namespaces. Error: %s", err)
	}

	fmt.Printf("-----> There are %d namespaces in the cluster\n", len(nses.Items))

	for _, ns := range (*nses).Items {
		handler.JsonPrettyPrint("namespace", &ns)
	}

	////////
	//////// Get Pod list
	////////

	pods, err := clientset.Core().Pods("").List(apiv1.ListOptions{})

	if err != nil {
		glog.Errorf("Error listing pods. Error: %s", err)
	}

	fmt.Printf("-----> There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range (*pods).Items {
		handler.JsonPrettyPrint("pod", &pod)
	}

	////////
	//////// Watch pods
	////////

	//Create a cache to store Pods
	// var store cache.Store
	// store, pController := handler.CreatePodController(clientset, "default", handler.PodCreated, handler.PodDeleted, handler.PodUpdated)

	_, pController := handler.CreatePodController(clientset, "default", handler.PodCreated, handler.PodDeleted, handler.PodUpdated)
	go pController.Run(wait.NeverStop)

	//Keep alive
	glog.Error(http.ListenAndServe(":8099", nil))

}
