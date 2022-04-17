Vue.component("http-access-log-partitions-box", {
	props: ["v-partition", "v-day"],
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
						isDisabled: false
					})
				})
				if (that.partitions.length > 0) {
					if (that.vPartition == null || that.vPartition < 0) {
						that.selectedPartition = that.partitions[0].code
					}
				}
			})
			.post()
	},
	data: function () {
		return {
			partitions: [],
			selectedPartition: this.vPartition
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
		}
	},
	template: `<div v-if="partitions.length > 1">
	<div class="ui divider" style="margin-bottom: 0"></div>
	<div class="ui menu text small blue" style="margin-bottom: 0; margin-top: 0">
		<a v-for="(p, index) in partitions" :href="url(p.code)" class="item" :class="{active: selectedPartition == p.code, disabled: p.isDisabled}">分表{{p.code+1}} &nbsp; &nbsp; <span class="disabled" v-if="index != partitions.length - 1">|</span></a>
	</div>
	<div class="ui divider" style="margin-top: 0"></div>
</div>`
})