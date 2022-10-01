BIN_DIR:= ./build/bin
CONTAINER_NAME:=azuki774/dropbox-uploader
.PHONY: build bin test start stop clean

build:
	make bin
	docker build -t $(CONTAINER_NAME):latest -f build/Dockerfile .

bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\"" -o build/bin/ ./...

test:
	go test -v ./...

start:
	docker compose -f deployment/compose-local.yml up -d

stop:
	docker compose -f deployment/compose-local.yml down

clean:
	rm -rf $(BIN_DIR)/*
	docker images -a | grep ${CONTAINER_NAME} | awk '{print $$3}' | xargs docker rmi
