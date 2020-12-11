Tea.context(function () {
	this.generateBills = function () {
		let that = this
		teaweb.confirm("确定要生成上个月的账单吗？", function () {
			that.$post(".generate")
				.params({
					month: that.month
				})
				.success(function () {
					teaweb.success("生成成功", function () {
						window.location = "/finance/bills"
					})
				})
		})
	}
})