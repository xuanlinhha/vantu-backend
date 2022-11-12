ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output
DOCKER_HUB_ID=xuanlinhha

BACKEND_IMAGE := go-vantu-backend
.PHONY: vantu-backend
vantu-backend:
	docker build -t $(BACKEND_IMAGE) --no-cache=true . && \
	docker tag $(BACKEND_IMAGE) $(DOCKER_HUB_ID)/$(BACKEND_IMAGE)
