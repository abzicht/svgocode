target=svgocode
package=github.com/abzicht/svgocode
builddir=build

BUILDLDFLAGS?=-ldflags="-s -w"
DEBUGLDFLAGS?=
TESTFLAGS?=

all: doc build

install:
	go install $(BUILDLDFLAGS) $(package)

build:
	-rm $(builddir)/*
	# Remove debug information from the build
	go build $(BUILDLDFLAGS) -o $(builddir)/$(target) main.go

build-dev:
	-rm $(builddir)/*
	# Keep debug information in the build
	go build $(DEBUGLDFLAGS) -o $(builddir)/$(target) main.go

test:
	go test $(TESTFLAGS) ./...

run: build
	# Run a clean build
	./$(builddir)/$(target)

dev: build-dev
	# Run with debug info
	./$(builddir)/$(target) -v 5

gdb: build-dev
	# Run via GDB
	gdb $(builddir)/$(target)

clean:
	# Remove temporary files
	-rm tags
	-rm $(builddir)/*

.PHONY:build build-dev run install clean gdb dev
