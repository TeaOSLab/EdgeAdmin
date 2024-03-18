// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build gcc

package injectionutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/waf/injectionutils"
	"github.com/iwind/TeaGo/assert"
	"runtime"
	"testing"
)

func TestDetectXSS(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsFalse(injectionutils.DetectXSS("", true))
	a.IsFalse(injectionutils.DetectXSS("abc", true))
	a.IsTrue(injectionutils.DetectXSS("<script>", true))
	a.IsTrue(injectionutils.DetectXSS("<link>", true))
	a.IsFalse(injectionutils.DetectXSS("<html><span>", true))
	a.IsFalse(injectionutils.DetectXSS("&lt;script&gt;", true))
	a.IsTrue(injectionutils.DetectXSS("/path?onmousedown=a", true))
	a.IsTrue(injectionutils.DetectXSS("/path?onkeyup=a", true))
	a.IsTrue(injectionutils.DetectXSS("onkeyup=a", true))
	a.IsTrue(injectionutils.DetectXSS("<iframe scrolling='no'>", true))
	a.IsFalse(injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span></body></html>", true))
	a.IsTrue(injectionutils.DetectXSS("name=s&description=%3Cscript+src%3D%22a.js%22%3Edddd%3C%2Fscript%3E", true))
	a.IsFalse(injectionutils.DetectXSS(`<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="XMP Core 6.0.0">
   <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
      <rdf:Description rdf:about=""
            xmlns:tiff="http://ns.adobe.com/tiff/1.0/">
         <tiff:Orientation>1</tiff:Orientation>
      </rdf:Description>
   </rdf:RDF>
</x:xmpmeta>`, true)) // included in some photo files
	a.IsFalse(injectionutils.DetectXSS(`<xml></xml>`, false))
}

func TestDetectXSS_Strict(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsFalse(injectionutils.DetectXSS(`<xml></xml>`, false))
	a.IsTrue(injectionutils.DetectXSS(`<xml></xml>`, true))
	a.IsFalse(injectionutils.DetectXSS(`<img src=\"\"/>`, false))
	a.IsFalse(injectionutils.DetectXSS(`<img src=\"test.jpg\"/>`, true))
	a.IsFalse(injectionutils.DetectXSS(`<a href="aaaa"></a>`, true))
	a.IsFalse(injectionutils.DetectXSS(`<span style="color: red"></span>`, false))
	a.IsTrue(injectionutils.DetectXSS(`<span style="color: red"></span>`, true))
	a.IsFalse(injectionutils.DetectXSS("https://example.com?style=list", false))
	a.IsTrue(injectionutils.DetectXSS("https://example.com?style=list", true))
}

func BenchmarkDetectXSS_MISS(b *testing.B) {
	var result = injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span></body></html>", false)
	if result {
		b.Fatal("'result' should not be 'true'")
	}

	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span></body></html>", false)
		}
	})
}

func BenchmarkDetectXSS_MISS_Cache(b *testing.B) {
	var result = injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span></body></html>", false)
	if result {
		b.Fatal("'result' should not be 'true'")
	}

	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectXSSCache("<html><body><span>RequestId: 1234567890</span></body></html>", false, 1800)
		}
	})
}

func BenchmarkDetectXSS_HIT(b *testing.B) {
	var result = injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span><script src=\"\"></script></body></html>", false)
	if !result {
		b.Fatal("'result' should not be 'false'")
	}

	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectXSS("<html><body><span>RequestId: 1234567890</span><script src=\"\"></script></body></html>", false)
		}
	})
}
