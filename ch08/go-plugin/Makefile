.DEFAULT_GOAL := help

help:
	# Commands:
	# make build  - Build the shared object (.so) library files
	# make duck   - Run the duck example
	# make frog   - Run the from example

build:
	@go build -buildmode=plugin -o duck/duck.so duck/duck.go
	@go build -buildmode=plugin -o frog/frog.so frog/frog.go
	@go build -buildmode=plugin -o fox/fox.so fox/fox.go

frog: build
	@go run main/main.go frog

duck: build
	@go run main/main.go duck
