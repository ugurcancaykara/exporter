BINARY_NAME := exporter
DIST := $(shell pwd)/_dist
BINARY_PATH := /usr/local/bin/exporter


.PHONY: build
build:
	go build -o bin/${BINARY_NAME}



.PHONY: dist
dist:
	mkdir -p ${DIST}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}
	tar -zcvf $(DIST)/${BINARY_NAME}-linux.tgz ${BINARY_NAME} README.md
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}
	tar -zcvf $(DIST)/${BINARY_NAME}-linux-arm.tgz ${BINARY_NAME} README.md 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}
	tar -zcvf $(DIST)/${BINARY_NAME}-macos.tgz ${BINARY_NAME} README.md

.PHONY: install
install: dist
	if [ "$$(uname)" = "Darwin" ]; then file="${BINARY_NAME}-macos"; \
 	elif [ "$$(uname)" = "Linux" ]; then file="${BINARY_NAME}-linux"; \
	else echo "doesn't support this operation system"; \
	fi; \
	mkdir -p $(DIST)/$$file; \
	tar -xf $(DIST)/$$file.tgz -C $(DIST)/$$file; \
	sudo cp -r exporter ${BINARY_PATH} 