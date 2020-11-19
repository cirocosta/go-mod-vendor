VERSION ?= "0.0.4"
REPOSITORY ?= "harbor-repo.vmware.com/tanzu_delivery_pipeline/go-mod-vendor"

build:
	./scripts/package.sh --version $(VERSION)
.PHONY: build

push:
	sudo skopeo copy \
		oci-archive:build/buildpackage.cnb \
		docker://$(REPOSITORY):$(VERSION)
