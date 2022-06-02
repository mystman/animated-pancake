# 📽 🥞 Animated Pancake
*(What a wonderful auto-generated GitHub repo name by the way. 😎)*

This is an example for a REST API server built in Go with a pesistent JSON document store running in Kubernetes.

*Please note that this project for experimental purposes only. It did not meant to be more than a small proof of concept demonstartion.*

## Requirements:
* \*nix based machine with Make (*perhaps you can try to this on Windows as well; I just tried*)
* **Go** - https://go.dev/doc/install to compile the project (it used fmt and vet as well)
  * **lint** (https://github.com/golang/lint) is required if you want to compile the project locally (not required for the image build though)
* **Minikube** to spin up and deploy to local cluster - https://minikube.sigs.k8s.io/docs/
* **kubectl** to manage the Kubernetes manifests - https://kubernetes.io/docs/tasks/tools/

* **Docker** (https://www.docker.com/) or **Podman** (https://podman.io/) depending what you use for building container images, and/or container runtim for your cluster.  
(Also a nice to have in case you run into issues with making `minikube tunnel` for the LoadBalancer)

* **curl** for simple testing - https://curl.se (optional)
* **hey** for loadtesting - https://github.com/rakyll/hey (optional)



## Build and deploy:
* `make help` - displays the available Makefile targets
* `make minikube-up` - creates a local cluster and sets up the default namespace & context
* `minikube-deploy` - once the cluster is up, you can us this target to build the image and deploy the k8s artifacts
* `pod-info` - gives you and overview of the k8s resources
* `pod-log` - displays the log of the pod
* `load-test-post` and `load-test-get` - can be used for load testing
* **testing** folder contains sample requests with `curl`

After starting the `minikube tunnel` (and it running) for the LoadBalancer, the API endpoint is available localhost on the 6543 port.  


## Demo video:
[demo/Full_demo_1080p.mp4](https://github.com/mystman/animated-pancake/blob/main/demo/Full_demo_1080p.mp4)
[![Demo](https://github.com/mystman/animated-pancake/blob/main/demo/Thumbnail.png)](https://github.com/mystman/animated-pancake/blob/main/demo/Full_demo_1080p.mp4)
