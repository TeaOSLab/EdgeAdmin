#!/usr/bin/env bash

VERSION=latest

docker build -t goedge/edge-admin:${VERSION} .
