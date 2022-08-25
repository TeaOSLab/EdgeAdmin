Vue.component("http-access-log-partitions-box", {
	props: ["v-partition", "v-day", "v-query"],
	mounted: function () {
		let that = this
		Tea.action("/servers/logs/partitionData")
			.params({
				day: this.vDay
			})
			.success(function (resp) {
				that.partitions = []
				resp.data.partitions.reverse().forEach(function (v) {
					that.partitions.push({
						code: v,
						isDisabled: false,
						hasLogs: false
					})
				})
				if (that.partitions.length > 0) {
					if (that.vPartition == null || that.vPartition < 0) {
						that.selectedPartition = that.partitions[0].code
					}

					if (that.partitions.length > 1) {
						that.checkLogs()
					}
				}
			})
			.post()
	},
	data: function () {
		return {
			partitions: [],
			selectedPartition: this.vPartition,
			checkingPartition: 0
		}
	},
	methods: {
		url: function (p) {
			let u = window.location.toString()
			u = u.replace(/\?partition=-?\d+/, "?")
			u = u.replace(/\?requestId=-?\d+/, "?")
			u = u.replace(/&partition=-?\d+/, "")
			u = u.replace(/&requestId=-?\d+/, "")
			if (u.indexOf("?") > 0) {
				u += "&partition=" + p
			} else {
				u += "?partition=" + p
			}
			return u
		},
		disable: function (partition) {
			this.partitions.forEach(function (p) {
				if (p.code == partition) {
					p.isDisabled = true
				}
			})
		},
		checkLogs: function () {
			let that = this
			let index = this.checkingPartition
			let params = {
				partition: index
			}
			let query = this.vQuery
			if (query == null || query.length == 0) {
				return
			}
			query.split("&").forEach(function (v) {
				let param = v.split("=")
				params[param[0]] = decodeURIComponent(param[1])
			})
			Tea.action("/servers/logs/hasLogs")
				.params(params)
				.post()
				.success(function (response) {
					if (response.data.hasLogs) {
						// 因为是倒序，所以这里需要使用总长度减去index
						that.partitions[that.partitions.length - 1 - index].hasLogs = true
					}

					index++
					if (index >= that.partitions.length) {
						return
					}
					that.checkingPartition = index
					that.checkLogs()
				})
		}
	},
	template: `<div v-if="partitions.length > 1">
	<div class="ui divider" style="margin-bottom: 0"></div>
	<div class="ui menu text small blue" style="margin-bottom: 0; margin-top: 0">
		<a v-for="(p, index) in partitions" :href="url(p.code)" class="item" :class="{active: selectedPartition == p.code, disabled: p.isDisabled}">分表{{p.code+1}} <span v-if="p.hasLogs">&nbsp; <dot></dot></span> &nbsp; &nbsp; <span class="disabled" v-if="index != partitions.length - 1">|</span></a>
	</div>
	<div class="ui divider" style="margin-top: 0"></div>
</div>`
})