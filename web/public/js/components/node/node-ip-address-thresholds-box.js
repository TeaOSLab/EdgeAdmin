// 节点IP阈值
Vue.component("node-ip-address-thresholds-box", {
	props: ["v-thresholds"],
	data: function () {
		let thresholds = this.vThresholds
		if (thresholds == null) {
			thresholds = []
		}

		let avgRequests = {
			duration: "",
			operator: "lte",
			value: ""
		}

		let avgTrafficOut = {
			duration: "",
			operator: "lte",
			value: ""
		}

		let avgTrafficIn = {
			duration: "",
			operator: "lte",
			value: ""
		}

		thresholds.forEach(function (v) {
			switch (v.item) {
				case "avgRequests":
					avgRequests.duration = v.duration
					avgRequests.operator = v.operator
					avgRequests.value = v.value.toString()
					break
				case "avgTrafficOut":
					avgTrafficOut.duration = v.duration
					avgTrafficOut.operator = v.operator
					avgTrafficOut.value = v.value.toString()
					break
				case "avgTrafficIn":
					avgTrafficIn.duration = v.duration
					avgTrafficIn.operator = v.operator
					avgTrafficIn.value = v.value.toString()
					break
			}
		})

		return {
			thresholds: thresholds,
			avgRequests: avgRequests,
			avgTrafficOut: avgTrafficOut,
			avgTrafficIn: avgTrafficIn
		}
	},
	watch: {
		"avgRequests.duration": function () {
			this.compose()
		},
		"avgRequests.operator": function () {
			this.compose()
		},
		"avgRequests.value": function () {
			this.compose()
		},
		"avgTrafficOut.duration": function () {
			this.compose()
		},
		"avgTrafficOut.operator": function () {
			this.compose()
		},
		"avgTrafficOut.value": function () {
			this.compose()
		},
		"avgTrafficIn.duration": function () {
			this.compose()
		},
		"avgTrafficIn.operator": function () {
			this.compose()
		},
		"avgTrafficIn.value": function () {
			this.compose()
		}
	},
	methods: {
		compose: function () {
			let thresholds = []

			// avg requests
			{
				let duration = parseInt(this.avgRequests.duration)
				let value = parseInt(this.avgRequests.value)
				if (!isNaN(duration) && duration > 0 && !isNaN(value) && value > 0) {
					thresholds.push({
						item: "avgRequests",
						operator: this.avgRequests.operator,
						duration: duration,
						durationUnit: "minute",
						value: value
					})
				}
			}

			// avg traffic out
			{
				let duration = parseInt(this.avgTrafficOut.duration)
				let value = parseInt(this.avgTrafficOut.value)
				if (!isNaN(duration) && duration > 0 && !isNaN(value) && value > 0) {
					thresholds.push({
						item: "avgTrafficOut",
						operator: this.avgTrafficOut.operator,
						duration: duration,
						durationUnit: "minute",
						value: value
					})
				}
			}

			// avg requests
			{
				let duration = parseInt(this.avgTrafficIn.duration)
				let value = parseInt(this.avgTrafficIn.value)
				if (!isNaN(duration) && duration > 0 && !isNaN(value) && value > 0) {
					thresholds.push({
						item: "avgTrafficIn",
						operator: this.avgTrafficIn.operator,
						duration: duration,
						durationUnit: "minute",
						value: value
					})
				}
			}

			this.thresholds = thresholds
		}
	},
	template: `<div>
	<input type="hidden" name="thresholdsJSON" :value="JSON.stringify(thresholds)"/>
	<table class="ui table celled">
		<thead>
			<tr>
				<td>统计项目</td>
				<th>统计周期</th>
				<th>操作符</th>
				<th>对比值</th>
			</tr>	
		</thead>
		<tr>
			<td>平均请求数/秒</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 5em" v-model="avgRequests.duration"/>
					<span class="ui label">分钟</span>
				</div>
			</td>
			<td>
				<select class="ui dropdown auto-width" v-model="avgRequests.operator">
					<option value="lte">小于等于</option>
					<option value="gt">大于</option>
				</select>
			</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 6em" v-model="avgRequests.value"/>
					<span class="ui label">个</span>
				</div>
			</td>
		</tr>
		<tr>
			<td class="title">平均下行流量/秒</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 5em" v-model="avgTrafficOut.duration"/>
					<span class="ui label">分钟</span>
				</div>
			</td>
			<td>
				<select class="ui dropdown auto-width" v-model="avgTrafficOut.operator">
					<option value="lte">小于等于</option>
					<option value="gt">大于</option>
				</select>
			</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 6em"  v-model="avgTrafficOut.value"/>
					<span class="ui label">MB</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>平均上行流量/秒</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 5em"  v-model="avgTrafficIn.duration"/>
					<span class="ui label">分钟</span>
				</div>
			</td>
			<td>
				<select class="ui dropdown auto-width" v-model="avgTrafficIn.operator">
					<option value="lte">小于等于</option>
					<option value="gt">大于</option>
				</select>
			</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" style="width: 6em" v-model="avgTrafficIn.value"/>
					<span class="ui label">MB</span>
				</div>
			</td>
		</tr>
	</table>
	<p class="comment">满足所有阈值条件时IP才会上线，否则下线。统计周期和对比值设置为0表示没有限制。各个输入项只支持整数数字。</p>
</div>`
})