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

    this.updateItem = function (listId, itemId) {
        teaweb.popup(Tea.url(".updateIPPopup?listId=" + listId, {itemId: itemId}), {
            height: "30em",
            callback: function () {
                teaweb.success("保存成功", function () {

                })
            }
        })
    }
})