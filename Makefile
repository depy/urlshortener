build:
	docker build -t urlshortener .
run:
	make build
	docker-compose up