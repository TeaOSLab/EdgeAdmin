// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build gcc

package injectionutils

/*
#cgo CFLAGS: -O2 -I./libinjection/src

#include <libinjection.h>
#include <stdlib.h>
*/
import "C"
import (
	"net/url"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// DetectSQLInjectionCache detect sql injection in string with cache
func DetectSQLInjectionCache(input string, isStrict bool, cacheLife int) bool {
	var l = len(input)

	if l == 0 {
		return false
	}

	if cacheLife <= 0 || l < 128 || l > MaxCacheDataSize {
		return DetectSQLInjection(input, isStrict)
	}

	var result = DetectSQLInjection(input, isStrict)
	return result
}

// DetectSQLInjection detect sql injection in string
func DetectSQLInjection(input string, isStrict bool) bool {
	if len(input) == 0 {
		return false
	}

	if !isStrict {
		if len(input) > 1024 {
			if !utf8.ValidString(input[:1024]) && !utf8.ValidString(input[:1023]) && !utf8.ValidString(input[:1022]) {
				return false
			}
		} else {
			if !utf8.ValidString(input) {
				return false
			}
		}
	}

	if detectSQLInjectionOne(input) {
		return true
	}

	// 兼容 /PATH?URI
	if (input[0] == '/' || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")) && len(input) < 1024 {
		var argsIndex = strings.Index(input, "?")
		if argsIndex > 0 {
			var args = input[argsIndex+1:]
			unescapeArgs, err := url.QueryUnescape(args)
			if err == nil && args != unescapeArgs {
				return detectSQLInjectionOne(args) || detectSQLInjectionOne(unescapeArgs)
			} else {
				return detectSQLInjectionOne(args)
			}
		}
	} else {
		unescapedInput, err := url.QueryUnescape(input)
		if err == nil && input != unescapedInput {
			return detectSQLInjectionOne(unescapedInput)
		}
	}

	return false
}

func detectSQLInjectionOne(input string) bool {
	if len(input) == 0 {
		return false
	}

	var fingerprint [8]C.char
	var fingerprintPtr = (*C.char)(unsafe.Pointer(&fingerprint[0]))
	var cInput = C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	return C.libinjection_sqli(cInput, C.size_t(len(input)), fingerprintPtr) == 1
}
