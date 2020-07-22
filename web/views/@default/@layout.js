Tea.context(function () {
	this.moreOptionsVisible = false;

	this.$delay(function () {
		if (this.$refs.focus != null) {
			this.$refs.focus.focus();
		}
	});

	/**
	 * 左侧子菜单
	 */
	this.showSubMenu = function (menu) {
		if (menu.alwaysActive) {
			return;
		}
		if (this.teaSubMenus.menus != null && this.teaSubMenus.menus.length > 0) {
			this.teaSubMenus.menus.$each(function (k, v) {
				if (menu.id == v.id) {
					return;
				}
				v.isActive = false;
			});
		}
		menu.isActive = !menu.isActive;
	};
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