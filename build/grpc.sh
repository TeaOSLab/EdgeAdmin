#!/usr/bin/env bash

function proto_compile() {
	protoc --go_out=plugins=grpc:../internal/rpc --proto_path=../internal/rpc ../internal/rpc/${1}/*.proto
}

proto_compile "admin"
#proto_compile "dns"
#proto_compile "log"
#proto_compile "provider"
#proto_compile "stat"
#proto_compile "user"
#proto_compile "monitor"
#proto_compile "node"