Tea.context(function () {
	this.isInstalling = false
	let installingNode = null

	this.$delay(function () {
		this.reload()
	})

	this.installNode = function (node) {
		let that = this
		teaweb.confirm("确定要开始安装此节点吗？", function () {
			installingNode = node
			that.isInstalling = true
			node.isInstalling = true

			that.$post("$")
				.params({
					nodeId: node.id
				})
		})
	}

	this.reload = function () {
		let that = this
		if (installingNode != null) {
			this.$post("/clusters/cluster/installStatus")
				.params({
					nodeId: installingNode.id
				})
				.success(function (resp) {
					if (resp.data.status != null) {
						installingNode.installStatus = resp.data.status
						if (installingNode.installStatus.isFinished) {
							if (installingNode.installStatus.isOk) {
								installingNode = null
								teaweb.success("安装成功", function () {
									window.location.reload()
								})
							} else {
								let nodeId = installingNode.id
								let errMsg = installingNode.installStatus.error
								that.isInstalling = false
								installingNode.isInstalling = false
								installingNode = null

								switch (resp.data.status.errorCode) {
									case "EMPTY_LOGIN":
									case "EMPTY_SSH_HOST":
									case "EMPTY_SSH_PORT":
									case "EMPTY_GRANT":
										teaweb.warn("需要填写SSH登录信息", function () {
											teaweb.popup("/clusters/cluster/updateNodeSSH?nodeId=" + nodeId, {
												callback: function () {
													teaweb.reload()
												}
											})
										})
										return
									default:
										teaweb.warn("安装失败：" + errMsg)
								}
							}
						}
					}
				})
				.done(function () {
					setTimeout(this.reload, 3000)
				})
		} else {
			setTimeout(this.reload, 3000)
		}
	}
})