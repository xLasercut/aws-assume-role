install:
	go mod tidy

build:
	go build -o ./dist/${BUILD_FILENAME} ./cmd/aws-assume-role

clean:
	rm -rf dist
