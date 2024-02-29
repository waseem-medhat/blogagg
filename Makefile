include .env

up:
	cd sql/schema && goose postgres ${DBURL} up

down:
	cd sql/schema && goose postgres ${DBURL} down

lint:
	golangci-lint run
