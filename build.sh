#!/bin/sh

BUILD_DIR="build"

_err() {
	echo "$@"
	exit 1
}

[ -d "./$BUILD_DIR" ] || _err "Missing build dir"
[ -d "frontend" ]     || _err "Frontend directory is missing"

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

_backend_build() {
	CGO_ENABLED=0 go build \
		-ldflags '-extldflags "-static"' \
		-trimpath \
		-v \
		-o "$BUILD_DIR/homed" \
		.
}

_frontend_build() {
	cd frontend || return
	npm install
	npm run-script build
}

_backend_build
_frontend_build
