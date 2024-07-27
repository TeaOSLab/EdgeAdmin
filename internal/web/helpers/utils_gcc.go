// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .
//go:build gcc

package helpers

import (
	"net/http"

	"github.com/TeaOSLab/EdgeAdmin/internal/waf/injectionutils"
)

// filter request
func safeFilterRequest(req *http.Request) bool {
	return !injectionutils.DetectXSS(req.RequestURI, false) && !injectionutils.DetectSQLInjection(req.RequestURI, false)
}
