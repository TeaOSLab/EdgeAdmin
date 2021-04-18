Tea.context(function () {
	this.isRequesting = false
	this.resp = null

	this.success = function (resp) {
		this.resp = resp.data
	}

	this.requestBefore = function () {
		this.isRequesting = true
		this.resp = null
	}

	this.requestDone = function () {
		this.isRequesting = false
	}
})