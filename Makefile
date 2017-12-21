#!/usr/bin/env make

ml:
	go build -o ml.exe -ldflags "-X main.version=0.1 -X main.binary=ml" -v -x


%.so:
	go build -buildtype=plugin -o %.so


clean:
	rm ml.exe