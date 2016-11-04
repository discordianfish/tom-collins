REVISION = $(shell git rev-parse --short HEAD)
VERSION  = $(shell git name-rev --tags --name-only $(REVISION))
BRANCH   = $(shell git symbolic-ref --short -q HEAD)
DATE     = $(shell date +%Y%m%d-%H:%M:%S)

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

all: bin/tom

bin/tom: $(SOURCES)
	mkdir -p bin
	cd tom && go build -o ../$@ -ldflags "\
		-X main.BuildVersion=${VERSION:undefined=v0.0.0-dev} \
		-X main.BuildRevision=${REVISION} \
		-X main.BuildBranch=${BRANCH} \
		-X main.BuildUser=${USER} \
		-X main.BuildDate=${DATE} \
	" ./...

.PHONY: build
