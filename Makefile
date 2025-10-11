target=svgocode
builddir=build

BUILDLDFLAGS?=-ldflags="-s -w"
DEBUGLDFLAGS?=

all: doc build

install:
	go install -ldflags="$(BUILDLDFLAGS)" $(target)

build:
	-rm $(builddir)/*
	# Remove debug information from the build
	go build $(BUILDLDFLAGS) -o $(builddir)/$(target) main.go

build-dev:
	-rm $(builddir)/*
	# Keep debug information in the build
	go build $(DEBUGLDFLAGS) -o $(builddir)/$(target) main.go

run: build
	# Run a clean build
	./$(builddir)/$(target)

dev: build-dev
	# Run with debug info
	DEBUG=TRUE ./$(builddir)/$(target)

gdb: build-dev
	# Run via GDB
	DEBUG=TRUE gdb $(builddir)/$(target)

clean:
	# Remove temporary files
	-rm tags
	-rm $(builddir)/*

.PHONY:build build-dev run install clean gdb dev
