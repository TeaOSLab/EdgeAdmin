Tea.context(function () {
    this.$delay(function () {
        this.reload()
    })

    this.reload = function () {
        this.$post("$")
            .success(function (resp) {
                this.tasks = resp.data.tasks
            })
            .done(function () {
                this.$delay(function () {
                    this.reload()
                }, 3000)
            })
    }

    this.deleteTask = function (taskId) {
        let that = this
        teaweb.confirm("确定要删除这个任务吗？", function () {
            that.$post(".delete")
                .params({
                    taskId: taskId
                })
                .success(function () {
                    teaweb.reload()
                })
        })
    }
})