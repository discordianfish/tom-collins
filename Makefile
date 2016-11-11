REVISION = $(shell git rev-parse --short HEAD)
VERSION  = $(shell git name-rev --tags --name-only $(REVISION))
BRANCH   = $(shell git symbolic-ref --short -q HEAD)
DATE     = $(shell date +%Y%m%d-%H:%M:%S)

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

all: bin/tom

bin/tom: $(SOURCES)
	mkdir -p bin
	cd tom && go build -o ../$@ -ldflags "\
		-X main.buildVersion=${VERSION:undefined=v0.0.0-dev} \
		-X main.buildRevision=${REVISION} \
		-X main.buildBranch=${BRANCH} \
		-X main.buildUser=${USER} \
		-X main.buildDate=${DATE} \
	" ./...

.PHONY: build
