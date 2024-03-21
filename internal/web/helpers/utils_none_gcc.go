// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !gcc

package helpers

import (
	"net/http"
)

// filter request
func safeFilterRequest(req *http.Request) bool {
	return true
}
