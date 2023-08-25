Tea.context(function () {
	this.success = NotifyPopup

	this.bodyType = this.pageConfig.bodyType

	this.newStatus = ""
	if (this.pageConfig.newStatus > 0) {
		this.newStatus = this.pageConfig.newStatus
	}

	this.addHTMLTemplate = function () {
		this.$refs.htmlBody.value = `<!DOCTYPE html>
<html lang="en">
<head>
\t<title>\${status} \${statusMessage}</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
\t<style>
\t\taddress { line-height: 1.8; }
\t</style>
</head>
<body>

<h1>\${status} \${statusMessage}</h1>
<p><!-- 内容 --></p>

<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>

</body>
</html>`
	}
})