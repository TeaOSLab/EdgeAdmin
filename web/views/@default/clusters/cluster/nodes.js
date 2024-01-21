Tea.context(function () {
	this.teaweb = teaweb

	// 显示的统计项
	this.windowWidth = window.innerWidth
	this.miniWidth = 760
	this.columnWidth1 = 800
	this.columnWidth2 = 900
	this.columnWidth3 = 1000
	this.columnWidth4 = 1100
	this.columnWidth5 = 1200

	let that = this
	window.addEventListener("resize", function () {
		that.windowWidth = window.innerWidth
	})

	this.deleteNode = function (nodeId) {
		teaweb.confirm("确定要从当前集群中删除这个节点吗？", function () {
			this.$post("/nodes/delete")
				.params({
					clusterId: this.clusterId,
					nodeId: nodeId
				})
				.refresh();
		})
	}

	this.upNode = function (nodeId) {
		teaweb.confirm("确定要手动上线此节点吗？", function () {
			this.$post("/clusters/cluster/node/up")
				.params({
					nodeId: nodeId
				})
				.refresh()
		})
	}

	this.updateNodeDNS = function (nodeId) {
		let that = this
		teaweb.popup("/clusters/cluster/node/updateDNSPopup?clusterId=" + this.clusterId + "&nodeId=" + nodeId, {
			width: "46em",
			height: "26em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateNodeOn = function (nodeId, isOn) {
		let that = this
		let op
		if (isOn) {
			op = "启用"
		} else {
			op = "停用"
		}
		teaweb.confirm("确定要" + op + "此节点吗？", function () {
			that.$post(".node.updateIsOn")
				.params({
					nodeId: nodeId,
					isOn: isOn
				})
				.success(function () {
					teaweb.successRefresh(op + "成功")
				})
		})
	}

	/**
	 * 显示和隐藏IP
	 */
	this.mostIPVisible = 4

	this.showMoreIP = function (nodeIndex, node) {
		if (typeof node.ipAddressesVisible != "boolean") {
			node.ipAddressesVisible = false
		}
		node.ipAddressesVisible = !node.ipAddressesVisible
		Vue.set(this.nodes, nodeIndex, node)
	}
})