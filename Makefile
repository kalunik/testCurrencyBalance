CONFIG_LOCAL = config/config_local.yml
COMPOSE_LOCAL = ./docker-compose.local.yml


local:      config
			@(docker compose -f $(COMPOSE_LOCAL) up -d --build)
			@(sleep 15)
			@(go run cmd/main.go)
run:
			@(docker compose -f $(COMPOSE_LOCAL) up -d --build)

config:
			@(echo "Creating configs for launch. Don't forget change sample credentials.")
			@(cp ./config/config_sample.yml $(CONFIG_LOCAL))
			@(cp .env_sample .env)

migrations:
			@migrate create -ext sql -dir database/migrate/ -seq init

.PHONY: local run config migrations