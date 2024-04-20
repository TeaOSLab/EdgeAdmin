// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package otputils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/otputils"
	"testing"
)

func TestFixIssuer(t *testing.T) {
	var beforeURL = "otpauth://totp/GoEdge%25E7%25AE%25A1%25E7%2590%2586%25E5%2591%2598%25E7%25B3%25BB%25E7%25BB%259F:admin?issuer=GoEdge%25E7%25AE%25A1%25E7%2590%2586%25E5%2591%2598%25E7%25B3%25BB%25E7%25BB%259F&secret=Q3J4WNOWBRFLP3HI"
	var afterURL = otputils.FixIssuer(beforeURL)
	t.Log(afterURL)

	if beforeURL == afterURL {
		t.Fatal("'afterURL' should not be equal to 'beforeURL'")
	}
}
