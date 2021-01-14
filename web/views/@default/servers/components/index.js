Tea.context(function () {
    this.tab = "domainMatch"

    this.$delay(function () {
        if (window.location.hash != null && window.location.hash.length > 1) {
            this.selectTab(window.location.hash.substr(1))
        }
    })

    this.selectTab = function (tab) {
        this.tab = tab
        window.location.hash = "#" + tab
    }

    this.success = function () {
        teaweb.success("保存成功", function () {
            teaweb.reload()
        })
    }

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

    /**
     * TCP端口
     */
    this.tcpAllPortRangeMin = 10000
    this.tcpAllPortRangeMax = 40000
    if (this.globalConfig.tcpAll.portRangeMin > 0) {
        this.tcpAllPortRangeMin = this.globalConfig.tcpAll.portRangeMin
    }
    if (this.globalConfig.tcpAll.portRangeMax > 0) {
        this.tcpAllPortRangeMax = this.globalConfig.tcpAll.portRangeMax
    }

    this.tcpAllDenyPorts = []
    if (this.globalConfig.tcpAll.denyPorts != null) {
        this.tcpAllDenyPorts = this.globalConfig.tcpAll.denyPorts
    }
})