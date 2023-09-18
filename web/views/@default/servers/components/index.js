Tea.context(function () {
    this.tab = "tcpPorts"

    this.$delay(function () {
        if (window.location.hash != null && window.location.hash.length > 1) {
            this.selectTab(window.location.hash.substring(1))
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