Tea.context(function () {
	this.createRecord = function () {
		teaweb.popup("/ns/domains/records/createPopup?domainId=" + this.domain.id, {
			height: "33em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.updateRecord = function (recordId) {
		teaweb.popup("/ns/domains/records/updatePopup?recordId=" + recordId, {
			height: "33em",
			callback: function () {
				teaweb.success("保存成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteRecord = function (recordId) {
		let that = this
		teaweb.confirm("确定要删除此记录吗？", function () {
			that.$post(".delete")
				.params({
					recordId: recordId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}

	this.formatTTL = function (ttl) {
		if (ttl % 86400 == 0) {
			let days = ttl / 86400
			return days + "天"
		}
		if (ttl % 3600 == 0) {
			let hours = ttl / 3600
			return hours + "小时"
		}
		if (ttl % 60 == 0) {
			let minutes = ttl / 60
			return minutes + "分钟"
		}
		return ttl + "秒"
	}
})