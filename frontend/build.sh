#!/bin/sh

buildir="build"

_build() {
	[ -d "$outdir" ] && rm -rf "$outdir"
	mkdir -p "$outdir"

	cp src/assets/* "$outdir"
	npx esbuild \
		--loader:.js=jsx \
		--bundle \
		--outdir="$outdir" \
		src/js/app.js \
		"$@"
}

case "$1" in
	--dev)
		outdir="$buildir"
		_build --watch --sourcemap
		;;
	*)
		outdir="../$buildir"
		_build --minify
		;;
esac
