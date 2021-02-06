Tea.context(function () {
    this.createAction = function () {
        teaweb.popup(Tea.url(".createPopup", {clusterId: this.clusterId}), {
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }

    this.updateAction = function (actionId) {
        teaweb.popup(Tea.url(".updatePopup", {actionId: actionId}), {
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }

    this.deleteAction = function (actionId) {
        let that = this
        teaweb.confirm("确定要删除此动作吗？", function () {
            that.$post(".delete")
                .params({
                    actionId: actionId
                })
                .success(function () {
                   teaweb.success("删除成功", function () {
                       teaweb.reload()
                   })
                })
        })
    }
})