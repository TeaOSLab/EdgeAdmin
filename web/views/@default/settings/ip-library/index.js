Tea.context(function () {
	this.upload = function () {
		teaweb.popup("/settings/ip-library/uploadPopup", {
			callback: function () {
				teaweb.success("上传成功", function () {
					teaweb.reload()
				})
			}
		})
	}

	this.deleteLibrary = function (libraryId) {
		let that = this
		teaweb.confirm("确定要删除此IP库吗？", function () {
			that.$post(".delete")
				.params({
					"libraryId": libraryId
				})
				.success(function () {
					teaweb.reload()
				})
		})
	}
})