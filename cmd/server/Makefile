# Default variables
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= server

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build
	cp ../../configs/server-template-dev.yml ../../configs/server-dev.yml

release:
	go build -o build/server main.go

clean:
	rm -rf build

run: release
	./build/server
