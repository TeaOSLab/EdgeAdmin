#!/usr/bin/env bash

JS_ROOT=../web/public/js

echo "generating component.src.js ..."
env CGO_ENABLED=0 go run -tags=community ../cmd/edge-admin/main.go generate

if [ "$(which uglifyjs)" ]; then
	echo "compress to component.js ..."
	uglifyjs --compress --mangle -- ${JS_ROOT}/components.src.js > ${JS_ROOT}/components.js

	echo "compress to utils.min.js ..."
	uglifyjs --compress --mangle -- ${JS_ROOT}/utils.js > ${JS_ROOT}/utils.min.js
else
	echo "copy to component.js ..."
	cp ${JS_ROOT}/components.src.js ${JS_ROOT}/components.js

	echo "copy to utils.min.js ..."
	cp ${JS_ROOT}/utils.js ${JS_ROOT}/utils.min.js
fi

echo "ok"