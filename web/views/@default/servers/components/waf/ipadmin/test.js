Tea.context(function () {
    this.ip = ""
    this.result = {
        isDone: false,
        isOk: false,
        isFound: false,
        isAllowed: false,
        error: "",
        province: null,
        country: null,
        ipItem: null,
        ipList: null
    }

    this.$delay(function () {
        this.$watch("ip", function () {
            this.result.isDone = false
        })
    })

    this.success = function (resp) {
        this.result = resp.data.result
    }

    this.updateItem = function (itemId) {
        teaweb.popup(Tea.url(".updateIPPopup?firewallPolicyId=" + this.firewallPolicyId, {itemId: itemId}), {
            height: "23em",
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }
})