Tea.context(function () {
	this.updateProvider = function (providerId) {
		teaweb.popup(Tea.url(".updatePopup?providerId=" + providerId), {
			height: "26em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.createDomain = function () {
		teaweb.popup("/dns/domains/createPopup?providerId=" + this.provider.id, {
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteDomain = function (domain) {
		let that = this
		teaweb.confirm("确定要删除域名\"" + domain.name + "\"吗？", function () {
			that.$post("/dns/domains/delete")
				.params({
					domainId: domain.id
				})
				.post()
				.refresh()
		})
	}
})