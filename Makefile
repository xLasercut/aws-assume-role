.PHONY: install
install:
	go mod tidy

.PHONY: build
build:
	go build -o ./dist/${BUILD_FILENAME} ./cmd/aws-assume-role

.PHONY: build-ci
build-ci: guard-BUILD_FILENAME guard-VERSION_NUMBER
	go build -ldflags="-X main.AppVersion=${VERSION_NUMBER}" -o ./dist/${BUILD_FILENAME} ./cmd/aws-assume-role

.PHONY: clean
clean:
	rm -rf dist

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
