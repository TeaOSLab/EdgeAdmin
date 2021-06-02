Vue.component("ns-access-log-box", {
	props: ["v-access-log"],
	data: function () {
		let accessLog = this.vAccessLog
		return {
			accessLog: accessLog
		}
	},
	methods: {
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()

			teaweb.popup("/ns/clusters/accessLogs/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "24em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		}
	},
	template: `<div class="access-log-row" :style="{'color': (accessLog.nsRecordId == null || accessLog.nsRecordId == 0) ? '#dc143c' : ''}" ref="box">
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> {{accessLog.remoteAddr}} [{{accessLog.timeLocal}}] [{{accessLog.networking}}] <em>{{accessLog.questionType}} {{accessLog.questionName}}</em> -&gt; <em>{{accessLog.recordType}} {{accessLog.recordValue}}</em><!-- &nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>-->
</div>`
})