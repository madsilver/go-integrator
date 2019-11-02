.PHONY: start build depend clean configure config

APP_NAME=realtime-location-integrator
APP_PATH=bitbucket.org/picolotec/realtime-location-integrator
APP_VERSION=1.0.0
OS=linux

start: configure
	@cp ./config/config.yml.example ./config/config.yml

build: depend configure config
	@CGO_ENABLED=0 GOOS=${OS} go build \
    -o ./build/${APP_NAME} bitbucket.org/picolotec/realtime-location-integrator/cmd

configure:
	@dep ensure

depend:
	@mkdir -p ./build/config

config:
	@cp ./config/config.yml ./build/config

clean:
	@rm -fR vendor/ ./build
