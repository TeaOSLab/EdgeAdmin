Tea.context(function () {
    this.createGroup = function () {
        teaweb.popup(Tea.url(".createPopup"), {
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }

    this.updateGroup = function (groupId) {
        teaweb.popup(Tea.url(".updatePopup", {groupId: groupId}), {
            callback: function () {
                teaweb.success("保存成功", function () {
                    teaweb.reload()
                })
            }
        })
    }

    this.deleteGroup = function (groupId) {
        teaweb.confirm("确定要删除此分组吗？", function () {
            this.$post(".delete")
                .params({groupId: groupId})
                .success(function () {
                    teaweb.success("删除成功", function () {
                        teaweb.reload()
                    })
                })
        })
    }
})