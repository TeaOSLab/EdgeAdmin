Tea.context(function () {
	this.success = NotifyPopup

	this.rewriteRule = {
		mode: "proxy",
		redirectStatus: 307
	}
	this.statusOptions = [
		{"code": 301, "text": "Moved Permanently"},
		{"code": 308, "text": "Permanent Redirect"},
		{"code": 302, "text": "Found"},
		{"code": 303, "text": "See Other"},
		{"code": 307, "text": "Temporary Redirect"}
	]
})