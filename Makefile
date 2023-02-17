install:
	cd src && go mod tidy

build: clean
	cd src && go build -o ../dist/ .

clean:
	rm -rf dist
