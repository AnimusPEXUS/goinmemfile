export GONOPROXY=https://github.com/AnimusPEXUS/*

all: get build

get:
		$(MAKE) -C tests/RO get
		go get -u -v "./..."
		go mod tidy

build:
		$(MAKE) -C tests/RO build
		go build

