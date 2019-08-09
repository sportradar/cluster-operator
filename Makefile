# Image URL to use all building/pushing image targets
CONTROLLER_IMAGE=eu.gcr.io/cf-rabbitmq-for-k8s-bunny/rabbitmq-for-kubernetes-controller
CI_IMAGE=eu.gcr.io/cf-rabbitmq-for-k8s-bunny/rabbitmq-for-kubernetes-ci
GCP_PROJECT=cf-rabbitmq-for-k8s-bunny
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Run unit tests
unit-tests: generate fmt vet manifests
	ginkgo -r internal/

# Run integration tests
integration-tests: generate fmt vet manifests
	ginkgo -r api/ controllers/

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./api/...;./controllers/..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Deploy manager
deploy-manager:
	kubectl apply -k config/default/base

# Deploy manager in CI
deploy-manager-ci:
	kubectl apply -k config/default/overlays/ci

# Deploy local rabbitmqcluster
deploy-sample:
	kubectl apply -k config/samples/base

configure-kubectl-ci: ci-cluster
	gcloud auth activate-service-account --key-file=$(KUBECTL_SECRET_TOKEN_PATH)
	gcloud container clusters get-credentials $(CI_CLUSTER) --region europe-west1 --project $(GCP_PROJECT)

# Cleanup all controller artefacts
destroy:
	kubectl delete -k config/default/base
	kubectl delete -k config/namespace/base

destroy-ci: configure-kubectl-ci
	kubectl delete -k config/default/overlays/ci --ignore-not-found=true
	kubectl delete -k config/namespace/overlays/ci --ignore-not-found=true

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate manifests fmt vet install deploy-namespace
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

deploy-namespace:
	kubectl apply -k config/namespace/base

deploy-namespace-ci:
	kubectl apply -k config/namespace/overlays/ci

deploy-master: install deploy-namespace gcr-viewer
	kubectl apply -k config/default/base

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests deploy-namespace gcr-viewer deploy-manager

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy-ci: configure-kubectl-ci patch-controller-image manifests deploy-namespace-ci gcr-viewer-ci deploy-manager-ci

# Build the docker image
docker-build:
	docker build . -t $(CONTROLLER_IMAGE):latest

docker-build-ci-image:
	docker build ci/ -t ${CI_IMAGE}
	docker push ${CI_IMAGE}

# Push the docker image
docker-push:
	docker push $(CONTROLLER_IMAGE):latest

docker-image-release: controller-image-tag
	docker build . -t $(CONTROLLER_IMAGE):$(CONTROLLER_IMAGE_TAG)
	docker push $(CONTROLLER_IMAGE):$(CONTROLLER_IMAGE_TAG)

system-tests:
	NAMESPACE="pivotal-rabbitmq-system" ginkgo -p --randomizeAllSpecs -r system_tests/

system-tests-ci:
	NAMESPACE="pivotal-rabbitmq-system-ci" ginkgo -p --randomizeAllSpecs -r system_tests/

GCR_VIEWER_ACCOUNT_EMAIL=gcr-viewer@cf-rabbitmq-for-k8s-bunny.iam.gserviceaccount.com
GCR_VIEWER_ACCOUNT_NAME=gcr-viewer
GCR_VIEWER_KEY=$(shell lpassd show "Shared-RabbitMQ for Kubernetes/ci-gcr-pull" --notes | jq -c)
gcr-viewer: operator-namespace
	echo "creating gcr-viewer secret and patching default service account"
	@kubectl -n $(K8S_OPERATOR_NAMESPACE) create secret docker-registry $(GCR_VIEWER_ACCOUNT_NAME) --docker-server=https://eu.gcr.io --docker-username=_json_key --docker-email=$(GCR_VIEWER_ACCOUNT_EMAIL) --docker-password='$(GCR_VIEWER_KEY)' || true
	@kubectl -n $(K8S_OPERATOR_NAMESPACE) patch serviceaccount default -p '{"imagePullSecrets": [{"name": "$(GCR_VIEWER_ACCOUNT_NAME)"}]}'

gcr-viewer-ci: operator-namespace
	echo "creating gcr-viewer secret and patching default service account"
	@kubectl -n $(K8S_OPERATOR_NAMESPACE) create secret docker-registry $(GCR_VIEWER_ACCOUNT_NAME) --docker-server=https://eu.gcr.io --docker-username=_json_key --docker-email=$(GCR_VIEWER_ACCOUNT_EMAIL) --docker-password='$(GCR_VIEWER_KEY_CI)' || true
	@kubectl -n $(K8S_OPERATOR_NAMESPACE) patch serviceaccount default -p '{"imagePullSecrets": [{"name": "$(GCR_VIEWER_ACCOUNT_NAME)"}]}'

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.1
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

patch-controller-image: controller-image-tag
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"$(CONTROLLER_IMAGE):$(CONTROLLER_IMAGE_TAG)"'@' ./config/default/base/manager_image_patch.yaml

operator-namespace:
ifeq (, $(K8S_OPERATOR_NAMESPACE))
K8S_OPERATOR_NAMESPACE=pivotal-rabbitmq-system
endif

ci-cluster:
ifeq (, $(CI_CLUSTER))
CI_CLUSTER=ci-bunny
endif

controller-image-tag:
ifeq (, $(CONTROLLER_IMAGE_TAG))
CONTROLLER_IMAGE_TAG=0.1.0
endif
