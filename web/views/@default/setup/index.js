Tea.context(function () {
	this.STEP_INTRO = "intro"
	this.STEP_API = "api"
	this.STEP_DB = "db"
	this.STEP_ADMIN = "admin"
	this.STEP_FINISH = "finish"

	this.step = this.STEP_INTRO

	this.$delay(function () {
		this.loadStatusText()
	})

	// 介绍
	this.goIntroNext = function () {
		this.step = this.STEP_API
	}

	// API节点
	this.apiNodeInfo = {}
	this.apiNodeMode = "new"
	this.newAPINodePort = "8001"
	this.apiRequesting = false

	this.apiHostInput = false // 是否手工输入

	this.apiSubmit = function () {
		this.apiRequesting = true
	}

	this.apiDone = function () {
		this.apiRequesting = false
	}

	this.apiSuccess = function (resp) {
		this.step = this.STEP_DB
		this.detectDB()
		this.apiNodeInfo = resp.data.apiNode

		if (this.apiNodeMode == "new") {
			this.$delay(function () {
				this.$refs.dbHost.focus()
			}, 200)
		}
	}

	this.goBackIntro = function () {
		this.step = this.STEP_INTRO
	}

	this.inputAPIHost = function () {
		this.apiHostInput = true
		this.$delay(function () {
			this.$refs.newHostRef.focus()
		})
	}

	// 数据库
	this.dbInfo = {}
	this.localDB = {"host": "", "port": "", "username": "", "port": "", "isLocal": true, "canInstall": false}
	this.localDBHost = ""
	this.dbRequesting = false

	this.detectDB = function () {
		this.$post(".detectDB")
			.success(function (resp) {
				this.localDB = resp.data.localDB
				this.localDB["isLocal"] = true
				this.localDBHost = this.localDB.host
			})
	}

	this.checkDBIP = function () {
		this.localDB["isLocal"] = true
		if (this.localDB.host.length == 0) {
			return
		}
		this.$post(".checkLocalIP")
			.params({
				host: this.localDB.host
			})
			.success(function (resp) {
				this.localDB["isLocal"] = resp.data.isLocal
				this.$forceUpdate()
			})
	}

	this.dbSubmit = function () {
		this.dbRequesting = true
	}

	this.dbSuccess = function (resp) {
		this.step = this.STEP_ADMIN
		this.dbInfo = resp.data.db

		this.$delay(function () {
			if (this.$refs.adminPasswordInput != null) {
				this.$refs.adminPasswordInput.focus()
			}
		})
	}

	this.dbDone = function () {
		this.dbRequesting = false
	}

	this.goBackAPI = function () {
		this.step = this.STEP_API
	}

	this.goDBNext = function () {
		this.step = this.STEP_ADMIN
	}

	// 管理员
	this.goBackDB = function () {
		this.step = this.STEP_DB
	}

	this.adminInfo = {}
	this.adminPassword = ""
	this.adminPassword2 = ""
	this.adminPasswordVisible = false

	this.showAdminPassword = function () {
		this.adminPasswordVisible = !this.adminPasswordVisible

		// TODO 切换密码显示的时候应该focus输入框
	}

	this.adminSuccess = function (resp) {
		this.step = this.STEP_FINISH
		this.adminInfo = resp.data.admin
	}

	// 结束
	this.goBackAdmin = function () {
		this.step = this.STEP_ADMIN
	}

	this.isInstalling = false
	this.finishSubmit = function () {
		this.isInstalling = true
	}

	this.finishDone = function () {
		this.isInstalling = false
	}

	this.finishSuccess = function () {
		teaweb.success("html:恭喜你！安装完成！<br/>请记住你创建的管理员账号，现在跳转到登录界面。", function () {
			window.location = "/"
		})
	}

	this.statusText = ""
	this.loadStatusText = function () {
		if (!this.isInstalling) {
			this.statusText = ""
			this.$delay(function () {
				this.loadStatusText()
			}, 1000)
			return
		}
		this.$post(".status")
			.success(function (resp) {
				this.statusText = resp.data.statusText
			})
			.done(function () {
				this.$delay(function () {
					this.loadStatusText()
				}, 1000)
			})
	}

	/**
	 * MySQL
	 */
	this.installMySQL = function () {
		let that = this
		teaweb.popup("/setup/mysql/installPopup", {
			height: "28em",
			onClose: function () {
				that.detectDB()
			}
		})
	}
})