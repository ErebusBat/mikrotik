
## Defines
SRC     = cmd/bandwidth/*.go
GOPKG   = github.com/ErebusBat/mikrotik/cmd/bandwidth
EXEFILE = bandwidth

## Input and Output Variables
GIT_VER    := $(shell git describe --always --dirty)
DATE_STAMP := $(shell date +%Y%m%d-%H%M%S)

# Use equals (as opposed to colon-equals) so that they are re-evaluated
# for each target, each time
BUILDTAG   = $(DATE_STAMP)-$(GIT_VER)
PLATARCH   = $(GOOS)-$(GOARCH)
OUTDIR     = bin/$(BUILDTAG)/$(PLATARCH)
GOLDFLAGS  = -ldflags '-X main.BUILD_TAG $(BUILDTAG)'
GOMAKE     = GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o "$(OUTDIR)/$(EXEFILE)$(EXEEXT)" $(GOLDFLAGS) $(SRC)
ZIPFILE    = $(OUTDIR)/../$(PLATARCH)-$(BUILDTAG).zip
ZIPCMD     = zip -j "$(ZIPFILE)" "$(OUTDIR)/$(EXEFILE)$(EXEEXT)"

################################################################################
# Generic Targets
################################################################################

default: build

# .PHONY: default mac
# .PHONY: build

build: mac

clean:
	rm -dfr bin/*
	go clean -i -x

uninstall:
	go clean -i -x $(GOPKG)

all: mac win linux

deps:
	go get -u -f $(GOPKG)

install: deps goinst
goinst:
	go install $(GOLDFLAGS) $(GOPKG)

################################################################################
# Mac OSX Targets
################################################################################
mac: mac64
mac64: GOOS=darwin
mac64: GOARCH=amd64
mac64:
	mkdir -p $(OUTDIR)
	$(GOMAKE)
	$(ZIPCMD)

################################################################################
# Windows Targets
################################################################################
win: win64 win32
win64: GOOS=windows
win64: GOARCH=amd64
win64: EXEEXT=.exe
win64:
	mkdir -p $(OUTDIR)
	$(GOMAKE)
	$(ZIPCMD)

win32: GOOS=windows
win32: GOARCH=386
win32: EXEEXT=.exe
win32:
	mkdir -p $(OUTDIR)
	$(GOMAKE)
	$(ZIPCMD)


################################################################################
# Linux Targets
################################################################################
linux: linux64 linux32
linux64: GOOS=linux
linux64: GOARCH=amd64
linux64:
	mkdir -p $(OUTDIR)
	$(GOMAKE)
	$(ZIPCMD)

linux32: GOOS=linux
linux32: GOARCH=386
linux32:
	mkdir -p $(OUTDIR)
	$(GOMAKE)
	$(ZIPCMD)

