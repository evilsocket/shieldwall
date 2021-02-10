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

install_api: api
	service shieldwall-api stop
	cp _build/shieldwall-api /usr/bin/
	setcap 'cap_net_bind_service=+ep' /usr/bin/shieldwall-api
	mkdir -p /etc/shieldwall/
	test -s /etc/shieldwall/config.yaml || echo cp api.example.yaml /etc/shieldwall/config.yaml
	cp shieldwall-api.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable shieldwall-api
	service shieldwall-api restart


