Tea.context(function () {
    this.deleteCluster = function (clusterId) {
        let that = this
        teaweb.confirm("html:确定要删除此集群吗？<br/>删除后不可恢复！", function () {
            that.$post("/clusters/cluster/delete")
                .params({
                    clusterId: clusterId
                })
                .success(function () {
                    teaweb.success("删除成功", function () {
                        window.location = "/clusters"
                    })
                })
        })
    }
})