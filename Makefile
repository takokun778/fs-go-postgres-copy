export CONTAINER_NAME := fs-go-postgres-copy

docker:
	@docker run --rm -d \
		-p 15432:5432 \
		-e TZ=UTC \
		-e LANG=ja_JP.UTF-8 \
		-e POSTGRES_HOST_AUTH_METHOD=trust \
		-e POSTGRES_DB=postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_INITDB_ARGS=--encoding=UTF-8 \
		--name $(CONTAINER_NAME) \
		postgres:14.2-alpine
psql:
	docker exec -it $(CONTAINER_NAME) psql -U postgres
stop:
	docker stop $(CONTAINER_NAME)
