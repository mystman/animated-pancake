#=======| Defaults |=======
.DEFAULT_GOAL := build

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
.PHONY: help
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	
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
debug: lint vet cover # [ Run ] Compile and run with gcflags
	go run -gcflags '-m -l' cmd/main.go

#=======| Build |=======
.PHONY: build 
build: lint vet cover ## [ Build ] Build the project binary
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

.PHONY: cover
cover: ## [ Test ] Test coverage
	go test -cover ./...

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


#=======| Cluster |=======
CLUSTER_NAME := pancake-cluster
NAPESPACE := pancake

.PHONY: kind-up
kind-up:	## [ Cluster ] Start a local Kind cluster
	kind create cluster --name $(CLUSTER_NAME)
	kubectl create namespace pancake
	kubectl config set-context --current --namespace=$(NAPESPACE)

.PHONY: kind-down
kind-down: ## [ Cluster ] Remove the local Kind cluster
	kind delete cluster --name $(CLUSTER_NAME)

.PHONY: kube-context
kind-ctx:
	kubectl config set-context --current --namespace=$(NAPESPACE)

kind-load:
	kind load docker-image ${BUILD_NAME}:${BUILD_VERSION} --name=$(CLUSTER_NAME)


.PHONY: pod-deploy
pod-deploy: kind-load	 ## [ k8s ] Deploy a pod with the image
	kubectl run ${BUILD_NAME} --image=${BUILD_NAME}:${BUILD_VERSION} --image-pull-policy=Never -n=$(NAPESPACE)
#	 -o=yaml --dry-run=client

.PHONY: pod-delete
pod-delete: 	## [ k8s ] Delete the deployed pod
	kubectl delete po ${BUILD_NAME}
