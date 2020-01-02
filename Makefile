#!/usr/bin/env make

ml:
	go build -o ml.exe -ldflags "-X main.version=0.1 -X main.binary=ml" -x

bindata:
	# ....

node_modules:
	npm install

js: node_modules
	npm run build

%.so:
	go build -buildtype=plugin -o %.so


clean:
	rm ml.exe
	rm node_modules


serve:
	go run main.go serve


.PHONY: clean
