#!/bin/bash

if [ -z $CI_PROJECT_DIR ]
then
    echo "Missing CI_PROJECT_DIR"
    exit 1
fi
WORMHOLE_WIN_64="wormhole.i64.exe"
WORMHOLE_LINUX_64="wormhole.i64"
WORMHOLE_LINUX_ARM="wormhole.arm.i386"
WORMHOLE_LINUX_ARM_64="wormhole.arm.i64"

build_windows() {
    GOOS="windows"
    GOARCH="amd64"
    go build -o $CI_PROJECT_DIR/bin/$WORMHOLE_WIN_64
}

build_linux() {
    GOOS="linux"
    GOARCH="amd64"
    go build -o $CI_PROJECT_DIR/bin/$WORMHOLE_LINUX_64
}

build_arm() {
    GOOS=linux
    GOARCH=arm
    go build -o $CI_PROJECT_DIR/bin/$WORMHOLE_LINUX_ARM
    GOOS=linux
    GOARCH=arm64
    go build -o $CI_PROJECT_DIR/bin/$WORMHOLE_LINUX_ARM_64
}

cd $CI_PROJECT_DIR/go
if [ "$1" == "win" ] || [ "$1" == "all" ]
then
    build_windows
fi
if [ "$1" == "linux" ] || [ "$1" == "all" ]
then
    build_linux
fi
if [ "$1" == "arm" ] || [ "$1" == "all" ]
then
    build_arm
fi