local-db:
	docker run \
  --name postgres \
  -it \
  --rm \
  -p 5432:5432 \
  -e POSTGRES_USER=${POSTGRES_USER} \
  -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
  -e POSTGRES_DB=${POSTGRES_DB} \
  -v ${PWD}/scripts/init-database.sh:/docker-entrypoint-initdb.d/init-database.sh \
  postgres:13-alpine