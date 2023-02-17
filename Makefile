install:
	cd src && go mod tidy

build: clean
	cd src && go build -o ../dist/${BUILD_FILENAME} .

clean:
	rm -rf dist
