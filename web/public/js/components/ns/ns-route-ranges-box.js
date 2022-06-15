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
			isAddingBatch: false,

			// IP范围
			ipRangeFrom: "",
			ipRangeTo: "",

			batchIPRange: ""
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
		addBatch: function () {
			this.isAddingBatch = true
			let that = this
			setTimeout(function () {
				that.$refs.batchIPRange.focus()
			}, 100)
		},
		cancelBatchIPRange: function () {
			this.isAddingBatch = false
			this.batchIPRange = ""
		},
		confirmBatchIPRange: function () {
			let that = this
			let rangesText = this.batchIPRange
			if (rangesText.length == 0) {
				teaweb.warn("请填写要加入的IP范围", function () {
					that.$refs.batchIPRange.focus()
				})
				return
			}

			let validRanges = []
			let invalidLine = ""
			rangesText.split("\n").forEach(function (line) {
				line = line.trim()
				if (line.length == 0) {
					return
				}
				line = line.replace("，", ",")
				let pieces = line.split(",")
				if (pieces.length != 2) {
					invalidLine = line
					return
				}
				let ipFrom = pieces[0].trim()
				let ipTo = pieces[1].trim()
				if (!that.validateIP(ipFrom) || !that.validateIP(ipTo)) {
					invalidLine = line
					return
				}
				validRanges.push({
					type: "ipRange",
					params: {
						ipFrom: ipFrom,
						ipTo: ipTo
					}
				})
			})
			if (invalidLine.length > 0) {
				teaweb.warn("'" + invalidLine + "'格式错误", function () {
					that.$refs.batchIPRange.focus()
				})
				return
			}
			validRanges.forEach(function (v) {
				that.ranges.push(v)
			})
			this.cancelBatchIPRange()
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
		<div class="ui label tiny basic" v-for="(range, index) in ranges" style="margin-bottom: 0.3em">
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			{{range.params.ipFrom}} - {{range.params.ipTo}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- 添加单个 -->
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

	<!-- 添加多个 -->
	<div style="margin-bottom: 1em" v-show="isAddingBatch">
		<div class="ui field">
			<textarea rows="5" ref="batchIPRange" v-model="batchIPRange"></textarea>	
			<p class="comment">每行一条，格式为<code-label>开始IP,结束IP</code-label>，比如<code-label>192.168.1.100,192.168.1.200</code-label>。</p>	
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmBatchIPRange">确定</button> &nbsp;
			<a href="" @click.prevent="cancelBatchIPRange" title="取消"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="!isAdding && !isAddingBatch">
		<button class="ui button tiny" type="button" @click.prevent="add">单个添加</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="addBatch">批量添加</button>
	</div>
</div>`
})