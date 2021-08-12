#!/usr/bin/env bash

ROOT=$(dirname $0)

# build all nodes
if [ -f $ROOT"/../../EdgeNode/build/build-all-plus.sh" ]; then
	echo "=============================="
	echo "build all edge-node"
	echo "=============================="
	cd $ROOT"/../../EdgeNode/build"
		./build-all-plus.sh
	cd -
fi

./build.sh linux amd64
./build.sh linux 386
./build.sh linux arm64
./build.sh linux mips64
./build.sh linux mips64le
./build.sh darwin amd64
./build.sh darwin arm64
