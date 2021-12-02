Tea.context(function () {
	this.success = NotifyPopup

	this.bodyType = this.pageConfig.bodyType

	this.newStatus = ""
	if (this.pageConfig.newStatus > 0) {
		this.newStatus = this.pageConfig.newStatus
	}

	this.addHTMLTemplate = function () {
		this.$refs.htmlBody.value = `<!DOCTYPE html>
<html>
<head>
\t<title>\${status} \${statusMessage}</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
</head>
<body>

<h1>\${status} \${statusMessage}</h1>
<p><!-- 内容 --></p>

<address>Request ID: \${requestId}, Powered by GoEdge.</address>

</body>
</html>`
	}
})