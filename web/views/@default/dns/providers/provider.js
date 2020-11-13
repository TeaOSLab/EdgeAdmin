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

	this.updateDomain = function (domainId) {
		teaweb.popup("/dns/domains/updatePopup?domainId=" + domainId, {
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

	this.syncDomain = function (index, domain) {
		domain.isSyncing = true
		Vue.set(this.domains, index, domain)

		this.$post("/dns/domains/sync")
			.params({
				domainId: domain.id
			})
			.success(function () {
				teaweb.success("同步成功", function () {
					teaweb.reload()
				})
			})
			.fail(function (resp) {
				teaweb.warn(resp.message, function () {
					if (resp.data.shouldFix) {
						window.location = "/dns/issues"
					}
				})
			})
			.done(function () {
				domain.isSyncing = false
				Vue.set(this.domains, index, domain)
			})
	}

	this.showRoutes = function (domainId) {
		teaweb.popup("/dns/domains/routesPopup?domainId=" + domainId)
	}
})