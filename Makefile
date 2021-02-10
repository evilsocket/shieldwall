all: agent api

agent: _build
	@go build -ldflags="-w -s" -o _build/shieldwall-agent cmd/agent/*.go

bindata:
	@go get -u github.com/go-bindata/go-bindata/...

frontend: bindata
	@rm -rf frontend/compiled.go
	@cd frontend && npm run build
	@go-bindata -o frontend/compiled.go -pkg frontend -prefix frontend/dist ./frontend/dist/...

api_and_frontend: _build frontend
	@go build -ldflags="-w -s" -o _build/shieldwall-api cmd/api/*.go

api: _build
	@go build -ldflags="-w -s" -o _build/shieldwall-api cmd/api/*.go

test:
	@go test -short ./...

_build:
	@mkdir -p _build

clean:
	@rm -rf _build

composer_build:
	@docker-compose build

composer_up: composer_build
	@docker-compose up

install_api:
	@service shieldwall-api stop || true
	@cp _build/shieldwall-api /usr/bin/
	@setcap 'cap_net_bind_service=+ep' /usr/bin/shieldwall-api
	@mkdir -p /etc/shieldwall/
	@test -s /etc/shieldwall/config.yaml || cp api.example.yaml /etc/shieldwall/config.yaml
	@cp shieldwall-api.service /etc/systemd/system/
	@systemctl daemon-reload
	@systemctl enable shieldwall-api
	@service shieldwall-api restart

install_agent:
	@service shieldwall-agent stop || true
	@cp _build/shieldwall-agent /usr/bin/
	@mkdir -p /etc/shieldwall/
	@test -s /etc/shieldwall/config.yaml || cp agent.example.yaml /etc/shieldwall/config.yaml
	@cp shieldwall-agent.service /etc/systemd/system/
	@systemctl daemon-reload
	@systemctl enable shieldwall-agent
	@service shieldwall-agent restart

self_update:
	clear && git checkout . && git pull && make clean && make api && sudo make install_api