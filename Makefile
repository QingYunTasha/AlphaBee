build:
	docker build -t alpha-bee:latest -f Dockerfile/web.dockerfile .

push:
	docker push web:latest

benchmark:
	k6 run benchmark/index.js

all: build push