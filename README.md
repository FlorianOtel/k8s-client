# k8s-client -- simple Kubernetes client for listing / listening to Kubernetes events 

A simple standalone Kubernetes client based on [Go client for Kubernetes](https://github.com/kubernetes/client-go "client-go").


Its purpose is to serve as scaffolding code for:

* Discovering server API capabilities: Listing API constructs

* Listing Kubernetes constructs. Currently supports: Pods, Services, Namespaces, Network Policies. 

* Watching CRUD operations for those constructs & performing actions on those operations. Currently: Only listing the object details (`ObjectMeta` and object pecific `Specs`) 


Based on various bits and pieces, and code samples found on the net. Thanks to all involved but particularly to our dear friends at [Aporeto](https://www.aporeto.com) and their [Trireme](https://www.aporeto.com/trireme/) OSS project.


## Usage 

Requires a `kubeconfig` file -- e.g. `kubelet.kubeconfig` for authentication to the Kubernetes cluster

Build accordingl (`go build`) then start as: 

```
/k8s-client  -alsologtostderr -kubeconfig /path/to/kubelet.kubeconfig 
```

## Comments, Questions, Issues, Contributions

Via Github. TIA for any


