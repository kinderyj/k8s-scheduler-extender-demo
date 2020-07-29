BIN_DIR=_output/bin


.EXPORT_ALL_VARIABLES:

all: local

init:
	mkdir -p ${BIN_DIR}

local: init
	go build -o=${BIN_DIR}/scheduler-extender-demo ./cmd/scheduler.go

build-linux: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=${BIN_DIR}/scheduler-extender-demo ./cmd/scheduler.go

image: build-linux
	docker build --no-cache . -t scheduler-extender-demo

update:
	go mod download
	go mod tidy
	go mod vendor

clean:
	rm -rf _output/
	rm -f *.log