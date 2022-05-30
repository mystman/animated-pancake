#=======| Defaults |=======
.DEFAULT_GOAL := run

#=======| Project settings |=======
BINARY_NAME=go-kub-stub
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
run: lint vet test  ## [ Run ] Run the main.go (with LDFLAGS)
	go run cmd/main.go ${LDFLAGS}
	
.PHONY: debug
debug: lint vet cover ## [ Run ] Compile and run with gcflags
	go run -gcflags '-m -l' cmd/main.go

.PHONY: start
start: build ## [ Run ] Build and execute binary
	./${BUILD_DIR}/${BINARY_NAME}

#=======| Build |=======
.PHONY: build 
build: lint vet cover ## [ Build ] Build the project binary
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/main.go

.PHONY: build_release
build_release: ## [ Build ] To be used from Dockerfile for building a release
	go build ${LDFLAGS} -o ${BUILD_DIR}/service cmd/main.go



.PHONY: clean
clean: ## [ Build ] Clean the project binary and trace files
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

.PHONY: external_test
external_test: ## [ Test ] Executes binary form shell with params and stdin
	echo "1 2 3 4 5" | ${BUILD_DIR}/${BINARY_NAME} param1 param2 param3 


#=======| Debug - Trace & Benchmark |=======
.PHONY: trace
trace: ## [ Debug ] Run tests with trace
	go test ./... -trace=copy_trace.out

.PHONY: trace_show 
trace_show: trace ## [ Debug ] Takes copy_trace.out file and displays it in Trace viewer
	go tool trace copy_trace.out
	
.PHONY: bench
bench: ## [ Debug ] Run benchmark
	go test -v -run="none" -bench="BenchmarkMain" -benchmem	

.PHONY: bench_long 
bench_long: ## [ Debug ] Run long benchmark
	go test -v -run="none" -bench="BenchmarkMain" -benchmem -benchtime=3s -count 3 ./cmd

.PHONY: fuzz
fuzz: ## [ Debug ] Run fuzzing
	go test -fuzz=. -fuzztime=10s ./...



#=======| Containers |=======
.PHONY: podman-build
podman-build: minikube-use-podman ## Build container image with Podman
	podman build -t ${BUILD_NAME}:${BUILD_VERSION} .

.PHONY: docker-build
docker-build: minikube-use-docker ## Build container image with Podman
	docker build -t ${BUILD_NAME}:${BUILD_VERSION} .


#=======| Cluster |=======
.PHONY: minikube-use-podman
minikube-use-podman:
	eval $(minikube -p minikube podman-env)	

.PHONY: minikube-use-docker
minikube-use-docker:
	eval $(minikube -p minikube docker-env)	


.PHONY: minikube-podman-start
minikube-podman-start: ## Starts local Minikube cluster with Podman
	podman machine start
	minikube start --driver=podman 

.PHONY: minikube-start
minikube-start: ## Starts a default local Minikube cluster
	minikube start


.PHONY: pod-deploy
pod-deploy: ## Deploy a pod with the image (to default namespace)
	kubectl run ${BUILD_NAME} --image=${BUILD_NAME}:${BUILD_VERSION} --image-pull-policy=Never

.PHONY: pod-delete
pod-delete: ## Delete the pod
	kubectl delete po ${BUILD_NAME}
