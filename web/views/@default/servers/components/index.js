Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	/**
	 * 域名不匹配动作
	 */
	this.domainMismatchAction = "page"

	this.domainMismatchActionPageOptions = {
		statusCode: 404,
		contentHTML: `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<title>404 not found</title>
</head>
<body>
<h4>找不到您要访问的页面.</h4>
<h4>Sorry, page not found.</h4>
</body>
</html>
`
	}
	if (this.globalConfig.httpAll.domainMismatchAction != null) {
		this.domainMismatchAction = this.globalConfig.httpAll.domainMismatchAction.code

		if (this.domainMismatchAction == "page") {
			this.domainMismatchActionPageOptions = this.globalConfig.httpAll.domainMismatchAction.options;
		}
	}

})