# https://github.com/terraform-docs/terraform-docs/blob/master/Makefile
# Project variables
PROJECT_NAME  := pddoc
PROJECT_OWNER := rqtx
DESCRIPTION   := A tool to generate documentation from cloud providers in markdown formats
PROJECT_URL   := https://github.com/$(PROJECT_OWNER)/$(PROJECT_NAME)
LICENSE       := MIT

# Build variables
BUILD_DIR    := bin
COMMIT_HASH  ?= $(shell git rev-parse --short HEAD 2>/dev/null)
CUR_VERSION  ?= $(shell $(GORUN) main.go version -s)
COVERAGE_OUT := coverage.out

# Go variables
GO          ?= go
GO_PACKAGE  := github.com/$(PROJECT_OWNER)/$(PROJECT_NAME)
GOOS        ?= $(shell $(GO) env GOOS)
GOARCH      ?= $(shell $(GO) env GOARCH)

GOLDFLAGS   += -X $(GO_PACKAGE)/internal/version.commit=$(COMMIT_HASH)

GOBUILD     ?= CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags="$(GOLDFLAGS)"
GORUN       ?= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) run

# Docker variables
DEFAULT_TAG  ?= $(shell echo "$(CUR_VERSION)" | tr -d 'v')
DOCKER_IMAGE := $(PROJECT_OWNER)/$(PROJECT_NAME)
DOCKER_TAG   ?= $(DEFAULT_TAG)

# Binary versions
GOLANGCI_VERSION  := v1.38.0

.PHONY: all
all: clean verify checkfmt lint test build

#########
##@ Build

.PHONY: go-build
go-build: clean ## Build binary for current OS/ARCH
	$(GOBUILD) -o ./$(BUILD_DIR)/$(PROJECT_NAME)

.PHONY: docker-build
docker-build:   ## Build Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker build --pull --tag $(DOCKER_IMAGE):$(DOCKER_TAG) --file Dockerfile .

.PHONY: docker-push
docker-push:   ## Push Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

###############
##@ Development

.PHONY: clean
clean:   ## Clean workspace
	@ $(MAKE) --no-print-directory log-$@	
	rm -rf ./$(BUILD_DIR) ./$(COVERAGE_OUT)

########################################################################
## Self-Documenting Makefile Help                                     ##
## https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html ##
########################################################################

########
##@ Help

.PHONY: help
help:   ## Display this help
	@awk \
		-v "col=\033[36m" -v "nocol=\033[0m" \
		' \
			BEGIN { \
				FS = ":.*##" ; \
				printf "Usage:\n  make %s<target>%s\n", col, nocol \
			} \
			/^[a-zA-Z_-]+:.*?##/ { \
				printf "  %s%-12s%s %s\n", col, $$1, nocol, $$2 \
			} \
			/^##@/ { \
				printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol \
			} \
		' $(MAKEFILE_LIST)

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'