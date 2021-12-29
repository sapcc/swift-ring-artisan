################################################################################
# This file is AUTOGENERATED with <https://github.com/sapcc/go-makefile-maker> #
# Edit Makefile.maker.yaml instead.                                            #
################################################################################

MAKEFLAGS=--warn-undefined-variables

# /bin/sh is dash on Debian which does not support all features of ash/bash
# to fix that we use /bin/bash only on Debian to not break Alpine
ifneq (,$(wildcard /etc/os-release)) # check file existence
	ifneq ($(shell grep -c debian /etc/os-release),0)
		SHELL := /bin/bash
	endif
endif

default: build-all

generate: build/swift-ring-artisan
	./build/swift-ring-artisan parse testing/builder-output.txt -o testing/builder-output.yaml

build-all: build/swift-ring-artisan

GO_BUILDFLAGS = -mod vendor
GO_LDFLAGS =
GO_TESTENV =

build/swift-ring-artisan: FORCE
	go build $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -o build/swift-ring-artisan .

DESTDIR =
ifeq ($(shell uname -s),Darwin)
	PREFIX = /usr/local
else
	PREFIX = /usr
endif

install: FORCE build/swift-ring-artisan
	install -D -m 0755 build/swift-ring-artisan "$(DESTDIR)$(PREFIX)/bin/swift-ring-artisan"

# which packages to test with static checkers
GO_ALLPKGS := $(shell go list ./...)
# which files to test with static checkers (this contains a list of globs)
GO_ALLFILES := $(addsuffix /*.go,$(patsubst $(shell go list .)%,.%,$(shell go list ./...)))
# which packages to test with "go test"
GO_TESTPKGS := $(shell go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)
# which packages to measure coverage for
GO_COVERPKGS := $(shell go list ./... | command grep -E '/internal|/pkg')
# to get around weird Makefile syntax restrictions, we need variables containing a space and comma
space := $(null) $(null)
comma := ,

check: build-all static-check build/cover.html FORCE
	@printf "\e[1;32m>> All checks successful.\e[0m\n"

prepare-static-check: FORCE
	@if ! hash staticcheck 2>/dev/null; then printf "\e[1;36m>> Installing staticcheck...\e[0m\n"; go install honnef.co/go/tools/cmd/staticcheck@latest; fi
	@if ! hash exportloopref 2>/dev/null; then printf "\e[1;36m>> Installing exportloopref...\e[0m\n"; go install github.com/kyoh86/exportloopref/cmd/exportloopref@latest; fi
	@if ! hash rowserrcheck 2>/dev/null; then printf "\e[1;36m>> Installing rowserrcheck...\e[0m\n"; go install github.com/jingyugao/rowserrcheck@latest; fi
	@if ! hash unconvert 2>/dev/null; then printf "\e[1;36m>> Installing unconvert...\e[0m\n"; go install github.com/mdempsky/unconvert@latest; fi
	@if ! hash unparam 2>/dev/null; then printf "\e[1;36m>> Installing unparam...\e[0m\n"; go install mvdan.cc/unparam@latest; fi

static-check: prepare-static-check FORCE
	@printf "\e[1;36m>> gofmt\e[0m\n"
	@if s="$$(gofmt -s -d $(GO_ALLFILES) 2>/dev/null)" && test -n "$$s"; then echo "$$s"; false; fi
	@printf "\e[1;36m>> staticcheck\e[0m\n"
	@staticcheck -checks 'inherit,-ST1015' $(GO_ALLPKGS)
	@printf "\e[1;36m>> exportloopref\e[0m\n"
	@exportloopref $(GO_ALLPKGS)
	@printf "\e[1;36m>> rowserrcheck\e[0m\n"
	@rowserrcheck $(GO_ALLPKGS)
	@printf "\e[1;36m>> unconvert\e[0m\n"
	@unconvert $(GO_ALLPKGS)
	@printf "\e[1;36m>> unparam\e[0m\n"
	@unparam $(GO_ALLPKGS)
	@printf "\e[1;36m>> go vet\e[0m\n"
	@go vet $(GO_BUILDFLAGS) $(GO_ALLPKGS)

build/cover.out: build FORCE
	@printf "\e[1;36m>> go test\e[0m\n"
	@env $(GO_TESTENV) go test $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -p 1 -coverprofile=$@ -covermode=count -coverpkg=$(subst $(space),$(comma),$(GO_COVERPKGS)) $(GO_TESTPKGS)

build/cover.html: build/cover.out
	@printf "\e[1;36m>> go tool cover > build/cover.html\e[0m\n"
	@go tool cover -html $< -o $@

build:
	@mkdir $@

vendor: FORCE
	go mod tidy
	go mod vendor
	go mod verify

vendor-compat: FORCE
	go mod tidy -compat=$(shell awk '$$1 == "go" { print $$2 }' < go.mod)
	go mod vendor
	go mod verify

license-headers: FORCE
	@if ! hash addlicense 2>/dev/null; then printf "\e[1;36m>> Installing addlicense...\e[0m\n"; go install github.com/google/addlicense@latest; fi
	find * \( -name vendor -type d -prune \) -o \( -name \*.go -exec addlicense -c "SAP SE" -- {} + \)

clean: FORCE
	git clean -dxf build

.PHONY: FORCE
