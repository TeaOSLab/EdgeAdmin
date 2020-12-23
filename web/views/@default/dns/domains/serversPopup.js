Tea.context(function () {
	this.keyword = ""
	this.status = ""

	let allServers = []
	this.clusters.forEach(function (cluster) {
		let servers = cluster.servers
		servers.forEach(function (server) {
			server.cluster = cluster
			allServers.push(server)
		})
	})

	this.servers = allServers

	this.$delay(function () {
		this.$watch("keyword", function () {
			this.reloadServers()
		})
		this.$watch("status", function () {
			this.reloadServers()
		})
	})

	this.reloadServers = function () {
		let that = this
		this.servers = allServers.$copy().$findAll(function (k, v) {
			if (that.keyword.length > 0
				&& !teaweb.match(v.cluster.name, that.keyword)
				&& !teaweb.match(v.cluster.dnsName, that.keyword)
				&& !teaweb.match(v.name, that.keyword)
				&& !teaweb.match(v.dnsName, that.keyword)) {
				return false
			}
			if (that.status == "ok" && !v.isOk) {
				return false
			}
			if (that.status == "notOk" && v.isOk) {
				return false
			}
			return true
		})
	}
})