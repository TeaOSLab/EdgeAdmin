Tea.context(function () {
	if (this.clusterId == null) {
		if (this.clusters.length > 0) {
			this.clusterId = this.clusters[0].id
		} else {
			this.clusterId = 0
		}
	}

	this.isRequestingWrite = false
	this.writeOk = false
	this.writeMessage = ""
	this.writeIsAllOk = false
	this.writeResults = []

	this.beforeWrite = function () {
		this.isRequestingWrite = true
		this.writeOk = false
		this.writeMessage = ""
		this.writeResult = {}
	}

	this.failWrite = function (resp) {
		this.writeOk = false
		this.writeMessage = resp.message
	}

	this.successWrite = function (resp) {
		this.writeOk = true
		this.writeIsAllOk = resp.data.isAllOk
		this.writeResults = resp.data.results
	}

	this.doneWrite = function () {
		this.isRequestingWrite = false
	}

	this.isRequestingRead = false
	this.readOk = false
	this.readMessage = ""
	this.readIsAllOk = false
	this.readResults = []

	this.beforeRead = function () {
		this.isRequestingRead = true
		this.readOk = false
		this.readMessage = ""
		this.readResult = {}
	}

	this.failRead = function (resp) {
		this.readOk = false
		this.readMessage = resp.message
	};

	this.successRead = function (resp) {
		this.readOk = true;
		this.readIsAllOk = resp.data.isAllOk
		this.readResults = resp.data.results
	}

	this.doneRead = function () {
		this.isRequestingRead = false
	}
});