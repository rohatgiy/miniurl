package main

func main() {
	redis := initRedis()
	postgres := initPostgres()
	router := initRouter(redis, postgres)

	router.Run()
}
