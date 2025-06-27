BINARY=.bin/backend
BACKEND_DIR=backend

.PHONY: all build run clean

all: build

build:
	mkdir -p .bin
	cd $(BACKEND_DIR) && go build -o ../$(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -rf .bin 