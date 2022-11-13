#!/usr/bin/env bash

VERSION=latest

docker run -d -p 7788:7788 -p 8001:8001 -p 3306:3306 --name edge-admin goedge/edge-admin:${VERSION}
