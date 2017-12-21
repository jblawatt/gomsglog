#!/usr/bin/env make

ml:
	go build -o ml.exe


%.so:
	go build -buildtype=plugin -o %.so


clean:
	rm ml.exe