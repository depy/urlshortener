build:
	docker build -t urlshortener .
run:
	docker run -p 8080:8080 -it urlshortener