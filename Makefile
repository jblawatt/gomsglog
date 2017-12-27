#!/usr/bin/env make

ml:
	go build -o ml.exe -ldflags "-X main.version=0.1 -X main.binary=ml" -v -x

js:
	node node_modules/.bin/babel static/js/gomsglog.jsx --out-dir . --source-maps

%.so:
	go build -buildtype=plugin -o %.so


clean:
	rm ml.exe