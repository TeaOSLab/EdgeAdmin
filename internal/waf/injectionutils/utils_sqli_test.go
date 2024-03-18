// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build gcc

package injectionutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/waf/injectionutils"
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"runtime"
	"strings"
	"testing"
)

func TestDetectSQLInjection(t *testing.T) {
	var a = assert.NewAssertion(t)
	for _, isStrict := range []bool{true, false} {
		a.IsTrue(injectionutils.DetectSQLInjection("' UNION SELECT * FROM myTable", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("id=1 ' UNION  select * from a", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("asdf asd ; -1' and 1=1 union/* foo */select load_file('/etc/passwd')--", isStrict))
		a.IsFalse(injectionutils.DetectSQLInjection("' UNION SELECT1 * FROM myTable", isStrict))
		a.IsFalse(injectionutils.DetectSQLInjection("1234", isStrict))
		a.IsFalse(injectionutils.DetectSQLInjection("", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("id=123 OR 1=1&b=2", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("id=123&b=456&c=1' or 2=2", isStrict))
		a.IsFalse(injectionutils.DetectSQLInjection("?", isStrict))
		a.IsFalse(injectionutils.DetectSQLInjection("/hello?age=22", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("/sql/injection?id=123 or 1=1", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("/sql/injection?id=123%20or%201=1", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("https://example.com/sql/injection?id=123%20or%201=1", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("id=123%20or%201=1", isStrict))
		a.IsTrue(injectionutils.DetectSQLInjection("https://example.com/' or 1=1", isStrict))
	}
}

func BenchmarkDetectSQLInjection(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("asdf asd ; -1' and 1=1 union/* foo */select load_file('/etc/passwd')--", false)
		}
	})
}

func BenchmarkDetectSQLInjection_URL(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("/sql/injection?id=123 or 1=1", false)
		}
	})
}

func BenchmarkDetectSQLInjection_Normal_Small(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("a/sql/injection?id=1234", false)
		}
	})
}

func BenchmarkDetectSQLInjection_URL_Normal_Small(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("/sql/injection?id="+types.String(rands.Int64()%10000), false)
		}
	})
}

func BenchmarkDetectSQLInjection_URL_Normal_Middle(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("/search?q=libinjection+fingerprint&newwindow=1&sca_esv=589290862&sxsrf=AMwHvKnxuLoejn2XlNniffC12E_xc35M7Q%3A1702090118361&ei=htvzzebfFZfo1e8PvLGggAk&ved=0ahUKEwjTsYmnq4GDAxUWdPOHHbwkCJAQ4ddDCBA&uact=5&oq=libinjection+fingerprint&gs_lp=Egxnd3Mtd2l6LXNlcnAiGIxpYmluamVjdGlvbmBmaW5nKXJwcmludTIEEAAYHjIGVAAYCBgeSiEaUPkRWKFZcAJ4AZABAJgBHgGgAfoEqgwDMC40uAEGyAEA-AEBwgIKEAFYTxjWMuiwA-IDBBgAVteIBgGQBgI&sclient=gws-wiz-serp#ip=1", false)
		}
	})
}

func BenchmarkDetectSQLInjection_URL_Normal_Small_Cache(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjectionCache("/sql/injection?id="+types.String(rands.Int64()%10000), false, 1800)
		}
	})
}

func BenchmarkDetectSQLInjection_Normal_Large(b *testing.B) {
	runtime.GOMAXPROCS(4)

	var s = strings.Repeat("A", 512)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("a/sql/injection?id="+types.String(rands.Int64()%10000)+"&s="+s+"&v=%20", false)
		}
	})
}

func BenchmarkDetectSQLInjection_Normal_Large_Cache(b *testing.B) {
	runtime.GOMAXPROCS(4)

	var s = strings.Repeat("A", 512)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjectionCache("a/sql/injection?id="+types.String(rands.Int64()%10000)+"&s="+s, false, 1800)
		}
	})
}

func BenchmarkDetectSQLInjection_URL_Unescape(b *testing.B) {
	runtime.GOMAXPROCS(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = injectionutils.DetectSQLInjection("/sql/injection?id=123%20or%201=1", false)
		}
	})
}
