#!/bin/sh

BUILD_DIR="build"

_err() {
	echo "$@"
	exit 1
}

_backend_build() {
	echo "Running tests..."
	CGO_ENABLED=0 go test ./... || _err "Backend test failed"

	echo "Embeding files in the binary:"
	ls -lah "$BUILD_DIR"
	CGO_ENABLED=0 go build \
		-ldflags '-extldflags "-static"' \
		-trimpath \
		-v \
		-o "homed" \
		. || _err "Backend build failed"
}

_frontend_build() {
	cd frontend || return
	npm install || _err "npm install failed"
	npm run-script build || _err "frondend build failed"
	cd .. || return
}

case "$1" in
	frontend)
		_frontend_build
		;;
	backend)
		_backend_build
		;;
	*)
		_frontend_build
		_backend_build
		;;
esac
