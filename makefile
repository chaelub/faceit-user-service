#!/bin/bash

run:
	make build
	make start

build:
	docker build -t faceit-user-service-image .

start:
	docker run -p 3000:3000 -it --rm faceit-user-service-image

stop:
	docker stop faceit-user-service-image

clean:
	docker images -a | grep "faceit-user-service-image" | awk '{print $$3}' | xargs docker rmi --force

