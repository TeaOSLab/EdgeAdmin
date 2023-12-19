Vue.component("server-traffic-limit-status-viewer", {
	props: ["value"],
	data: function () {
		let targetTypeName = "流量"
		if (this.value != null) {
			targetTypeName = this.targetTypeToName(this.value.targetType)
		}

		return {
			status: this.value,
			targetTypeName: targetTypeName
		}
	},
	methods: {
		targetTypeToName: function (targetType) {
			switch (targetType) {
				case "traffic":
					return "流量"
				case "request":
					return "请求数"
				case "websocketConnections":
					return "Websocket连接数"
			}
			return "流量"
		}
	},
	template: `<span v-if="status != null">
	<span v-if="status.dateType == 'day'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>当日{{targetTypeName}}限制</span>
	<span v-if="status.dateType == 'month'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>当月{{targetTypeName}}限制</span>
	<span v-if="status.dateType == 'total'" class="small red">已达到<span v-if="status.planId > 0">套餐</span>总体{{targetTypeName}}限制</span>
</span>`
})