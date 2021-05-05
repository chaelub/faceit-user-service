#!/bin/bash

run:
	docker build -t faceit-user-service-image .
	docker run -p 3000:3000 faceit-user-service-image