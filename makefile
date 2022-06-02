#=======| Defaults |=======
.DEFAULT_GOAL := help

#=======| Project settings |=======
BINARY_NAME=animated-pancake
BUILD_DIR=./bin

#=======| Build flags |=======
BUILD_NAME = "${BINARY_NAME}"
BUILD_VERSION := 1.0
BUILD_DATETIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ") # ISO-8601 date

LDFLAGS=-ldflags="-w -s \
-X 'main.buildName=${BUILD_NAME}' \
-X 'main.buildVersion=${BUILD_VERSION}' \
-X 'main.buildDate=${BUILD_DATETIME}'"

#=======| Help |=======
## Display help
.PHONY: help
help:  
	@awk 'BEGIN {FS = ":.*##"; printf "\n\"$(BINARY_NAME) $(BUILD_VERSION)\" makefile\n\nUsage:  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-17s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	
#=======| Compile |=======
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: fmt
	golint ./...

.PHONY: vet
vet: fmt
	go vet ./...

#=======| Run |=======
.PHONY: run 
run: lint vet test  # [ Run ] Run the main.go (with LDFLAGS)
	go run cmd/main.go ${LDFLAGS}
	
.PHONY: debug
debug: lint vet # [ Run ] Compile and run with gcflags
	go run -gcflags '-m -l' cmd/main.go

#=======| Build |=======
.PHONY: build 
build: lint vet ## [ Build ] Build the project binary
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/main.go

.PHONY: build_release
build_release: # [ Build ] To be used from Dockerfile for building a release
	go build ${LDFLAGS} -o ${BUILD_DIR}/service cmd/main.go



.PHONY: clean
clean: # [ Build ] Clean the project binary and trace files
	go clean
	rm -rf ./bin
	rm -f copy_trace.out

#=======| Test |=======
.PHONY: test
test: ## [ Test ] Execute tests
	go test ./... -v

#=======| Debug - Trace & Benchmark |=======
.PHONY: bench
bench: # [ Debug ] Run benchmark
	go test -v -run="none" -bench="BenchmarkMain" -benchmem	

.PHONY: bench_long 
bench_long: # [ Debug ] Run long benchmark
	go test -v -run="none" -bench="BenchmarkMain" -benchmem -benchtime=3s -count 3 ./cmd

#=======| Containers |=======
.PHONY: docker-build
docker-build: ## [ Container] Build container image with Docker
	docker build -t ${BUILD_NAME}:${BUILD_VERSION} .


.PHONY: minikube-docker
minikube-docker: ## [ Container] Build container image with Minikube's Docker
	./minikube-image-build.sh ${BUILD_NAME} ${BUILD_VERSION}

#=======| Cluster |=======
CLUSTER_NAME := pancake-cluster
NAPESPACE := pancake

## [ Cluster ] Start a local Kind cluster
.PHONY: kind-up
kind-up:	
	kind create cluster --name $(CLUSTER_NAME) --config config/cluster/kind-config.yaml
	kubectl create namespace $(NAPESPACE)
	kubectl config set-context --current --namespace=$(NAPESPACE)


## [ Cluster ] Remove the local Kind cluster
.PHONY: kind-down
kind-down: 
	kind delete cluster --name $(CLUSTER_NAME)


.PHONY: kind-ctx
kind-ctx:
	kubectl config set-context --current --namespace=$(NAPESPACE)


.PHONY: kind-load
kind-load:
	kind load docker-image ${BUILD_NAME}:${BUILD_VERSION} --name=$(CLUSTER_NAME)


.PHONY: kind-deploy ## [ Cluster ] (Re)deploy pods to the local Kind cluster
kind-deploy: pod-delete docker-build kind-load pod-deploy


#=======
.PHONY: minikube-up
minikube-up:	## [ Cluster ] Start a local Minikube cluster
	minikube start
#	kubectl create namespace $(NAPESPACE)
	kubectl apply -f config/k8s/pancake-ns.yaml
	kubectl config set-context --current --namespace=$(NAPESPACE)

.PHONY: minikube-shutdown
minikube-shutdown: ## [ Cluster ] Remove the local Minikube cluster
	minikube delete

.PHONY: minikube-deploy	
minikube-deploy: pod-delete minikube-docker docker-build pod-deploy ## [ Cluster ] Deploy pods to the local Minikube cluster

.PHONY: minikube-tunnel 
minikube-tunnel: ## [ Cluster ] Starts Minikube tunnel for LoadBalancer support
	minikube tunnel

#=======| Kubernetes |=======
.PHONY: pod-deploy
pod-deploy:	## [ k8s ] Deploy a pod with the image
#	kubectl run ${BUILD_NAME} --image=${BUILD_NAME}:${BUILD_VERSION} --image-pull-policy=Never -n=$(NAPESPACE)
	kubectl apply -f config/k8s/pancake-ns.yaml
	kubectl apply -f config/k8s/pancake-storage.yaml
	kubectl apply -f config/k8s/pancake-app.yaml


.PHONY: pod-delete
pod-delete: 	## [ k8s ] Delete the deployed pod
	kubectl delete -f config/k8s/pancake-app.yaml --ignore-not-found --grace-period=3

.PHONY: pod-info
pod-info: 		## [ k8s ] Get kubectl info the cluster elements
	kubectl get po,svc,pvc,pv -o wide -n $(NAPESPACE)

.PHONY: pod-log
pod-log: 		## [ k8s ] Get logs of the pod streamed
	kubectl logs -f animated-pancake


#=======| Load testing |=======
.PHONY: load-test-get
load-test-get:	## [ Load testing ] Send 1000 request with 100 concurrent workers 
	hey -m GET -c 100 -n 1000 http://localhost:6543/v1/

.PHONY: load-test-post
load-test-post:	## [ Load testing ] Send 1000 request with 100 concurrent workers 
	hey -m POST -c 100 -n 1000  -H "Content-Type: application/json" -d '{ "subnetName": "testSubnet", "ipRange": "10.2.1.0/24" }' http://localhost:6543/v1/network
