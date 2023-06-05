Tea.context(function () {
	this.nodeId = 0
	this.node = {}
	this.sshHost = ""
	this.sshPort = ""
	this.grantId = 0
	this.step = "info"
	this.name = ""

	this.success = function (resp) {
		this.node = resp.data.node
		this.nodeId = this.node.id
		this.sshHost = this.node.login.params.host
		if (this.node.login.params.port > 0) {
			this.sshPort = this.node.login.params.port
		}
		if (this.node.addresses.length > 0) {
			this.sshHost = this.node.addresses[0]
		}
		this.step = "install"
		this.$delay(function () {
			this.$refs.installSSHHostRef.focus()
			this.$refs.nodeLoginSuggestPortsRef.reload(this.sshHost)
		})
	}

	/**
	 * 安装
	 */
	this.isInstalled = false
	this.installMethod = "remote" // remote | manual
	this.isInstalling = false

	this.switchInstallMethod = function (method) {
		this.installMethod = method
	}

	this.selectSSHHost = function (host) {
		this.sshHost = host
		this.changeSSHHost()
	}

	this.changeSSHHost = function () {
		if (this.$refs.nodeLoginSuggestPortsRef != null) {
			this.$refs.nodeLoginSuggestPortsRef.reload(this.sshHost)
		}
	}

	this.selectLoginPort = function (port) {
		this.sshPort = port
	}

	this.autoSelectLoginPort = function (port) {
		if (this.sshPort == null || this.sshPort <= 0) {
			this.sshPort = port
		}
	}

	this.install = function () {
		if (this.node.grant != null) {
			this.grantId = this.node.grant.id
		}

		this.isInstalling = true
		this.$post(".createNodeInstall")
			.params({
				nodeId: this.node.id,
				sshHost: this.sshHost,
				sshPort: this.sshPort,
				grantId: this.grantId
			})
			.timeout(30)
			.success(function () {
				this.$delay(function () {
					this.isInstalling = true
					this.reloadStatus(this.node.id)
				})
			})
			.done(function () {
				this.isInstalling = false
			})
	}

	this.changeGrant = function (grant) {
		if (grant != null) {
			this.grantId = grant.id
		} else {
			this.grantId = 0
		}
	}

	// 刷新状态
	this.installStatus = null
	this.reloadStatus = function (nodeId) {
		let that = this

		this.$post("/clusters/cluster/node/status")
			.params({
				nodeId: nodeId
			})
			.success(function (resp) {
				this.installStatus = resp.data.installStatus
				this.node.isInstalled = resp.data.isInstalled
				if (this.node.isInstalled) {
					this.isInstalling = false
					this.isInstalled = true
					this.finish()
				}

				if (!this.isInstalling) {
					return
				}

				let nodeId = this.node.id
				let errMsg = this.installStatus.error

				if (this.installStatus.errorCode.length > 0 || errMsg.length > 0) {
					this.isInstalling = false
				}

				switch (this.installStatus.errorCode) {
					case "EMPTY_LOGIN":
					case "EMPTY_SSH_HOST":
					case "EMPTY_SSH_PORT":
					case "EMPTY_GRANT":
						teaweb.warn("需要填写SSH登录信息", function () {
							teaweb.popup("/clusters/cluster/updateNodeSSH?nodeId=" + nodeId, {
								height: "30em",
								callback: function () {
									that.install()
								}
							})
						})
						return
					case "SSH_LOGIN_FAILED":
						teaweb.warn("SSH登录失败，请检查设置")
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
						teaweb.confirm("html:要安装的节点到API服务之间的RPC通讯测试失败，具体错误：" + errMsg + "，<br/>现在修改API信息？", function () {
							window.location = "/settings/api"
						})
						return
					default:
						shouldReload = true
					//teaweb.warn("安装失败：" + errMsg)
				}
			})
			.done(function () {
				this.$delay(function () {
					this.reloadStatus(nodeId)
				}, 1000)
			});
	}

	/**
	 * 完成
	 */
	this.finish = function () {
		this.step = "finish"
	}

	this.createNext = function () {
		teaweb.reload()
	}

	this.defaultIP = ""
	this.changeName = function () {
		let matchIP = this.name.match(/(\d{1,3}\.){3}\d{1,3}/)
		if (matchIP != null) {
			if (this.validateIP(matchIP[0])) {
				this.defaultIP = matchIP[0]
			} else {
				this.defaultIP = ""
			}
		}
	}

	this.validateIP = function (v) {
		// 目前只支持ipv4
		let pieces = v.split(".")
		if (pieces.length != 4) {
			return false
		}
		for (let i = 0; i < pieces.length; i++) {
			if (!/^\d{1,3}$/.test(pieces[i])) {
				return false
			}
			let p = parseInt(pieces[i], 10)
			if (p > 255) {
				return false
			}
		}

		return true
	}
})