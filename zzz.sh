#!/bin/bash

goreleaser release --skip-publish --rm-dist --snapshot

if [ -f dist/nezha-darwin-amd64.zip ];then
	cd dist
	mkdir -p tmp
	cp -rf nezha-darwin-amd64.zip tmp/nezha-darwin-amd64.zip
	cd tmp && unzip nezha-darwin-amd64.zip
	./nezha web
fi