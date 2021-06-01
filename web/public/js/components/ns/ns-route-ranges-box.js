Vue.component("ns-route-ranges-box", {
	props: ["v-ranges"],
	data: function () {
		let ranges = this.vRanges
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,

			// IP范围
			ipRangeFrom: "",
			ipRangeTo: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipRangeFrom.focus()
			}, 100)
		},
		remove: function (index) {
			this.ranges.$remove(index)
		},
		cancelIPRange: function () {
			this.isAdding = false
			this.ipRangeFrom = ""
			this.ipRangeTo = ""
		},
		confirmIPRange: function () {
			// 校验IP
			let that = this
			this.ipRangeFrom = this.ipRangeFrom.trim()
			if (!this.validateIP(this.ipRangeFrom)) {
				teaweb.warn("开始IP填写错误", function () {
					that.$refs.ipRangeFrom.focus()
				})
				return
			}

			this.ipRangeTo = this.ipRangeTo.trim()
			if (!this.validateIP(this.ipRangeTo)) {
				teaweb.warn("结束IP填写错误", function () {
					that.$refs.ipRangeTo.focus()
				})
				return
			}

			this.ranges.push({
				type: "ipRange",
				params: {
					ipFrom: this.ipRangeFrom,
					ipTo: this.ipRangeTo
				}
			})
			this.cancelIPRange()
		},
		validateIP: function (ip) {
			if (!ip.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/)) {
				return false
			}
			let pieces = ip.split(".")
			let isOk = true
			pieces.forEach(function (v) {
				let v1 = parseInt(v)
				if (v1 > 255) {
					isOk = false
				}
			})
			return isOk
		}
	},
	template: `<div>
	<input type="hidden" name="rangesJSON" :value="JSON.stringify(ranges)"/>
	<div v-if="ranges.length > 0">
		<div class="ui label tiny basic" v-for="(range, index) in ranges">
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			{{range.params.ipFrom}} - {{range.params.ipTo}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- IP 范围 -->
	<div style="margin-bottom: 1em" v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="开始IP" maxlength="15" size="15" v-model="ipRangeFrom" ref="ipRangeFrom"  @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">-</div>
			<div class="ui field">
				<input type="text" placeholder="结束IP" maxlength="15" size="15" v-model="ipRangeTo" ref="ipRangeTo" @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirmIPRange">确定</button> &nbsp;
				<a href="" @click.prevent="cancelIPRange" title="取消"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`
})