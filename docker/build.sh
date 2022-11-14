#!/usr/bin/env bash

VERSION=latest

docker build --no-cache -t goedge/edge-admin:${VERSION} .
