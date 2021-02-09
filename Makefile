all: agent api

agent: _build
	@go build -ldflags="-w -s" -o _build/shieldwall-agent cmd/agent/*.go

api: _build
	@go build -ldflags="-w -s" -o _build/shieldwall-api cmd/api/*.go

test:
	@go test -short ./...

_build:
	@mkdir -p _build

clean:
	@rm -rf _build

composer_build:
	docker-compose build

composer_up: composer_build
	docker-compose up