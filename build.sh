#!/bin/sh
set -e

BUILD_DIR="build"
mkdir -p "$BUILD_DIR"

_cleanup_buildir() {
	rm -rf "$BUILD_DIR"
	mkdir -p "$BUILD_DIR"
	[ "$1" ] || return 0
	touch "$BUILD_DIR/$1"
}

_backend_build() {
	echo "Embeding files in the binary:"
	ls -lah "$BUILD_DIR"
	CGO_ENABLED=0 go test -v ./...
	CGO_ENABLED=0 go build \
		-ldflags '-extldflags "-static"' \
		-trimpath \
		-v \
		-o "homed" \
		.
}

_frontend_build() {
	cd frontend || return
	npm install
	npm run-script build
	cd .. || return
}

case "$1" in
	frontend)
		_cleanup_buildir
		_frontend_build
		;;
	backend)
		_backend_build
		;;
	*)
		_cleanup_buildir
		_frontend_build
		_backend_build
		_cleanup_buildir keep
		;;
esac
