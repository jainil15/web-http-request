SHELL=/bin/bash
build:
	@templ generate
	@go build -o tmp/bin/ cmd/main.go

run: build
	@./tmp/bin/main

watch:
	@air -c .air.toml
