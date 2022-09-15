OS = darwin linux windows
ARCH = amd64 arm64
export DRIVER=$(shell which chromedriver)

default:
	@echo "Build current platform executable..."
	go build  .

static:
	@echo "Build static files..."
	CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' .

all:
	make clean
	@echo "Build all platform executables..."
	@for o in $(OS) ; do            \
    		for a in $(ARCH) ; do     \
    		  	echo "Building $$o-$$a..."; \
    		  	if [ "$$o" = "windows" ]; then \
                	CGO_ENABLED=0 GOOS=$$o GOARCH=$$a go build -ldflags="-s -w" -o builds/archiver-$$o-$$a.exe .;    \
                else \
    				CGO_ENABLED=0 GOOS=$$o GOARCH=$$a go build -ldflags="-s -w" -o builds/archiver-$$o-$$a .;    \
    			fi; \
    		done   \
    	done

	@make universal
	@make checksum

clean:
	@rm -rf builds
	@rm -f archiver

test:
	go test -v -coverprofile=coverage.txt -covermode=atomic

checksum: builds/*
	@echo "Generating checksums..."
	if [ "$(shell uname)" = "Darwin" ]; then \
		shasum -a 256 $^ >>  builds/checksum-sha256sum.txt ;\
	else \
		sha256sum  $^ >> builds/checksum-sha256sum.txt; \
	fi


universal:
	@echo "Building macOS universal binary..."
	docker run --rm -v $(shell pwd)/builds:/app/ bennythink/lipo-linux -create -output \
		archiver-darwin-universal \
		archiver-darwin-amd64    archiver-darwin-arm64

	file builds/archiver-darwin-universal

release:
	git tag $(shell git rev-parse --short HEAD)
	git push --tags