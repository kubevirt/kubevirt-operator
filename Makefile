SHELL=/bin/bash -o pipefail

REPO?=docker.io/kubevirt/kubevirt-operator
TAG?=$(shell git rev-parse --short HEAD)

GOLANG_FILES:=$(shell find . -name \*.go -print)
pkgs = $(shell go list ./... | grep -v /vendor/ )

dep:
	dep ensure -v
	wget -o kubevirt.yaml https://github.com/kubevirt/kubevirt/releases/download/v0.6.4/kubevirt.yaml

all: format dep compile build deploy

clean:
	# Remove all files and directories ignored by git.
	git clean -Xfd .

############
# Compile  #
############

compile: kubevirt-operator

kubevirt-operator: $(GOLANG_FILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
	-o $@ cmd/kubevirt-operator/main.go

##############
# Build      #
##############

build:
	$(GOPATH)/bin/operator-sdk build $(REPO):$(TAG)

push:
	docker push $(REPO):$(TAG)

##############
# OLM        #
##############

olm:
	oc create -f kubevirt-operator-configmap.yaml -f kubevirt-catalog-source.yaml
	# Temporary Hack: Need to fix rbac
	oc adm policy add-cluster-role-to-user -z kubevirt-operator cluster-admin

rm-olm:
	oc delete -f kubevirt-operator-configmap.yaml -f kubevirt-catalog-source.yaml

##############
# Deploy     #
##############

setprivileges:
	# Give the kubevirt operator user 'default' cluster-admin privileges
	oc adm policy add-cluster-role-to-user cluster-admin -z default

deploy: setprivileges
	# Deploy the app-operator
	kubectl create -f deploy/rbac.yaml
	kubectl create -f deploy/crd.yaml
	kubectl create -f deploy/operator.yaml
	# Create the operator resource (kubectl get apps)
	kubectl create -f deploy/cr.yaml

##############
# Undeploy   #
##############

undeploy:
	# Undeploy the app-operator
	kubectl delete -f deploy/crd.yaml --ignore-not-found
	kubectl delete -f deploy/operator.yaml --ignore-not-found
	kubectl delete -f deploy/rbac.yaml --ignore-not-found

##############
# Formatting #
##############

format: go-fmt

go-fmt:
	go fmt $(pkgs)

.PHONY: dep all clean compile build push deploy undeploy format olm rm-olm
