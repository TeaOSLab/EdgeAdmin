Tea.context(function () {
	this.hasHTTPS = false
	this.grpcAddrs = []
	this.restAddrs = []

	this.changeListens = function (addrs) {
		this.grpcAddrs = addrs

		this.hasHTTPS = this.grpcAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}) || (this.restIsOn && this.restAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}))
	}

	this.changeRestListens = function (addrs) {
		this.restAddrs = addrs

		this.hasHTTPS = this.grpcAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}) || (this.restIsOn && this.restAddrs.$any(function (k, v) {
			return v.protocol == "https"
		}))
	}

	this.restIsOn = false
})