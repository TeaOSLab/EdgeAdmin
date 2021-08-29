Tea.context(function () {
	this.isRequesting = false
	this.isOk = false
	this.message = ""
	this.results = []

	this.before = function () {
		this.isRequesting = true
		this.isOk = false
		this.message = ""
		this.results = []
	}

	this.success = function (resp) {
		this.isOk = true
		this.results = resp.data.results
	}

	this.fail = function (resp) {
		this.message = resp.message
	}

	this.done = function () {
		this.isRequesting = false
	}
});