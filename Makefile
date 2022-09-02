BIN_DIR:= ./build/bin
CONTAINER_NAME:=azuki774/dropbox-uploader
.PHONY: build bin

build:
	make bin
	docker build -t azuki774/$(CONTAINER_NAME):latest -f build/Dockerfile .

bin:
	CGO_ENABLED=0 go build -o $(BIN_DIR)/ .

test:
	go test -v ./...

clean:
	rm -rf $(BIN_DIR)/*
