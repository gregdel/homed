#!/bin/sh
set -e

BUILD_DIR="build"

_err() {
	echo "$@"
	exit 1
}

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

[ -d "./$BUILD_DIR" ] || _err "Missing build dir"
[ -d "frontend" ]     || _err "Frontend directory is missing"

_backend_build() {
	echo "Embeding files in the binary:"
	ls -lah "$BUILD_DIR"
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

_frontend_build
_backend_build

rm -rf "$BUILD_DIR"
