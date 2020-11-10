Tea.context(function () {
	this.moreOptionsVisible = false
	this.globalChangedClusters = []
	this.globalMessageBadge = 0

	if (typeof this.leftMenuItemIsDisabled == "undefined") {
		this.leftMenuItemIsDisabled = false
	}

	this.$delay(function () {
		if (this.$refs.focus != null) {
			this.$refs.focus.focus()
		}

		// 检查变更
		this.checkClusterChanges()

		// 检查消息
		this.checkMessages()
	})

	/**
	 * 左侧子菜单
	 */
	this.showSubMenu = function (menu) {
		if (menu.alwaysActive) {
			return
		}
		if (this.teaSubMenus.menus != null && this.teaSubMenus.menus.length > 0) {
			this.teaSubMenus.menus.$each(function (k, v) {
				if (menu.id == v.id) {
					return
				}
				v.isActive = false
			})
		}
		menu.isActive = !menu.isActive
	};

	/**
	 * 检查集群变更
	 */
	this.checkClusterChanges = function () {
		this.$post("/clusters/checkChange")
			.params({
				isNotifying: (this.globalChangedClusters.length > 0) ? 1 : 0
			})
			.timeout(60)
			.success(function (resp) {
				this.globalChangedClusters = resp.data.clusters;
			})
			.fail(function () {
				this.globalChangedClusters = [];
			})
			.done(function () {
				let delay = 3000
				if (this.globalChangedClusters.length > 0) {
					delay = 30000
				}
				this.$delay(function () {
					this.checkClusterChanges()
				}, delay)
			})
	}

	/**
	 * 检查消息
	 */
	this.checkMessages = function () {
		this.$post("/messages/badge")
			.params({})
			.success(function (resp) {
				this.globalMessageBadge = resp.data.count
			})
			.done(function () {
				let delay = 6000
				if (this.globalMessageBadge > 0) {
					delay = 30000
				}
				this.$delay(function () {
					this.checkMessages()
				}, delay)
			})
	}

	/**
	 * 同步集群配置
	 */
	this.syncClustersConfigs = function () {
		teaweb.confirm("html:有若干个集群配置已变更！<br/>确定要同步配置到边缘节点吗？", function () {
			this.$post("/clusters/sync")
				.success(function () {
					this.globalChangedClusters = [];
				})
		})
	};

	/**
	 * 底部伸展框
	 */
	this.showQQGroupQrcode = function () {
		teaweb.popup("/about/qq", {
			width: "21em",
			height: "24em"
		})
	}
});

window.NotifySuccess = function (message, url, params) {
	if (typeof (url) == "string" && url.length > 0) {
		if (url[0] != "/") {
			url = Tea.url(url, params);
		}
	}
	return function () {
		teaweb.success(message, function () {
			window.location = url;
		});
	};
};

window.NotifyReloadSuccess = function (message) {
	return function () {
		teaweb.success(message, function () {
			window.location.reload()
		})
	}
}

window.NotifyDelete = function (message, url, params) {
	teaweb.confirm(message, function () {
		Tea.Vue.$post(url)
			.params(params)
			.refresh();
	});
};

window.NotifyPopup = function (resp) {
	window.parent.teaweb.popupFinish(resp);
};

window.ChangePageSize = function (size) {
	let url = window.location.toString();
	if (url.indexOf("pageSize") > 0) {
		url = url.replace(/pageSize=\d+/g, "pageSize=" + size);
	} else {
		if (url.indexOf("?") > 0) {
			url += "&pageSize=" + size;
		} else {
			url += "?pageSize=" + size;
		}
	}
	window.location = url;
};