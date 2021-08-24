Tea.context(function () {
    this.createRecipient = function () {
        teaweb.popup(Tea.url(".createPopup"), {
            height: "27em",
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }

    this.deleteRecipient = function (recipientId) {
        teaweb.confirm("确定要删除此接收媒介吗？", function () {
            this.$post(".delete")
                .params({recipientId: recipientId})
                .success(function () {
                    teaweb.success("删除成功", function () {
                        teaweb.reload()
                    })
                })
        })
    }
})