bin := papermc-manager
dist := dist

.PHONY: all
all: build

.PHONY: build
build: clean
	go build -o $(dist)/$(bin)

.PHONY: clean
clean:
	rm -rf $(dist)/*

.PHONY: run
run: build
	$(dist)/$(bin)
