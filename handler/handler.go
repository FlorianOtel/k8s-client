/*

Attribution for this code: Our dearest friends at Aporeto -- see https://www.aporeto.com/trireme/.
Original code: https://github.com/aporeto-inc/trireme-kubernetes/blob/master/kubernetes/handler.go

*/

package handler

import (
	"github.com/golang/glog"
	//
	"github.com/FlorianOtel/client-go/kubernetes"
	apiv1 "github.com/FlorianOtel/client-go/pkg/api/v1"
	apiv1beta1 "github.com/FlorianOtel/client-go/pkg/apis/extensions/v1beta1"
	"github.com/FlorianOtel/client-go/pkg/fields"
	"github.com/FlorianOtel/client-go/pkg/runtime"
	"github.com/FlorianOtel/client-go/tools/cache"
)

const errorLogLevel = 2

// CreateResourceController creates a controller for a specific ressource and namespace.
// The parameter function will be called on Add/Delete/Update events
func CreateResourceController(client cache.Getter, resource string, namespace string, apiStruct runtime.Object, selector fields.Selector,
	addFunc func(addedApiStruct interface{}), deleteFunc func(deletedApiStruct interface{}), updateFunc func(oldApiStruct, updatedApiStruct interface{})) (cache.Store, *cache.Controller) {

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		DeleteFunc: deleteFunc,
		UpdateFunc: updateFunc,
	}

	listWatch := cache.NewListWatchFromClient(client, resource, namespace, selector)
	store, controller := cache.NewInformer(listWatch, apiStruct, 0, handlers)
	return store, controller
}

// CreateNamespaceController creates a controller specifically for Namespaces.
func (c *kubernetes.Clientset) CreateNamespaceController(
	addFunc func(addedApiStruct *apiv1.Namespace) error, deleteFunc func(deletedApiStruct *apiv1.Namespace) error, updateFunc func(oldApiStruct, updatedApiStruct *apiv1.Namespace) error) (cache.Store, *cache.Controller) {

	return CreateResourceController(c.KubeClient().Core().RESTClient(), "namespaces", "", &apiv1.Namespace{}, fields.Everything(),
		func(addedApiStruct interface{}) {
			if err := addFunc(addedApiStruct.(*apiv1.Namespace)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Add NameSpace: %s ", err)
			}
		},
		func(deletedApiStruct interface{}) {
			if err := deleteFunc(deletedApiStruct.(*apiv1.Namespace)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Delete NameSpace: %s ", err)

			}
		},
		func(oldApiStruct, updatedApiStruct interface{}) {
			if err := updateFunc(oldApiStruct.(*apiv1.Namespace), updatedApiStruct.(*apiv1.Namespace)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Update NameSpace: %s ", err)

			}
		})
}

// CreateLocalPodController creates a controller specifically for Pods.
func (c *kubernetes.Clientset) CreateLocalPodController(namespace string,
	addFunc func(addedApiStruct *apiv1.Pod) error, deleteFunc func(deletedApiStruct *apiv1.Pod) error, updateFunc func(oldApiStruct, updatedApiStruct *apiv1.Pod) error) (cache.Store, *cache.Controller) {

	return CreateResourceController(c.KubeClient().Core().RESTClient(), "pods", namespace, &apiv1.Pod{}, c.localNodeSelector(),
		func(addedApiStruct interface{}) {
			if err := addFunc(addedApiStruct.(*apiv1.Pod)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Add Pod: %s ", err)
			}
		},
		func(deletedApiStruct interface{}) {
			if err := deleteFunc(deletedApiStruct.(*apiv1.Pod)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Delete Pod: %s ", err)
			}
		},
		func(oldApiStruct, updatedApiStruct interface{}) {
			if err := updateFunc(oldApiStruct.(*apiv1.Pod), updatedApiStruct.(*apiv1.Pod)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Update Pod: %s ", err)
			}
		})
}

// CreateNetworkPoliciesController creates a controller specifically for NetworkPolicies.
func (c *kubernetes.Clientset) CreateNetworkPoliciesController(namespace string,
	addFunc func(addedApiStruct *apiv1beta1.NetworkPolicy) error, deleteFunc func(deletedApiStruct *apiv1beta1.NetworkPolicy) error, updateFunc func(oldApiStruct, updatedApiStruct *apiv1beta1.NetworkPolicy) error) (cache.Store, *cache.Controller) {
	return CreateResourceController(c.Extensions().RESTClient(), "networkpolicies", namespace, &apiv1beta1.NetworkPolicy{}, fields.Everything(),
		func(addedApiStruct interface{}) {
			if err := addFunc(addedApiStruct.(*apiv1beta1.NetworkPolicy)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Add NetworkPolicy: %s ", err)
			}
		},
		func(deletedApiStruct interface{}) {
			if err := deleteFunc(deletedApiStruct.(*apiv1beta1.NetworkPolicy)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Delete NetworkPolicy: %s ", err)
			}
		},
		func(oldApiStruct, updatedApiStruct interface{}) {
			if err := updateFunc(oldApiStruct.(*apiv1beta1.NetworkPolicy), updatedApiStruct.(*apiv1beta1.NetworkPolicy)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Update NetworkPolicy: %s ", err)
			}
		})
}

// CreateNodeController creates a controller specifically for Nodes.
func (c *kubernetes.Clientset) CreateNodeController(
	addFunc func(addedApiStruct *apiv1.Node) error, deleteFunc func(deletedApiStruct *apiv1.Node) error, updateFunc func(oldApiStruct, updatedApiStruct *apiv1.Node) error) (cache.Store, *cache.Controller) {
	return CreateResourceController(c.KubeClient().Core().RESTClient(), "nodes", "", &apiv1.Node{}, fields.Everything(),
		func(addedApiStruct interface{}) {
			if err := addFunc(addedApiStruct.(*apiv1.Node)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Add Node: %s ", err)
			}
		},
		func(deletedApiStruct interface{}) {
			if err := deleteFunc(deletedApiStruct.(*apiv1.Node)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Delete Node: %s ", err)
			}
		},
		func(oldApiStruct, updatedApiStruct interface{}) {
			if err := updateFunc(oldApiStruct.(*apiv1.Node), updatedApiStruct.(*apiv1.Node)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Update Node: %s ", err)
			}
		})
}

// CreateServiceController creates a controller specifically for Services.
func (c *kubernetes.Clientset) CreateServiceController(
	addFunc func(addedApiStruct *apiv1.Service) error, deleteFunc func(deletedApiStruct *apiv1.Service) error, updateFunc func(oldApiStruct, updatedApiStruct *apiv1.Service) error) (cache.Store, *cache.Controller) {
	return CreateResourceController(c.KubeClient().Core().RESTClient(), "services", "", &apiv1.Service{}, fields.Everything(),
		func(addedApiStruct interface{}) {
			if err := addFunc(addedApiStruct.(*apiv1.Service)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Add service: %s ", err)
			}
		},
		func(deletedApiStruct interface{}) {
			if err := deleteFunc(deletedApiStruct.(*apiv1.Service)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Delete service: %s ", err)
			}
		},
		func(oldApiStruct, updatedApiStruct interface{}) {
			if err := updateFunc(oldApiStruct.(*apiv1.Service), updatedApiStruct.(*apiv1.Service)); err != nil {
				glog.V(errorLogLevel).Infof("Error while handling Update service: %s ", err)
			}
		})
}
