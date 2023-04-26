Tea.context(function () {
    this.tables = []
    this.isLoading = true

    this.$delay(function () {
        this.reload()
    })

    this.reload = function () {
        this.isLoading = true
        this.$post("$")
			.params({
				orderTable: this.orderTable,
				orderSize: this.orderSize
			})
            .success(function (resp) {
                this.tables = resp.data.tables;
            })
            .done(function () {
                this.isLoading = false
            })
    }

    this.deleteTable = function (tableName) {
        let that = this
        teaweb.confirm("html:确定要删除此数据表吗？<br/>删除后数据不能恢复！", function () {
            that.$post(".deleteTable")
                .params({
                    table: tableName
                })
                .success(function () {
                    teaweb.success("操作成功", function () {
                        that.reload()
                    })
                })
        })
    }

    this.truncateTable = function (tableName) {
        let that = this
        teaweb.confirm("html:确定要清空此数据表吗？<br/>清空后数据不能恢复！", function () {
            that.$post(".truncateTable")
                .params({
                    table: tableName
                })
                .success(function () {
                    teaweb.success("操作成功", function () {
                        that.reload()
                    })
                })
        })
    }
})