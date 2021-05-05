# faceit-user-service

### How to run
```make run```

The command builds a docker image and runs a container in an attached mode.

### How to test
There is no special commands to run tests, they will be run during running the service (Please check Dockerfile for details)

### How to stop
Press ```Ctrl+C```. The container will be deleted automatically.

In order to delete docker images please run following command

```make clean```

### How to use
In order to test the service with swagger-ui please navigate to
[this page](http://localhost:3000/docs/index.html) after starting the docker container


### Explanations
Current service implementation shows the basic idea of splitting internal service parts to separate standalone components and using those components as a dependencies.

It's possible to solve that problem by different ways from project structure point of view.
I prefer to declare dependencies interfaces at the same place/module where it will be used, for example the user provider and event service interfaces were declared in api/api.go..

All parts of the service implemented in the simplest way without input data validation, strong error handling or other requirements for production ready code.
Tests are just an example how it's easy to write test code, mock some server dependencies (e.g the event service component).

### Possible future improvements
- add health check endpoints
- use a database instead of an in-memory storage
- add a users cache
- add auth/logging middlewares
- add an opportunity for configuring service (e.g changing port/log level/enabling profiler/DB address/etc)

