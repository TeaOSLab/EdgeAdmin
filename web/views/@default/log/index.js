Tea.context(function () {
	this.$delay(function () {
		teaweb.datepicker("day-from-picker")
		teaweb.datepicker("day-to-picker")
	})

	this.logs.forEach(function (v) {
		v.moreVisible = false
	})

	this.showMore = function (log) {
		log.moreVisible = !log.moreVisible
	}

	this.exportExcel = function () {
		let that = this
		teaweb.confirm("确定要将当前列表导出到Excel吗？", function () {
			window.location = "/log/exportExcel?dayFrom=" + that.dayFrom + "&dayTo=" + that.dayTo + "&keyword=" + that.keyword + "&userType=" + that.userType + '&level=' + that.level
		})
	}

	this.deleteLog = function (logId) {
		let that = this
		teaweb.confirm("确定要删除此日志吗？", function () {
			that.$post(".delete")
				.params({
					logId: logId
				})
				.refresh()
		})
	}
})