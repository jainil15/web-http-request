SHELL=/bin/bash
tailwind:
	@npx tailwindcss -i ./cmd/static/css/input.css -o ./cmd/static/css/output.css 

build: tailwind
	@templ generate
	@go build -o tmp/bin/ cmd/main.go


run: build
	@./tmp/bin/main

watch:
	@air -c .air.toml
