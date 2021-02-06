Tea.context(function () {
    this.updateAllRead = function () {
        let that = this
        teaweb.confirm("确定要设置所有的未读消息为已读吗？", function () {
            that.$post("/messages/readAll")
                .success(function () {
                    // 刷新父级页面Badge
                    if (window.parent.Tea != null && window.parent.Tea.Vue != null) {
                        window.parent.Tea.Vue.checkMessagesOnce()
                    }

                    window.location = "/messages"
                })
        })
    }

    this.updatePageRead = function () {
        let that = this
        teaweb.confirm("确定要设置当前页的未读消息为已读吗？", function () {
            let messageIds = []
            that.messages.forEach(function (v) {
                messageIds.push(v.id)
            })
            that.$post("/messages/readPage")
                .params({
                    messageIds: messageIds
                })
                .success(function () {
                    // 刷新父级页面Badge
                    if (window.parent.Tea != null && window.parent.Tea.Vue != null) {
                        window.parent.Tea.Vue.checkMessagesOnce()
                    }

                    teaweb.reload()
                })
        })
    }
})