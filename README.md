# go-integrator
This is a listener for events from the realtime_locations table in Postgres that sends records to a queue on RabbitMQ server and Redis server.

## Requirements
App requires Golang 1.8 or later and Dep Package Manager

## Installation
- Install [Golang](https://golang.org/doc/install)
- Install [Dep](https://github.com/golang/dep)

```
# Prepare the project for development
make start
```

## Build
```
# Build the binary in your environment
make build

# Build with another OS. Default Linux
$ make OS=darwin build

# Clean Up
$ make clean

# Configure. Install app dependencies.
$ make configure
```

## Run
```
docker-compose up -d
```

## Development
```
# Run the application without build
go run ./cmd/*.go -client=<client> -env=dev
```

## Web GUI
### Supervisior
http://localhost:8091
### RabbitMQ
http://localhost:15672