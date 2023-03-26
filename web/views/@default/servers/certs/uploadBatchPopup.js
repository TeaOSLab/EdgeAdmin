Tea.context(function () {
	this.isRequesting = false

	this.before = function () {
		this.isRequesting = true
	}

	this.done = function () {
		this.isRequesting = false
	}

	this.successUpload = function (resp) {
		let msg = "html:成功上传" + resp.data.count + "个证书"
		if (this.userId > 0) {
			msg += "<br/>由于你选择了证书用户，所以只有此用户才能在用户系统中查看到这些证书。"
		}
		teaweb.success(msg, function () {
			NotifyPopup(resp)
		})
	}

	this.changeUserId = function (userId) {
		this.userId = userId
	}
})