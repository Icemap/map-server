.PHONY: all build clean

all:
	make clean build

build:
	go build -o bin/map-server

clean:
	rm -rf bin

yaml:
	go install github.com/Icemap/yaml2go-cli@latest
	yaml2go-cli -i config/config.yaml -o config/config.bean.go -p config -s Config