APP_NAME=apimock-go

all: clean build

clean:
	rm -f $(APP_NAME)

fmt:
	gofmt -w ./
	goimports -w ./

run:
	go run ./cmd/server.go -root ./root

build: clean
	go build -o $(APP_NAME) ./cmd/server.go

test: install
	GOPATH=$(VENDOR_GOPATH):$(ROOT_GOPATH) go test -v ./src/...
