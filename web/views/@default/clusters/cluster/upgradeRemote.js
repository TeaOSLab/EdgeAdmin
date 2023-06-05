Tea.context(function () {
    this.isInstalling = false
    this.isBatch = false
    let installingNode = null

    this.nodes.forEach(function (v) {
        v.isChecked = false
    })

    this.$delay(function () {
        this.reload()
    })

    let that = this

    this.checkNodes = function (isChecked) {
        this.nodes.forEach(function (v) {
            v.isChecked = isChecked
        })
    }

    this.countCheckedNodes = function () {
        return that.nodes.$count(function (k, v) {
            return v.isChecked
        })
    }

    this.installNode = function (node) {
        let that = this
        if (this.isBatch) {
            installingNode = node
            that.isInstalling = true
            node.isInstalling = true

            that.$post("$")
                .params({
                    nodeId: node.id
                })
        } else {
            teaweb.confirm("确定要开始升级此节点吗？", function () {
                installingNode = node
                that.isInstalling = true
                node.isInstalling = true

                that.$post("$")
                    .params({
                        nodeId: node.id
                    })
            })
        }
    }

    this.installBatch = function () {
        let that = this
        this.isBatch = true
        teaweb.confirm("确定要批量升级选中的节点吗？", function () {
            that.installNext()
        })
    }

    /**
     * 安装下一个
     */
    this.installNext = function () {
        let nextNode = this.nodes.$find(function (k, v) {
            return v.isChecked
        })

        if (nextNode == null) {
            teaweb.success("全部升级成功", function () {
                teaweb.reload()
            })
        } else {
            this.installNode(nextNode)
        }
        return
    }

    /**
     * 重新加载状态
     */
    this.reload = function () {
        let that = this
        if (installingNode != null) {
            this.$post("/clusters/cluster/upgradeStatus")
                .params({
                    nodeId: installingNode.id
                })
                .success(function (resp) {
                    if (resp.data.status != null) {
                        installingNode.installStatus = resp.data.status
                        if (installingNode.installStatus.isFinished) {
                            if (installingNode.installStatus.isOk) {
                                installingNode.isChecked = false // 取消选中
                                installingNode = null
                                if (that.isBatch) {
                                    that.installNext()
                                } else {
                                    teaweb.success("升级成功", function () {
                                        teaweb.reload()
                                    })
                                }
                            } else {
                                let nodeId = installingNode.id
                                let errMsg = installingNode.installStatus.error
                                that.isInstalling = false
                                installingNode.isInstalling = false
                                installingNode = null

                                switch (resp.data.status.errorCode) {
                                    case "EMPTY_LOGIN":
                                    case "EMPTY_SSH_HOST":
                                    case "EMPTY_SSH_PORT":
                                    case "EMPTY_GRANT":
                                        teaweb.warn("需要填写SSH登录信息", function () {
                                            teaweb.popup("/clusters/cluster/updateNodeSSH?nodeId=" + nodeId, {
												height: "30em",
                                                callback: function () {
                                                    teaweb.reload()
                                                }
                                            })
                                        })
                                        return
                                    case "CREATE_ROOT_DIRECTORY_FAILED":
                                        teaweb.warn("创建根目录失败，请检查目录权限或者手工创建：" + errMsg)
                                        return
                                    case "INSTALL_HELPER_FAILED":
                                        teaweb.warn("安装助手失败：" + errMsg)
                                        return
                                    case "TEST_FAILED":
                                        teaweb.warn("环境测试失败：" + errMsg)
                                        return
                                    case "RPC_TEST_FAILED":
                                        teaweb.confirm("html:要升级的节点到API服务之间的RPC通讯测试失败，具体错误：" + errMsg + "，<br/>现在修改API信息？", function () {
                                            window.location = "/settings/api"
                                        })
                                        return
                                    default:
                                        teaweb.warn("升级失败：" + errMsg)
                                }
                            }
                        }
                    }
                })
                .done(function () {
                    setTimeout(this.reload, 3000)
                })
        } else {
            setTimeout(this.reload, 3000)
        }
    }
})