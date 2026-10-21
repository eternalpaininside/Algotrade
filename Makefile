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
	  docker compose down algotrade-postgres && \
	  rm -rf out/pgdata &&\
	  echo "environment files was cleaned up"; \
	else \
	  echo "cleaning was canceled"; \
	fi
