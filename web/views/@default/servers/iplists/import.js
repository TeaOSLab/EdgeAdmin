Tea.context(function () {
	this.success = function (resp) {
		teaweb.success("成功导入" + resp.data.count + "个IP", function () {
			teaweb.reload()
		})
	}
})