docker_run_pg:
	docker run --name miniurl_postgres \
	-e POSTGRES_USER=admin \
	-e POSTGRES_PASSWORD=root \
	-p 5432:5432 \
	-v $(PWD)/config/init.sql:/docker-entrypoint-initdb.d/init.sql \
	-d postgres:latest

docker_run_redis:
	docker run --name miniurl_redis \
	-p 6379:6379 \
	-v $(PWD)/config/redis.conf:/usr/local/etc/redis/redis.conf \
	-d redis:latest \
	redis-server /usr/local/etc/redis/redis.conf

dev:
	go run *.go
