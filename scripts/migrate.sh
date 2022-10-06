#!/bin/sh

set -e
echo "specified environment is: ${ENVIRONMENT}"

(bin/yq --version >/dev/null && bin/goose --version >/dev/null) || (echo "installing yq to parse yaml files" && make bin-deps)

if [ "${ENVIRONMENT}" != "local" ] && [ "${ENVIRONMENT}" != "test" ] && [ "${ENVIRONMENT}" != "production" ]; then
  echo "Unknown environment: ${ENVIRONMENT}"
  exit 1
fi

if [ "${ENVIRONMENT}" = "local" ]; then
  CONFIG_FILE="deploy/configs/values_local.yaml"
else
  if [ "${ENVIRONMENT}" = "production" ]; then
    CONFIG_FILE="deploy/configs/values_production.yaml"
  else
    if [ "${ENVIRONMENT}" = "test" ]; then
      CONFIG_FILE="deploy/configs/values_docker.yaml"
    fi
  fi
fi
MIGRATION_DIR=$(bin/yq eval '.postgres.migrations*' "${CONFIG_FILE}")

# билдим DB_DSN
DB_HOST=$(bin/yq eval '.postgres.migration_host' "${CONFIG_FILE}")
DB_PORT=$(bin/yq eval '.postgres.migration_port' "${CONFIG_FILE}")
DB_USER=$(bin/yq eval '.postgres.user' "${CONFIG_FILE}")
DB_PASSWORD=$(bin/yq eval '.postgres.password' "${CONFIG_FILE}")
DB_NAME=$(bin/yq eval '.postgres.database' "${CONFIG_FILE}")
echo "DB_HOST: ${DB_HOST}"
echo "DB_PORT: ${DB_PORT}"
echo "DB_USER: ${DB_USER}"
echo "DB_NAME: ${DB_NAME}"
DB_DSN="user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} host=${DB_HOST} port=${DB_PORT} sslmode=disable"

echo "running: ${MIGRATION_DIR}" postgres "${DB_DSN}"
if [ "$1" = "--dryrun" ]; then
  ./bin/goose -dir "${MIGRATION_DIR}" postgres "${DB_DSN}" status
else
  ./bin/goose -dir "${MIGRATION_DIR}" postgres "${DB_DSN}" up
fi
