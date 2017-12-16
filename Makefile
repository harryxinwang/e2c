.PHONY: default
default: build

.PHONY: build
build: _output/e2c
	echo "Build is done!"

_output/e2c:
	mkdir -p ./_output
	go build -o ./_output/e2c .

.PHONY: clear
clear:
	rm -rf _output
