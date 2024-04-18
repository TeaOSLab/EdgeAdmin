Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	// Menu
	this.currentItem = ""

	this.titleMenus = []
	this.$delay(function () {
		let elements = document.querySelectorAll("#global-config-form h4")
		this.currentItem = elements[0].getAttribute("id")
		for (let i = 0; i < elements.length; i++) {
			let textContent = elements[i].textContent
			if (textContent == null || textContent.length == 0) {
				textContent = elements[i].innerText
			}
			let itemId = elements[i].getAttribute("id")
			if (window.location.hash == "#" + itemId) {
				this.currentItem = itemId
			}
			this.titleMenus.push({
				name: textContent,
				id: itemId
			})
		}
	})

	this.selectItem = function (item) {
		this.currentItem = item.id
	}

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