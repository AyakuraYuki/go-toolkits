GOOS?=darwin
GOARCH?=arm64
NAME=ffmpeg-make
BASE_BUILD_DIR=build
BUILD_NAME=$(GOOS)-$(GOARCH)$(GOARM)
BUILD_DIR=$(BASE_BUILD_DIR)/$(BUILD_NAME)
VERSION?=0.0.0

ifeq ($(GOOS),windows)
  ext=.exe
  archiveCmd=zip -9 -r $(NAME)-$(BUILD_NAME)-$(VERSION).zip $(BUILD_NAME)
else
  ext=
  archiveCmd=tar czpvf $(NAME)-$(BUILD_NAME)-$(VERSION).tar.gz $(BUILD_NAME)
endif

build: clean test
	@go build -o $(BUILD_DIR)/$(NAME)$(ext)

test:
	@go test -race -v -bench=. ./...

release: check-env-release
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(NAME)$(ext)
	@cd $(BASE_BUILD_DIR) ; $(archiveCmd)

check-env-release:
	@if [ "$(GOOS)" = "" ]; then \
		echo "Env variable GOOS not set"; \
		exit 1; \
   	fi
	@if [ "$(GOARCH)" = "" ]; then \
		echo "Env variable GOARCH not set"; \
		exit 1; \
   	fi

clean:
	@go clean
	@rm -vrf $(BASE_BUILD_DIR)
