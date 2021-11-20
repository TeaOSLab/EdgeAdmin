Tea.context(function () {
	this.STEP_PREPARE = "prepare"
	this.STEP_DATABASE = "database"
	this.STEP_ADMIN = "admin"
	this.STEP_API = "api"
	this.STEP_ADDRESS = "address"
	this.STEP_UPGRADE = "upgrade"
	this.STEP_FINISH = "finish"

	this.step = this.STEP_PREPARE

	this.doBack = function (step) {
		this.step = step

		switch (step) {
			case this.STEP_UPGRADE:
				if (this.apiNodeChanged == 0) {
					this.doBack(this.STEP_ADMIN)
				}
		}
	}

	/**
	 * 准备工作
	 */
	this.doPrepare = function () {
		this.step = this.STEP_DATABASE
	}

	/**
	 * 数据库
	 */
	this.databaseChanged = 1
	this.databaseTransferred = 0

	this.doDatabase = function () {
		if (this.databaseChanged == 1 && this.databaseTransferred == 0) {
			teaweb.warn("请先将当前的数据导入到新的数据库中。")
			return
		}
		this.step = this.STEP_API
	}

	/**
	 * API
	 */
	this.apiNodeChanged = 1
	this.apiNodeHost = ""
	this.apiNodePort = ""
	this.apiNodeProtocol = "http"
	this.apiNodeInstalled = 1

	this.doAPI = function () {
		if (this.apiNodeChanged == 0) {
			this.step = this.STEP_ADMIN
			return
		}
		if (this.apiNodeInstalled == 0) {
			teaweb.warn("请先安装新的API节点")
			return
		}

		this.$post(".validateAPI")
			.params({
				host: this.apiNodeHost,
				port: this.apiNodePort,
				protocol: this.apiNodeProtocol
			})
			.timeout(30)
			.success(function (resp) {
				if (this.apiNodeChanged == 1) {
					this.step = this.STEP_ADDRESS
					this.apiAddressHosts = resp.data.hosts
				} else {
					this.step = this.STEP_ADMIN
				}
			})
	}

	/**
	 * 修改地址
	 */
	this.apiAddresses = []
	this.apiAddressHosts = []

	this.doAddress = function () {
		this.step = this.STEP_ADMIN
	}

	/**
	 * 管理平台
	 */
	this.adminNodeChanged = 1
	this.adminNodeInstalled = 1

	this.doAdmin = function () {
		if (this.adminNodeChanged == 1 && this.adminNodeInstalled == 0) {
			teaweb.warn("请先安装新的管理平台")
			return
		}

		if (this.apiNodeChanged == 0) {
			this.step = this.STEP_FINISH
		} else {
			this.step = this.STEP_UPGRADE
		}
	}

	/**
	 * 边缘节点
	 */
	this.isUpgrading = false
	this.percentNodes = 0
	this.countNodes = 0
	this.countFinishedNodes = 0

	this.doStartUpgrade = function () {
		this.percentNodes = 0
		this.countNodes = 0
		this.countFinishedNodes = 0

		this.$post(".statNodes")
			.success(function (resp) {
				this.countNodes = resp.data.countNodes
				if (this.countNodes == 0) {
					this.isUpgrading = true
					this.percentNodes = 100
					return
				}

				this.isUpgrading = true
				this.upgradeNodeTimer()
			})
			.fail(function () {

			})
			.error(function () {

			})
	}

	this.upgradeNodeTimer = function () {
		if (!this.isUpgrading) {
			return
		}
		if (this.percentNodes == 100) {
			return
		}
		this.$post(".upgradeNodes")
			.params({
				apiNodeProtocol: this.apiNodeProtocol,
				apiNodeHost: this.apiNodeHost,
				apiNodePort: this.apiNodePort
			})
			.success(function (resp) {
				this.countFinishedNodes += resp.data.count
				if (this.countNodes > 0) {
					this.percentNodes = this.countFinishedNodes * 100 / this.countNodes
					if (this.percentNodes > 100) {
						this.percentNodes = 100
					}
				}

				if (resp.data.hasNext) {
					this.$delay(function () {
						this.upgradeNodeTimer()
					}, 5000)
				}
			})
			.fail(function (resp) {
				this.isUpgrading = false

				teaweb.warn(resp.message)
			})
			.error(function (err) {
				teaweb.warn("请求错误：" + err.message)
				this.isUpgrading = false
			})
	}

	this.doUpgrade = function () {
		this.step = this.STEP_FINISH
	}

	/**
	 * 完成
	 */
	this.doFinish = function () {
		window.location = "/"
	}
})