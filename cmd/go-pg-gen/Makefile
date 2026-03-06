NAME=go-pg-gen
OUT_DIR=build
TARGET=$(GOOS)-$(GOARCH)$(GOARM)
BIN_DIR=$(OUT_DIR)/$(TARGET)

ifeq ($(GOOS),windows)
  EXT=.exe
  PACK_CMD=zip -9 -r $(NAME)-$(TARGET)-$(VERSION).zip $(TARGET)
else
  EXT=
  PACK_CMD=tar czpvf $(NAME)-$(TARGET)-$(VERSION).tar.gz $(TARGET)
endif

.PHONY: build
build: clean
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BIN_DIR)/$(NAME)$(EXT)
	@chmod +x $(BIN_DIR)/$(NAME)$(EXT)
	@echo "binary file: $(BIN_DIR)/$(NAME)$(EXT)"

.PHONY: clean
clean:
	@go clean
	@rm -rf $(OUT_DIR)
