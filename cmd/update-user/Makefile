# Default variables
INSTALL_FOLDER ?= ../../bin
APPLICATION ?= update-user

BRANCH ?= master
TAG ?= ${BRANCH}

install: release
	cp -r build/* ${INSTALL_FOLDER}

setup:
	mkdir -p build

release: setup
	go build -o build/update-user main.go

clean:
	rm -rf build

run: install
	./build/update-user
