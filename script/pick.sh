#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

# https://github.com/FiloSottile/homebrew-musl-cross
# brew install FiloSottile/musl-cross/musl-cross --without-x86_64 --with-i486 --with-aarch64 --with-arm

# brew install mingw-w64
# sudo port install mingw-w64

VERSION=v0.15.2
curPath=`pwd`
rootPath=$(dirname "$curPath")

PACK_NAME=nezha

# go tool dist list
mkdir -p $rootPath/tmp/build
mkdir -p $rootPath/tmp/package

source ~/.bash_profile

cd $rootPath

echo $LDFLAGS
build_app(){

	if [ -f $rootPath/tmp/build/${PACK_NAME} ]; then
		rm -rf $rootPath/tmp/build/${PACK_NAME}
		rm -rf $rootPath/${PACK_NAME}
	fi

	if [ -f $rootPath/tmp/build/${PACK_NAME}.exe ]; then
		rm -rf $rootPath/tmp/build/${PACK_NAME}.exe
		rm -rf $rootPath/${PACK_NAME}.exe
	fi

	echo "build_app" $1 $2

	echo "export CGO_ENABLED=1 GOOS=$1 GOARCH=$2"
	echo "cd $rootPath && go build ${PACK_NAME}.go"

	# export CGO_ENABLED=1 GOOS=linux GOARCH=amd64

	if [ $1 != "darwin" ];then
		export CGO_ENABLED=1 GOOS=$1 GOARCH=$2
		export CGO_LDFLAGS="-static"
	fi

	if [ $1 == "windows" ];then
		
		if [ $2 == "amd64" ]; then
			export CC=x86_64-w64-mingw32-gcc
			export CXX=x86_64-w64-mingw32-g++
		else
			export CC=i686-w64-mingw32-gcc
			export CXX=i686-w64-mingw32-g++
		fi

		cd $rootPath && go build -o ${PACK_NAME}.exe -ldflags "${LDFLAGS}" ${PACK_NAME}.go
	fi

	if [ $1 == "linux" ]; then
		export CC=x86_64-linux-musl-gcc
		if [ $2 == "amd64" ]; then
			export CC=x86_64-linux-musl-gcc

		fi

		if [ $2 == "386" ]; then
			export CC=i486-linux-musl-gcc
		fi

		if [ $2 == "arm64" ]; then
			export CC=aarch64-linux-musl-gcc
		fi

		if [ $2 == "arm" ]; then
			export CC=arm-linux-musleabi-gcc
		fi

		cd $rootPath && go build -ldflags "${LDFLAGS}"  main.go 
	fi

	if [ $1 == "darwin" ]; then
		echo "cd $rootPath && go build -v -ldflags '${LDFLAGS}'"
		cd $rootPath && go build -v -ldflags "${LDFLAGS}"
		
		cp $rootPath/${PACK_NAME} $rootPath/tmp/build
	fi
	

	cp -rf $rootPath/resource $rootPath/tmp/build
	cp -rf $rootPath/data $rootPath/tmp/build
	rm -rf $rootPath/tmp/build/data/sqlite.db


	if [ $1 == "windows" ];then
		cp $rootPath/${PACK_NAME}.exe $rootPath/tmp/build
	else
		cp $rootPath/${PACK_NAME} $rootPath/tmp/build
	fi

	# zip
	cd $rootPath/tmp/build && zip -r -q -o ${PACK_NAME}_$1_$2.zip  ./ && mv ${PACK_NAME}_$1_$2.zip $rootPath/tmp/package
}

golist=`go tool dist list`
echo $golist

# build_app linux amd64
# build_app linux 386
# build_app linux arm64
# build_app linux arm
build_app darwin amd64
# build_app windows 386
# build_app windows amd64

