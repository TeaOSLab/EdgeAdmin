package helpers

import (
	"bytes"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/regionconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

var ipCacheMap = map[string]bool{} // ip => bool
var ipCacheLocker = sync.Mutex{}

func init() {
	events.On(events.EventSecurityConfigChanged, func() {
		ipCacheLocker.Lock()
		ipCacheMap = map[string]bool{}
		ipCacheLocker.Unlock()
	})
}

// DisableXFFPrompt 停用XFF提示
func DisableXFFPrompt() {
	securityXFFPromptDisabled = true
}

// 检查用户IP并支持缓存
func checkIP(config *systemconfigs.SecurityConfig, ipAddr string) bool {
	ipCacheLocker.Lock()
	ipCache, ok := ipCacheMap[ipAddr]
	if ok && ipCache {
		ipCacheLocker.Unlock()
		return ipCache
	}
	ipCacheLocker.Unlock()

	result := checkIPWithoutCache(config, ipAddr)
	ipCacheLocker.Lock()

	// 缓存的内容不能过多
	if len(ipCacheMap) > 100_000 {
		ipCacheMap = map[string]bool{}
	}

	ipCacheMap[ipAddr] = result
	ipCacheLocker.Unlock()
	return result
}

// 检查用户IP
func checkIPWithoutCache(config *systemconfigs.SecurityConfig, ipAddr string) bool {
	if config == nil {
		return true
	}

	// 本地IP
	ipObj := net.ParseIP(ipAddr)
	if ipObj == nil {
		logs.Println("[USER_MUST_AUTH]parse ip: invalid client address: " + ipAddr)
		return false
	}
	ip := ipObj.To4()
	if ip == nil {
		// IPv6
		ip = ipObj.To16()
		if ip == nil {
			logs.Println("[USER_MUST_AUTH]invalid client address: " + ipAddr)
			return false
		}
	}
	if config.AllowLocal && utils.IsLocalIP(ip) {
		return true
	}

	// 检查位置
	if len(config.AllowCountryIds) > 0 || len(config.AllowProvinceIds) > 0 {
		var userRegion = iplibrary.Lookup(ip)
		if userRegion == nil || !userRegion.IsOk() {
			return false
		}
		if len(config.AllowCountryIds) > 0 {
			// 检查大中华区
			var found = false
			for _, countryId := range config.AllowCountryIds {
				if regionconfigs.MatchUserRegion(userRegion.CountryId(), userRegion.ProvinceId(), countryId) {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}
		if len(config.AllowProvinceIds) > 0 && !lists.ContainsInt64(config.AllowProvinceIds, userRegion.ProvinceId()) {
			return false
		}
	}

	// 检查单独允许的IP
	if len(config.AllowIPRanges()) > 0 {
		for _, r := range config.AllowIPRanges() {
			if r.Contains(ipAddr) {
				return true
			}
		}
		return false
	}

	return true
}

// 请求检查相关正则
var searchEngineRegex = regexp.MustCompile(`(?i)(60spider|adldxbot|adsbot-google|applebot|admantx|alexa|baidu|bingbot|bingpreview|facebookexternalhit|googlebot|proximic|slurp|sogou|twitterbot|yandex)`)
var spiderRegexp = regexp.MustCompile(`(?i)(python|pycurl|http-client|httpclient|apachebench|nethttp|http_request|java|perl|ruby|scrapy|php|rust|curl|wget)`) // 其中增加了curl和wget

// 检查请求
func checkRequestSecurity(securityConfig *systemconfigs.SecurityConfig, req *http.Request) bool {
	if securityConfig == nil {
		return true
	}

	var userAgent = req.UserAgent()
	var refererURL = req.Referer()
	var referHost = ""
	u, err := url.Parse(refererURL)
	if err == nil {
		referHost = u.Host
	}

	// 检查搜索引擎
	if securityConfig.DenySearchEngines && (len(userAgent) == 0 || searchEngineRegex.MatchString(userAgent) || (len(referHost) > 0 && searchEngineRegex.MatchString(referHost))) {
		return false
	}

	// 检查爬虫
	if securityConfig.DenySpiders && (len(userAgent) == 0 || spiderRegexp.MatchString(userAgent) || (len(referHost) > 0 && spiderRegexp.MatchString(referHost))) {
		return false
	}

	// 检查允许访问的域名
	if len(securityConfig.AllowDomains) > 0 {
		var domain = req.Host
		realDomain, _, err := net.SplitHostPort(domain)
		if err == nil && len(realDomain) > 0 {
			domain = realDomain
		}
		if !configutils.MatchDomains(securityConfig.AllowDomains, domain) {
			return false
		}
	}

	return true
}

// 检查是否为禁止的请求
func isEvilRequest(req *http.Request) bool {
	var headersJSON, _ = json.Marshal(req.Header)
	return bytes.Contains(headersJSON, []byte("fofa."))
}
