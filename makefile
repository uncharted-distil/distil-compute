version=0.1.0

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  install       - install dependencies"

install:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
