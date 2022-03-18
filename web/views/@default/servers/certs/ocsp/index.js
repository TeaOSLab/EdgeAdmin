Tea.context(function () {
	this.certIds = []
	this.allChecked = false

	this.$delay(function () {
		let that = this
		this.$watch("allChecked", function (b) {
			let boxes = that.$refs.certCheckboxes
			boxes.forEach(function (box) {
				if (b) {
					box.check()
				} else {
					box.uncheck()
				}
				that.changeCerts()
			})
		})
	})

	this.changeCerts = function () {
		let boxes = this.$refs.certCheckboxes
		let that = this
		this.certIds = []
		boxes.forEach(function (box) {
			if (box.isChecked()) {
				let boxId = box.id
				that.certIds.push(parseInt(boxId.split("_")[1]))
			}
		})
	}

	this.resetAllCerts = function () {
		this.$post(".resetAll")
			.success(function () {
				teaweb.successRefresh("操作成功，将很快开始重试")
			})
	}

	this.resetCerts = function () {
		this.$post(".reset")
			.params({ certIds: this.certIds })
			.success(function () {
				teaweb.successRefresh("操作成功，将很快开始重试")
			})
	}

	this.ignoreCerts = function () {
		this.$post(".ignore")
			.params({ certIds: this.certIds })
			.success(function () {
				teaweb.successRefresh("忽略成功")
			})
	}

	// 查看证书详情
	this.viewCert = function (certId) {
		teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
			height: "28em",
			width: "48em"
		})
	}
})