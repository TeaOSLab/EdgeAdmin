Tea.context(function () {
    let checkedAll = false

    this.$delay(function () {
        this.reload()
    })

    this.checkAll = function (b) {
        checkedAll = b
        let that = this
        this.clusters.forEach(function (cluster, index) {
            cluster.tasks.forEach(function (task) {
                task.isChecked = checkedAll
            })
            Vue.set(that.clusters, index, cluster)
        })
    }

    this.checkTask = function (b) {
        this.clusters.forEach(function (cluster, index) {
            Vue.set(that.clusters, index, cluster)
        })
    }

    let that = this
    this.countCheckedTasks = function () {
        let count = 0
        that.clusters.forEach(function (cluster) {
            cluster.tasks.forEach(function (task) {
                if (task.isChecked) {
                    count++
                }
            })
        })
        return count
    }

    this.reload = function () {
        this.$post("$")
            .success(function (resp) {
                this.countTasks = resp.data.countTasks
                this.clusters = resp.data.clusters
            })
            .done(function () {
                this.$delay(function () {
                    // 没有选中任务的时候才重新刷新
                    if (this.countCheckedTasks() == 0) {
                        this.reload()
                    }
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

    this.deleteBatch = function () {
        let taskIds = []
        this.clusters.forEach(function (cluster) {
            cluster.tasks.forEach(function (task) {
                if (task.isChecked) {
                    taskIds.push(task.id)
                }
            })
        })

        let that = this
        teaweb.confirm("确定要批量删除选中的任务吗？", function () {
            that.$post(".deleteBatch")
                .params({
                    taskIds: taskIds
                })
                .success(function () {
                    teaweb.reload()
                })
        })
    }

	this.deleteAllTasks = function () {
		let that = this
		teaweb.confirm("确定要清空所有的任务吗？", function () {
			that.$post(".deleteAll")
				.success(function () {
					teaweb.reload()
				})
		})
	}
})