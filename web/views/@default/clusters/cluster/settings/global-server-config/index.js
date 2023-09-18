Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	/**
	 * TCP端口
	 */
	this.tcpAllPortRangeMin = 10000
	this.tcpAllPortRangeMax = 40000
	if (this.config.tcpAll.portRangeMin > 0) {
		this.tcpAllPortRangeMin = this.config.tcpAll.portRangeMin
	}
	if (this.config.tcpAll.portRangeMax > 0) {
		this.tcpAllPortRangeMax = this.config.tcpAll.portRangeMax
	}

	this.tcpAllDenyPorts = []
	if (this.config.tcpAll.denyPorts != null) {
		this.tcpAllDenyPorts = this.config.tcpAll.denyPorts
	}
})