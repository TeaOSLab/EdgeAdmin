#!/usr/bin/env bash

JS_ROOT=../web/public/js

echo "generate component.src.js ..."
go run -tags=community ../cmd/edge-admin/main.go  generate

if [ `which uglifyjs` ]; then
	echo "compress to component.js ..."
	uglifyjs --compress --mangle -- ${JS_ROOT}/components.src.js > ${JS_ROOT}/components.js
else
	echo "copy to component.js ..."
	cp ${JS_ROOT}/components.src.js ${JS_ROOT}/components.js
fi

echo "ok"