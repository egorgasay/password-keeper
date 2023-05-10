BINARY_NAME=main.out

all: build test

build:
	go build -o ${BINARY_NAME} cmd/keeper/main.go

test:
	go test ./...

run:
	go build -o ${BINARY_NAME} cmd/keeper/main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}