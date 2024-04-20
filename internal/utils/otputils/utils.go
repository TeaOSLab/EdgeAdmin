// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package otputils

import (
	"net/url"
)

// FixIssuer fix issuer in otp provisioning url
func FixIssuer(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		return urlString
	}

	var query = u.Query()

	if query != nil {
		var issuerName = query.Get("issuer")
		if len(issuerName) > 0 {
			unescapedIssuerName, unescapeErr := url.QueryUnescape(issuerName)
			if unescapeErr == nil {
				query.Set("issuer", unescapedIssuerName)
				u.RawQuery = query.Encode()
			}
		}
		return u.String()
	}

	return urlString
}
