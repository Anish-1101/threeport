#help: @ List available make targets
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-19s\033[0m - %s\n", $$1, $$2}'

#build-codegen: @ Build codegen binary
build-codegen:
	go build -o bin/threeport-codegen cmd/codegen/main.go

#build-rest-api: @ Build REST API binary
build-rest-api:
	@export GOFLAGS=-mod=mod; export CGO_ENABLED=0; export GOOS=linux; export GOARCH=amd64; go build -a -o bin/threeport-rest-api cmd/rest-api/main.go

#generate: @ Run code generation
generate: build-codegen
	go generate ./...

## dev environment targets

#dev-down: @ Delete the local development environment
dev-down:
	kind delete cluster --name threeport-dev

#dev-cluster: @ Start a kind cluster for local development
dev-cluster:
	./dev/kind-config.sh
	kind create cluster --config dev/kind-config.yaml
	./dev/kind-install-metallb.sh

#dev-image-build: @ Build a development docker image for API that supports hot reloading
#dev-image-build: build-rest-api
dev-image-build:
	docker image build -t threeport-rest-api-dev:latest -f cmd/rest-api/image/Dockerfile-dev .

#dev-image-load: @ Build and load a development API image onto the kind cluster
dev-image-load: dev-image-build
	kind load docker-image threeport-rest-api-dev:latest -n threeport-dev

#dev-deps: @ Start a local kube cluster, build a new API image, deploy the message broker and API database
dev-deps: dev-cluster dev-image-load
	#kubectl apply -f dev/dependencies-cr.yaml -f dev/crdb.yaml -f dev/db-create.yaml -f dev/db-load.yaml
	kubectl apply -f dev/dependencies-cr.yaml -f dev/crdb.yaml -f dev/db-create.yaml

#dev-up: @ Run a local development environment
dev-up: dev-down dev-deps
	./dev/cr-wait.sh
	sleep 10  # allow time for all dependency containers to start
	kubectl apply -f dev/api-cr.yaml

#dev-forward-api: @ Forward local port 1323 to the local dev API
dev-forward-api:
	kubectl port-forward -n threeport-control-plane service/threeport-api-server 1323:80

