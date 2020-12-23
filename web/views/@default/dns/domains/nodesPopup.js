Tea.context(function () {
	this.keyword = ""
	this.status = ""

	let allNodes = []
	this.clusters.forEach(function (cluster) {
		let nodes = cluster.nodes
		nodes.forEach(function (node) {
			node.cluster = cluster
			allNodes.push(node)
		})
	})

	this.nodes = allNodes

	this.$delay(function () {
		this.$watch("keyword", function () {
			this.reloadNodes()
		})
		this.$watch("status", function () {
			this.reloadNodes()
		})
	})

	this.reloadNodes = function () {
		let that = this
		this.nodes = allNodes.$copy().$findAll(function (k, v) {
			if (that.keyword.length > 0
				&& !teaweb.match(v.cluster.name, that.keyword)
				&& !teaweb.match(v.cluster.dnsName, that.keyword)
				&& !teaweb.match(v.name, that.keyword)
				&& !teaweb.match(v.ipAddr, that.keyword)
				&& !teaweb.match(v.route.name, that.keyword)) {
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