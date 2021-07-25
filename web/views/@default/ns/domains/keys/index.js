Tea.context(function () {
	this.createKey = function () {
		teaweb.popup(Tea.url(".createPopup?domainId=" + this.domain.id), {
			height: "24em",
			callback: function () {
				teaweb.successRefresh("保存成功")
			}
		})
	}

	this.updateKey = function (keyId) {
		teaweb.popup(Tea.url(".updatePopup?keyId=" + keyId), {
			height: "27em",
			callback: function () {
				teaweb.successRefresh("保存成功")
			}
		})
	}

	this.deleteKey = function (keyId) {
		teaweb.confirm("确定要删除这个密钥吗？", function () {
			this.$post(".delete")
				.params({
					keyId: keyId
				})
				.success(function () {
					teaweb.successRefresh("删除成功")
				})
		})
	}
})