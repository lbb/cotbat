OUT:=main
DOCKER_NAME:=realfake/cotbat
CC=docker run --rm -e GOOS=linux -e CGO_ENABLED=0 -v "$(PWD)":/usr/src/$(OUT):z -w /usr/src/$(OUT) golang:1.9.2 go
GOBUILDFLAGS:=-i -v
LDFLAGS:=-extldflags '-static'

.PHONY: default
default: all

.PHONY: all
all: clean build

.PHONY: build
build: $(OUT)

.PHONY: docker
docker: build
	docker build . -t $(DOCKER_NAME):latest

$(OUT):
	$(CC) build $(GOBUILDFLAGS) -ldflags '-w $(LDFLAGS)'

.PHONY: clean
clean:
	rm -rf $(OUT)

.PHONY: run
run: docker
	docker run -it --rm --name cotbat -p 8080:80 -v `pwd`/log.log:/log.log $(DOCKER_NAME)

