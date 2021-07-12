Tea.context(function () {
    this.moreOptionsVisible = false
    this.globalMessageBadge = 0

    if (typeof this.leftMenuItemIsDisabled == "undefined") {
        this.leftMenuItemIsDisabled = false
    }

    this.$delay(function () {
        if (this.$refs.focus != null) {
            this.$refs.focus.focus()
        }

        // 检查消息
        this.checkMessages()

        // 检查集群节点同步
        this.loadNodeTasks();

        // 检查DNS同步
        this.loadDNSTasks()
    })

	/**
	 * 切换模板
	 */
	this.changeTheme = function () {
		this.$post("/ui/theme")
			.success(function (resp) {
				teaweb.successToast("界面风格已切换")
				this.teaTheme = resp.data.theme
			})
	}

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

    this.checkMessagesOnce = function () {
        this.$post("/messages/badge")
            .params({})
            .success(function (resp) {
                this.globalMessageBadge = resp.data.count
            })
    }

    this.showMessages = function () {
        teaweb.popup("/messages", {
            height: "24em",
            width: "50em"
        })
    }

    /**
     * 底部伸展框
     */
    this.showQQGroupQrcode = function () {
        teaweb.popup("/about/qq", {
            width: "21em",
            height: "24em"
        })
    }

    /**
     * 弹窗中默认成功回调
     */
    if (window.IS_POPUP === true) {
        this.success = window.NotifyPopup
    }

    /**
     * 节点同步任务
     */
    this.doingNodeTasks = {
        isDoing: false,
        hasError: false,
        isUpdated: false
    }

    this.loadNodeTasks = function () {
        if (!Tea.Vue.teaCheckNodeTasks) {
            return
        }
        this.$post("/clusters/tasks/check")
            .success(function (resp) {
                this.doingNodeTasks.isDoing = resp.data.isDoing
                this.doingNodeTasks.hasError = resp.data.hasError
                this.doingNodeTasks.isUpdated = true
            })
            .done(function () {
                this.$delay(function () {
                    this.loadNodeTasks()
                }, 3000)
            })
    }

    this.showNodeTasks = function () {
        teaweb.popup("/clusters/tasks/listPopup", {
            height: "24em",
            width: "50em"
        })
    }

    /**
     * DNS同步任务
     */
    this.doingDNSTasks = {
        isDoing: false,
        hasError: false,
        isUpdated: false
    }

    this.loadDNSTasks = function () {
        if (!Tea.Vue.teaCheckDNSTasks) {
            return
        }
        this.$post("/dns/tasks/check")
            .success(function (resp) {
                this.doingDNSTasks.isDoing = resp.data.isDoing
                this.doingDNSTasks.hasError = resp.data.hasError
                this.doingDNSTasks.isUpdated = true
            })
            .done(function () {
                this.$delay(function () {
                    this.loadDNSTasks()
                }, 3000)
            })
    }

    this.showDNSTasks = function () {
        teaweb.popup("/dns/tasks/listPopup", {
            height: "24em",
            width: "50em"
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