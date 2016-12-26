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

	"github.com/golang/glog"

	"github.com/FlorianOtel/client-go/kubernetes"
	"github.com/FlorianOtel/client-go/pkg/runtime"

	"github.com/FlorianOtel/client-go/tools/clientcmd"

	// apiv1 "k8s.io/kubernetes/pkg/api/v1"
	apiv1 "github.com/FlorianOtel/client-go/pkg/api/v1"
	apiv1beta1 "github.com/FlorianOtel/client-go/pkg/apis/extensions/v1beta1"
	// "k8s.io/kubernetes/pkg/apis/extensions"
	// k8sfields "k8s.io/kubernetes/pkg/fields"
	// k8slabels "k8s.io/kubernetes/pkg/labels"
)

const errorLogLevel = 2

var (
	kubeconfig     = flag.String("kubeconfig", "./kubeconfig", "absolute path to the kubeconfig file")
	UseNetPolicies = false
)

// Pretty Prints (JSON) for a Kubernetes API object:
// - The "ObjectMeta"  is common to all the API objects and is handled identically, disregarding of the underlying type
// - The "Spec" is specific to each reasource and is handled on per-object specific basis (even if the field -- "Spec" -- is named the same for all objects)

func JsonPrettyPrint(resource string, obj runtime.Object) error {
	var meta apiv1.ObjectMeta
	var jsonmeta, jsonspec []byte
	var err error

	switch resource {
	case "pod":
		meta = obj.(*apiv1.Pod).ObjectMeta
		jsonspec, err = json.MarshalIndent(obj.(*apiv1.Pod).Spec, "", " ")
	case "namespace":
		meta = obj.(*apiv1.Namespace).ObjectMeta
		jsonspec, err = json.MarshalIndent(obj.(*apiv1.Namespace).Spec, "", " ")
	case "networkpolicy":
		meta = obj.(*apiv1beta1.NetworkPolicy).ObjectMeta
		jsonspec, err = json.MarshalIndent(obj.(*apiv1beta1.NetworkPolicy).Spec, "", " ")
	default:
		glog.Errorf("Don't know how to pretty-print resource: %s", resource)
	}

	// If any of the JSON marshalling of the object-specific specs above returned an error
	if err != nil {
		return err
	}

	// JSON pretty-print the ObjectMeta -- unlike "Spect" it's the same for all type of objects (has it's own type)
	jsonmeta, err = json.MarshalIndent(meta, "", " ")

	if err != nil {
		return err
	}

	fmt.Printf("====> %s <====\n ######## %s ObjectMetadata ########\n%s\n ######## %s Spec ########\n%s\n\n ", resource, resource, string(jsonmeta), resource, string(jsonspec))

	return err
}

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
			JsonPrettyPrint("networkpolicy", &netpolicy)
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

	fmt.Printf(" -----> There are %d namespaces in the cluster\n", len(nses.Items))

	for _, ns := range (*nses).Items {
		JsonPrettyPrint("namespace", &ns)
	}

	////////
	//////// Get Pod list
	////////

	pods, err := clientset.Core().Pods("").List(apiv1.ListOptions{})

	if err != nil {
		glog.Errorf("Error listing pods. Error: %s", err)
	}

	fmt.Printf(" -----> There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range (*pods).Items {
		JsonPrettyPrint("pod", &pod)
	}

}
