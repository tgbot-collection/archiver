OS = darwin linux windows
ARCH = amd64
default:
	git pull
	@echo "Build current platform executable..."
	go build  .

static:
	@echo "Build static files..."
	CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' .

all:
	git pull
	@echo "Build all platform executables..."
	@for o in $(OS) ; do            \
    		for a in $(ARCH) ; do     \
    			CGO_ENABLED=0 GOOS=$$o GOARCH=$$a go build -ldflags="-s -w" -o builds/DailyGakki-$$o-$$a .;    \
    		done                              \
    	done




clean:
	@rm -rf builds
	@rm -f archiver
