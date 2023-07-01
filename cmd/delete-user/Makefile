# Default variables
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= delete-user

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build

release: setup
	go build -o build/delete-user main.go

clean:
	rm -rf build

run: install
	./build/delete-user
