include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d algotrade-postgres

env-down:
	@docker compose down algotrade-postgres

env-cleanup:
	@read -p "clean all environment volumes[y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down port-forwarder algotrade-postgres && \
	  rm -rf out/pgdata && \
	  echo "environment files was cleaned up"; \
	else \
	  echo "cleaning was canceled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "not found parameter: seq. Example: make migrate-create seq=init"; \
  		exit 1; \
  	fi; \
	docker compose run --rm algotrade-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)";

migrate-init:
	@make migrate-create seq=init

migrate-action:
	@if [ -z "$(action)" ]; then \
  		echo "not found action. Example: make migrate-action action="; \
  		exit 1; \
  	fi; \
	docker compose run --rm algotrade-postgres-migrate \
    		-path /migrations \
    		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@algotrade-postgres:5432/${POSTGRES_DB}?sslomode=disable \
    		"$(action)";

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action acton=down

algotrade-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ./backend/cmd/main.go
