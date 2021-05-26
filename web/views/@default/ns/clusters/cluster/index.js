Tea.context(function () {
    this.deleteNode = function (nodeId) {
        teaweb.confirm("确定要删除这个节点吗？", function () {
            this.$post("/ns/clusters/cluster/deleteNode")
                .params({
                    nodeId: nodeId
                })
                .refresh();
        })
    }

    this.upNode = function (nodeId) {
        teaweb.confirm("确定要手动上线此节点吗？", function () {
            this.$post("/ns/clusters/cluster/node/up")
                .params({
                    nodeId: nodeId
                })
                .refresh()
        })
    }
})