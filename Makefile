# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

#Swagger ui
SWAGGERUIVERSION := 3.20.5

# Setup name variables for the package/tool
NAME := license-manager
DOCKERTAG := lic-man

# Set any default go build tags
BUILDTAGS :=

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := ${PREFIX}/build

GO_LDFLAGS_STATIC=-ldflags "-w -s -extldflags -static"


# Set our default go compiler
GO := go

# List the GOOS and GOARCH to build
GOOSARCHES = linux/amd64

.PHONY: build
build: $(NAME) ## Builds a dynamic executable or package

$(NAME): $(wildcard *.go) $(wildcard */*.go) VERSION.txt
	@echo "+ $@"
	$(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) .


.PHONY: docker
docker: static ## Builds a docker image
	@echo "+ $@"
	@docker build -t ${DOCKERTAG} --build-arg BINARY=$(NAME) .


.PHONY: static
static: ## Builds a static executable
	@echo "+ $@"
	 $(GO) build \
		-tags " netgo $(BUILDTAGS) static_build" \
		${GO_LDFLAGS_STATIC} -o $(NAME) .

all: clean build test   ## Runs a clean, build, fmt, lint, test, staticcheck, vet and install


.PHONY: test
test: ## Runs the go tests
	@echo "+ $@"
	@$(GO) test $(RACEFLAG) -v -tags "$(BUILDTAGS) cgo" $(shell $(GO) list ./... | grep -v vendor)


.PHONY: install
install: ## Installs the executable or package
	@echo "+ $@"
	$(GO) install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} .


define buildpretty
mkdir -p $(BUILDDIR)/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) $(GO) build \
	 -o $(BUILDDIR)/$(1)/$(2)/$(NAME) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).md5;
sha256sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).sha256;
endef

.PHONY: cross
cross: *.go VERSION.txt ## Builds the cross-compiled binaries, creating a clean directory structure (eg. GOOS/GOARCH/binary)
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildpretty,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

define buildrelease
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 $(GO) build \
	 -o $(BUILDDIR)/$(NAME)-$(1)-$(2) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} .;
md5sum $(BUILDDIR)/$(NAME)-$(1)-$(2) > $(BUILDDIR)/$(NAME)-$(1)-$(2).md5;
sha256sum $(BUILDDIR)/$(NAME)-$(1)-$(2) > $(BUILDDIR)/$(NAME)-$(1)-$(2).sha256;
endef

.PHONY: clean
clean: ## Cleanup any build binaries or packages
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: swaggerui-download
swaggerui-download: ## Downloads swaggerui release
	@curl -L https://github.com/swagger-api/swagger-ui/archive/v$(SWAGGERUIVERSION).zip --output /tmp/v$(SWAGGERUIVERSION).zip
	@unzip -o -d /tmp/ /tmp/v$(SWAGGERUIVERSION).zip
	@mkdir -p swaggerui

.PHONY: swaggerui
swaggerui: swaggerui-download ## Generates static files for swaggerui
	@GO111MODULE=off $(GO) get github.com/rakyll/statik
	cp openapi.yml /tmp/swagger-ui-$(SWAGGERUIVERSION)/dist/openapi.yml
	# Change the swaggerui config file to point to our openapi.yml definition file and also set the redirect url
	sed -i -e 's|url: "https://petstore.swagger.io/v2/swagger.json",|url: "./openapi.yml",\n\toauth2RedirectUrl: window.location.origin + window.location.pathname + "oauth2-redirect.html",|g' /tmp/swagger-ui-3.20.5/dist/index.html
	statik -f -c="Static assets for swaggerui v$(SWAGGERUIVERSION)" -src=/tmp/swagger-ui-$(SWAGGERUIVERSION)/dist -dest=. -p=swaggerui
