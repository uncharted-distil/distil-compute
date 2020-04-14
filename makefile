VERSION=`git describe --tags`
TIMESTAMP=`date +%FT%T%z`

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  test          - test the source code"
	@echo "  lint          - lint the source code"
	@echo "  fmt           - format the source code"
	@echo "  install       - install dependencies"

lint:
	@go vet ./...
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

fmt:
	@go fmt ./...

build_static:
	env CGO_ENABLED=0 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS}

compile: lint
	@go build ./...

test: compile
	@go test ./...

proto:
	@protoc -I /usr/local/include \
		-I ./pipeline \
		--go_out=plugins=grpc:pipeline \
		./pipeline/*.proto

peg:
	@peg -inline ./primitive/compute/result/complex_field.peg

install:
	@go get -u github.com/golang/protobuf/protoc-gen-go@d3c38a4eb4970272b87a425ae00ccc4548e2f9bb
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/pointlander/peg@21bead84a59870739b2ee9eac3125ff9e5767e00
