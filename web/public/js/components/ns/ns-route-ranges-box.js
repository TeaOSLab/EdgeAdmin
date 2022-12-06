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

			// 类型
			rangeType: "ipRange",
			isReverse: false,

			// IP范围
			ipRangeFrom: "",
			ipRangeTo: "",

			batchIPRange: "",

			// CIDR
			ipCIDR: "",
			batchIPCIDR: "",

			// region
			regions: [],
			regionType: "country",
			regionConnector: "OR"
		}
	},
	methods: {
		addIPRange: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipRangeFrom.focus()
			}, 100)
		},
		addCIDR: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipCIDR.focus()
			}, 100)
		},
		addRegions: function () {
			this.isAdding = true
		},
		addRegion: function (regionType) {
			this.regionType = regionType
		},
		remove: function (index) {
			this.ranges.$remove(index)
		},
		cancelIPRange: function () {
			this.isAdding = false
			this.ipRangeFrom = ""
			this.ipRangeTo = ""
			this.isReverse = false
		},
		cancelIPCIDR: function () {
			this.isAdding = false
			this.ipCIDR = ""
			this.isReverse = false
		},
		cancelRegions: function () {
			this.isAdding = false
			this.regions = []
			this.regionType = "country"
			this.regionConnector = "OR"
			this.isReverse = false
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
					ipTo: this.ipRangeTo,
					isReverse: this.isReverse
				}
			})
			this.cancelIPRange()
		},
		confirmIPCIDR: function () {
			let that = this
			if (this.ipCIDR.length == 0) {
				teaweb.warn("请填写CIDR", function () {
					that.$refs.ipCIDR.focus()
				})
				return
			}
			if (!this.validateCIDR(this.ipCIDR)) {
				teaweb.warn("请输入正确的CIDR", function () {
					that.$refs.ipCIDR.focus()
				})
				return
			}


			this.ranges.push({
				type: "cidr",
				params: {
					cidr: this.ipCIDR,
					isReverse: this.isReverse
				}
			})
			this.cancelIPCIDR()
		},
		confirmRegions: function () {
			if (this.regions.length == 0) {
				this.cancelRegions()
				return
			}
			this.ranges.push({
				type: "region",
				connector: this.regionConnector,
				params: {
					regions: this.regions,
					isReverse: this.isReverse
				}
			})
			this.cancelRegions()
		},
		addBatchIPRange: function () {
			this.isAddingBatch = true
			let that = this
			setTimeout(function () {
				that.$refs.batchIPRange.focus()
			}, 100)
		},
		addBatchCIDR: function () {
			this.isAddingBatch = true
			let that = this
			setTimeout(function () {
				that.$refs.batchIPCIDR.focus()
			}, 100)
		},
		cancelBatchIPRange: function () {
			this.isAddingBatch = false
			this.batchIPRange = ""
			this.isReverse = false
		},
		cancelBatchIPCIDR: function () {
			this.isAddingBatch = false
			this.batchIPCIDR = ""
			this.isReverse = false
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
						ipTo: ipTo,
						isReverse: that.isReverse
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
		confirmBatchIPCIDR: function () {
			let that = this
			let rangesText = this.batchIPCIDR
			if (rangesText.length == 0) {
				teaweb.warn("请填写要加入的CIDR", function () {
					that.$refs.batchIPCIDR.focus()
				})
				return
			}

			let validRanges = []
			let invalidLine = ""
			rangesText.split("\n").forEach(function (line) {
				let cidr = line.trim()
				if (cidr.length == 0) {
					return
				}
				if (!that.validateCIDR(cidr)) {
					invalidLine = line
					return
				}
				validRanges.push({
					type: "cidr",
					params: {
						cidr: cidr,
						isReverse: that.isReverse
					}
				})
			})
			if (invalidLine.length > 0) {
				teaweb.warn("'" + invalidLine + "'格式错误", function () {
					that.$refs.batchIPCIDR.focus()
				})
				return
			}
			validRanges.forEach(function (v) {
				that.ranges.push(v)
			})
			this.cancelBatchIPCIDR()
		},
		selectRegionCountry: function (country) {
			if (country == null) {
				return
			}
			this.regions.push({
				type: "country",
				id: country.id,
				name: country.name
			})
			this.$refs.regionCountryComboBox.clear()
		},
		selectRegionProvince: function (province) {
			if (province == null) {
				return
			}
			this.regions.push({
				type: "province",
				id: province.id,
				name: province.name
			})
			this.$refs.regionProvinceComboBox.clear()
		},
		selectRegionCity: function (city) {
			if (city == null) {
				return
			}
			this.regions.push({
				type: "city",
				id: city.id,
				name: city.name
			})
			this.$refs.regionCityComboBox.clear()
		},
		selectRegionProvider: function (provider) {
			if (provider == null) {
				return
			}
			this.regions.push({
				type: "provider",
				id: provider.id,
				name: provider.name
			})
			this.$refs.regionProviderComboBox.clear()
		},
		removeRegion: function (index) {
			this.regions.$remove(index)
		},
		validateIP: function (ip) {
			if (ip.length == 0) {
				return
			}

			// IPv6
			if (ip.indexOf(":") >= 0) {
				let pieces = ip.split(":")
				if (pieces.length > 8) {
					return false
				}
				let isOk = true
				pieces.forEach(function (piece) {
					if (!/^[\da-fA-F]{0,4}$/.test(piece)) {
						isOk = false
					}
				})

				return isOk
			}

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
		},
		validateCIDR: function (cidr) {
			let pieces = cidr.split("/")
			if (pieces.length != 2) {
				return false
			}
			let ip = pieces[0]
			if (!this.validateIP(ip)) {
				return false
			}
			let mask = pieces[1]
			if (!/^\d{1,3}$/.test(mask)) {
				return false
			}
			mask = parseInt(mask, 10)
			if (cidr.indexOf(":") >= 0) { // IPv6
				return mask <= 128
			}
			return mask <= 32
		},
		updateRangeType: function (rangeType) {
			this.rangeType = rangeType
		}
	},
	template: `<div>
	<input type="hidden" name="rangesJSON" :value="JSON.stringify(ranges)"/>
	<div v-if="ranges.length > 0">
		<div class="ui label tiny basic" v-for="(range, index) in ranges" style="margin-bottom: 0.3em">
			<span class="red" v-if="range.params.isReverse">[排除]</span>
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			<span v-if="range.type == 'cidr'">CIDR：</span>
			<span v-if="range.type == 'region'"></span>
			<span v-if="range.type == 'ipRange'">{{range.params.ipFrom}} - {{range.params.ipTo}}</span>
			<span v-if="range.type == 'cidr'">{{range.params.cidr}}</span>
			<span v-if="range.type == 'region'">
				<span v-for="(region, index) in range.params.regions">
					<span v-if="region.type == 'country'">国家/地区</span>
					<span v-if="region.type == 'province'">省份</span>
					<span v-if="region.type == 'city'">城市</span>
					<span v-if="region.type == 'provider'">ISP</span>
					：{{region.name}}
					<span v-if="index < range.params.regions.length - 1" class="grey">
						&nbsp;
						<span v-if="range.connector == 'OR' || range.connector == '' || range.connector == null">或</span>
						<span v-if="range.connector == 'AND'">且</span>
						&nbsp;
					</span>
				</span>
			</span>
			 &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- IP范围 -->
	<div v-if="rangeType == 'ipRange'">
		<!-- 添加单个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAdding">
			<table class="ui table">
				<tr>
					<td class="title">开始IP *</td>
					<td>
						<input type="text" placeholder="开始IP" maxlength="40" size="40" style="width: 15em" v-model="ipRangeFrom" ref="ipRangeFrom"  @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td>结束IP *</td>
					<td>
						<input type="text" placeholder="结束IP" maxlength="40" size="40" style="width: 15em" v-model="ipRangeTo" ref="ipRangeTo" @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmIPRange">确定</button> &nbsp;
					<a href="" @click.prevent="cancelIPRange" title="取消"><i class="icon remove small"></i></a>
		</div>
	
		<!-- 添加多个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAddingBatch">
			<table class="ui table">
				<tr>
					<td class="title">IP范围列表 *</td>
					<td>
						<textarea rows="5" ref="batchIPRange" v-model="batchIPRange"></textarea>	
						<p class="comment">每行一条，格式为<code-label>开始IP,结束IP</code-label>，比如<code-label>192.168.1.100,192.168.1.200</code-label>。</p>	
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmBatchIPRange">确定</button> &nbsp;
				<a href="" @click.prevent="cancelBatchIPRange" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addIPRange">添加单个IP范围</button> &nbsp;
			<button class="ui button tiny" type="button" @click.prevent="addBatchIPRange">批量添加IP范围</button>
		</div>
	</div>
	
	<!-- CIDR -->
	<div v-if="rangeType == 'cidr'">
		<!-- 添加单个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAdding">
			<table class="ui table">
				<tr>
					<td class="title">CIDR *</td>
					<td>
						<input type="text" placeholder="IP/MASK" maxlength="40" size="40" style="width: 15em" v-model="ipCIDR" ref="ipCIDR"  @keyup.enter="confirmIPCIDR" @keypress.enter.prevent="1"/>
						<p class="comment">类似于<code-label>192.168.2.1/24</code-label>。</p>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmIPCIDR">确定</button> &nbsp;
					<a href="" @click.prevent="cancelIPCIDR" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<!-- 添加多个IP范围 -->
		<div style="margin-bottom: 1em" v-show="isAddingBatch">
			<table class="ui table">
				<tr>
					<td class="title">IP范围列表 *</td>
					<td>
						<textarea rows="5" ref="batchIPCIDR" v-model="batchIPCIDR"></textarea>	
						<p class="comment">每行一条，格式为<code-label>IP/MASK</code-label>，比如<code-label>192.168.2.1/24</code-label>。</p>	
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmBatchIPCIDR">确定</button> &nbsp;
				<a href="" @click.prevent="cancelBatchIPCIDR" title="取消"><i class="icon remove small"></i></a>
		</div>
		
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addCIDR">添加单个CIDR</button> &nbsp;
			<button class="ui button tiny" type="button" @click.prevent="addBatchCIDR">批量添加CIDR</button>
		</div>
	</div>
	
	<!-- 区域 -->
	<div v-if="rangeType == 'region'">
		<!-- 添加区域 -->
		<div v-if="isAdding">
			<table class="ui table">
				<tr>
					<td>已添加</td>
					<td>
						<span v-for="(region, index) in regions">
							<span class="ui label small basic">
								<span v-if="region.type == 'country'">国家/地区</span>
								<span v-if="region.type == 'province'">省份</span>
								<span v-if="region.type == 'city'">城市</span>
								<span v-if="region.type == 'provider'">ISP</span>
								：{{region.name}} <a href="" title="删除" @click.prevent="removeRegion(index)"><i class="icon remove small"></i></a>
							</span>
							<span v-if="index < regions.length - 1" class="grey">
								&nbsp;
								<span v-if="regionConnector == 'OR' || regionConnector == ''">或</span>
								<span v-if="regionConnector == 'AND'">且</span>
								&nbsp;
							</span>
						</span>
					</td>
				</tr>
				<tr>
					<td class="title">添加新<span v-if="regionType == 'country'">国家/地区</span><span v-if="regionType == 'province'">省份</span><span v-if="regionType == 'city'">城市</span><span v-if="regionType == 'provider'">ISP</span>
					
					 *</td>
					<td>
					 	<!-- region country name -->
						<div v-if="regionType == 'country'">
							<combo-box title="" width="14em" data-url="/ui/countryOptions" data-key="countries" placeholder="点这里选择国家/地区" @change="selectRegionCountry" ref="regionCountryComboBox" key="combo-box-country"></combo-box>
						</div>
			
						<!-- region province name -->
						<div v-if="regionType == 'province'" >
							<combo-box title="" data-url="/ui/provinceOptions" data-key="provinces" placeholder="点这里选择省份" @change="selectRegionProvince" ref="regionProvinceComboBox" key="combo-box-province"></combo-box>
						</div>
			
						<!-- region city name -->
						<div v-if="regionType == 'city'" >
							<combo-box title="" data-url="/ui/cityOptions" data-key="cities" placeholder="点这里选择城市" @change="selectRegionCity" ref="regionCityComboBox" key="combo-box-city"></combo-box>
						</div>
			
						<!-- ISP Name -->
						<div v-if="regionType == 'provider'" >
							<combo-box title="" data-url="/ui/providerOptions" data-key="providers" placeholder="点这里选择ISP" @change="selectRegionProvider" ref="regionProviderComboBox" key="combo-box-isp"></combo-box>
						</div>
						
						<div style="margin-top: 1em">
							<button class="ui button tiny basic" :class="{blue: regionType == 'country'}" type="button" @click.prevent="addRegion('country')">添加国家/地区</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'province'}" type="button" @click.prevent="addRegion('province')">添加省份</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'city'}" type="button" @click.prevent="addRegion('city')">添加城市</button> &nbsp;
							<button class="ui button tiny basic" :class="{blue: regionType == 'provider'}" type="button" @click.prevent="addRegion('provider')">ISP</button> &nbsp;
						</div>
					</td>	
				</tr>
				<tr>
					<td>区域之间关系</td>
					<td>
						<select class="ui dropdown auto-width" v-model="regionConnector">
							<option value="OR">或</option>
							<option value="AND">且</option>
						</select>
						<p class="comment" v-if="regionConnector == 'OR'">匹配所选任一区域即认为匹配成功。</p>
						<p class="comment" v-if="regionConnector == 'AND'">匹配所有所选区域才认为匹配成功。</p>
					</td>
				</tr>
				<tr>
					<td>排除</td>
					<td>
						<checkbox v-model="isReverse"></checkbox>
						<p class="comment">选中后表示线路中排除当前条件。</p>
					</td>
				</tr>
			</table>
			<button class="ui button tiny" type="button" @click.prevent="confirmRegions">确定</button> &nbsp;
				<a href="" @click.prevent="cancelRegions" title="取消"><i class="icon remove small"></i></a>
		</div>
		<div v-if="!isAdding && !isAddingBatch">
			<button class="ui button tiny" type="button" @click.prevent="addRegions">添加区域</button> &nbsp;
		</div>	
	</div>
</div>`
})