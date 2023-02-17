install:
	cd src && go mod tidy

build:
	cd src && go build -o ../dist/${BUILD_FILENAME} .

clean:
	rm -rf dist
