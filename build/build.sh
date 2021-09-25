#!/usr/bin/env bash

function build() {
	ROOT=$(dirname $0)
	NAME="edge-admin"
	DIST=$ROOT/"../dist/${NAME}"
	OS=${1}
	ARCH=${2}
	TAG=${3}

	if [ -z $OS ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z $ARCH ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z $TAG ]; then
		TAG="community"
	fi

	VERSION=$(lookup-version $ROOT/../internal/const/const.go)
	ZIP="${NAME}-${OS}-${ARCH}-${TAG}-v${VERSION}.zip"

	# build edge-api
	APINodeVersion=$(lookup-version $ROOT"/../../EdgeAPI/internal/const/const.go")
	echo "building edge-api v${APINodeVersion} ..."
	EDGE_API_BUILD_SCRIPT=$ROOT"/../../EdgeAPI/build/build.sh"
	if [ ! -f $EDGE_API_BUILD_SCRIPT ]; then
		echo "unable to find edge-api build script 'EdgeAPI/build/build.sh'"
		exit
	fi

	cd $ROOT"/../../EdgeAPI/build"
	echo "=============================="
	./build.sh $OS $ARCH $TAG
	echo "=============================="
	cd -

	# create dir & copy files
	echo "copying ..."
	if [ ! -d $DIST ]; then
		mkdir $DIST
		mkdir $DIST/bin
		mkdir $DIST/configs
		mkdir $DIST/logs
	fi

	cp -R $ROOT/../web $DIST/
	rm -f $DIST/web/tmp/*
	cp $ROOT/configs/server.template.yaml $DIST/configs/

	EDGE_API_ZIP_FILE=$ROOT"/../../EdgeAPI/dist/edge-api-${OS}-${ARCH}-${TAG}-v${APINodeVersion}.zip"
	cp $EDGE_API_ZIP_FILE $DIST/
	cd $DIST/
	unzip -q $(basename $EDGE_API_ZIP_FILE)
	rm -f $(basename $EDGE_API_ZIP_FILE)
	cd -

	# build
	echo "building "${NAME}" ..."
	env GOOS=$OS GOARCH=$GOARCH go build -tags $TAG -ldflags="-s -w" -o $DIST/bin/${NAME} $ROOT/../cmd/edge-admin/main.go

	# delete hidden files
	find $DIST -name ".DS_Store" -delete
	find $DIST -name ".gitignore" -delete
	find $DIST -name "*.less" -delete
	find $DIST -name "*.css.map" -delete

	# zip
	echo "zip files ..."
	cd "${DIST}/../" || exit
	if [ -f "${ZIP}" ]; then
		rm -f "${ZIP}"
	fi
	zip -r -X -q "${ZIP}" ${NAME}/
	rm -rf ${NAME}
	cd - || exit

	echo "[done]"
}

function lookup-version() {
	FILE=$1
	VERSION_DATA=$(cat $FILE)
	re="Version[ ]+=[ ]+\"([0-9.]+)\""
	if [[ $VERSION_DATA =~ $re ]]; then
		VERSION=${BASH_REMATCH[1]}
		echo $VERSION
	else
		echo "could not match version"
		exit
	fi
}

build $1 $2 $3
