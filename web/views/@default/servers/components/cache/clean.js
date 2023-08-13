Tea.context(function () {
	this.reason = 0

	this.REASON_NEW_PIE = 0
	this.REASON_ISSUE_REPORT = 1
	this.REASON_BATCH_DELETE = 2
	this.REASON_MAINTAINS = 3

	this.isReasonable = function () {
		return this.reason == this.REASON_ISSUE_REPORT || this.reason == this.REASON_BATCH_DELETE || this.reason == this.REASON_MAINTAINS
	}

	if (this.clusterId == null) {
		if (this.clusters.length > 0) {
			this.clusterId = this.clusters[0].id
		} else {
			this.clusterId = 0
		}
	}

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