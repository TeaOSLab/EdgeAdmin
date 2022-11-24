Vue.component("ns-access-log-box", {
	props: ["v-access-log", "v-keyword"],
	data: function () {
		let accessLog = this.vAccessLog
		let isFailure = false

		if (accessLog.isRecursive) {
			if (accessLog.recordValue == null || accessLog.recordValue.length == 0) {
				isFailure = true
			}
		} else {
			if (accessLog.recordType == "SOA" || accessLog.recordType == "NS") {
				if (accessLog.recordValue == null || accessLog.recordValue.length == 0) {
					isFailure = true
				}
			}

			// 没有找到记录的不需要高亮显示，防止管理员看到红色就心理恐慌
		}

		return {
			accessLog: accessLog,
			isFailure: isFailure
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
	template: `<div class="access-log-row" :style="{'color': isFailure ? '#dc143c' : ''}" ref="box">
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword> [{{accessLog.timeLocal}}] [{{accessLog.networking}}] <em>{{accessLog.questionType}} <keyword :v-word="vKeyword">{{accessLog.questionName}}</keyword></em> -&gt; 
	
	<span v-if="accessLog.recordType != null && accessLog.recordType.length > 0"><em>{{accessLog.recordType}} <keyword :v-word="vKeyword">{{accessLog.recordValue}}</keyword></em></span>
	<span v-else class="disabled">&nbsp;[没有记录]</span>
	
	<!-- &nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>-->
	<div v-if="(accessLog.nsRoutes != null && accessLog.nsRoutes.length > 0) || accessLog.isRecursive" style="margin-top: 0.3em">
		<span class="ui label tiny basic grey" v-for="route in accessLog.nsRoutes">线路: {{route.name}}</span>
		<span class="ui label tiny basic grey" v-if="accessLog.isRecursive">递归DNS</span>
	</div>
	<div v-if="accessLog.error != null && accessLog.error.length > 0" style="color:#dc143c">
		<i class="icon warning circle"></i>错误：[{{accessLog.error}}]
	</div>
</div>`
})