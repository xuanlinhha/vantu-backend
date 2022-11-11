ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output
DOCKER_HUB_ID=xuanlinhha

BACKEND_IMAGE := go-vantu-backend
.PHONY: vantu-backend
vantu-backend:
	mkdir -p $(OUTPUT_DIR) && \
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR) $(ROOTDIR)/cmd/vantu-backend && \
	cp $(ROOTDIR)/Dockerfile $(OUTPUT_DIR) && \
	cp $(ROOTDIR)/data/* $(OUTPUT_DIR) && \
	docker build $(OUTPUT_DIR) -t $(BACKEND_IMAGE) -f $(OUTPUT_DIR)/Dockerfile --no-cache=true && \
	docker tag $(BACKEND_IMAGE) $(DOCKER_HUB_ID)/$(BACKEND_IMAGE)
