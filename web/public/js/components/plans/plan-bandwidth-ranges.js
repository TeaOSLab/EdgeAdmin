Vue.component("plan-bandwidth-ranges", {
	props: ["v-ranges"],
	data: function () {
		let ranges = this.vRanges
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,

			minMB: "",
			maxMB: "",
			pricePerMB: "",
			addingRange: {
				minMB: 0,
				maxMB: 0,
				pricePerMB: 0,
				totalPrice: 0
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = !this.isAdding
			let that = this
			setTimeout(function () {
				that.$refs.minMB.focus()
			})
		},
		cancelAdding: function () {
			this.isAdding = false
		},
		confirm: function () {
			this.isAdding = false
			this.minMB = ""
			this.maxMB = ""
			this.pricePerMB = ""
			this.ranges.push(this.addingRange)
			this.ranges.$sort(function (v1, v2) {
				if (v1.minMB < v2.minMB) {
					return -1
				}
				if (v1.minMB == v2.minMB) {
					return 0
				}
				return 1
			})
			this.change()
			this.addingRange = {
				minMB: 0,
				maxMB: 0,
				pricePerMB: 0,
				totalPrice: 0
			}
		},
		remove: function (index) {
			this.ranges.$remove(index)
			this.change()
		},
		change: function () {
			this.$emit("change", this.ranges)
		}
	},
	watch: {
		minMB: function (v) {
			let minMB = parseInt(v.toString())
			if (isNaN(minMB) || minMB < 0) {
				minMB = 0
			}
			this.addingRange.minMB = minMB
		},
		maxMB: function (v) {
			let maxMB = parseInt(v.toString())
			if (isNaN(maxMB) || maxMB < 0) {
				maxMB = 0
			}
			this.addingRange.maxMB = maxMB
		},
		pricePerMB: function (v) {
			let pricePerMB = parseFloat(v.toString())
			if (isNaN(pricePerMB) || pricePerMB < 0) {
				pricePerMB = 0
			}
			this.addingRange.pricePerMB = pricePerMB
		}
	},
	template: `<div>
	<!-- 已有价格 -->
	<div v-if="ranges.length > 0">
		<div class="ui label basic small" v-for="(range, index) in ranges" style="margin-bottom: 0.5em">
			{{range.minMB}}MB - <span v-if="range.maxMB > 0">{{range.maxMB}}MB</span><span v-else>&infin;</span> &nbsp;  价格：{{range.pricePerMB}}元/MB
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">带宽下限</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="最小带宽" style="width: 7em" maxlength="10" ref="minMB" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="minMB"/>
						<span class="ui label">MB</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">带宽上限</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="最大带宽" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="maxMB"/>
						<span class="ui label">MB</span>
					</div>
					<p class="comment">如果填0，表示上不封顶。</p>
				</td>
			</tr>
			<tr>
				<td class="title">单位价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="单位价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="pricePerMB"/>
						<span class="ui label">元/MB</span>
					</div>
				</td>
			</tr>
		</table>
		<button class="ui button small" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove small"></i></a>
	</div>
	
	<!-- 按钮 -->
	<div v-if="!isAdding">
		<button class="ui button small" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})