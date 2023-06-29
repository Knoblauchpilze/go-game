# Default variables
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= cli

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build

release: setup
	go build -o build/cli main.go

clean:
	rm -rf build
