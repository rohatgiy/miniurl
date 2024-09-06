d_pg:
	docker run --name miniurl_postgres \
	-e POSTGRES_USER=admin \
	-e POSTGRES_PASSWORD=root \
	-p 5432:5432 \
	-v $(PWD)/config/init.sql:/docker-entrypoint-initdb.d/init.sql \
	-d postgres:latest

d_redis:
	docker run --name miniurl_redis \
	-p 6379:6379 \
	-v $(PWD)/config/redis.conf:/usr/local/etc/redis/redis.conf \
	-d redis:latest \
	redis-server /usr/local/etc/redis/redis.conf

dc_up:
	ENV=prod docker-compose up --build

dev:
	go run *.go
