Vue.component("traffic-map-box", {
	props: ["v-stats", "v-is-attack"],
	mounted: function () {
		this.render()
	},
	data: function () {
		let maxPercent = 0
		let isAttack = this.vIsAttack
		this.vStats.forEach(function (v) {
			let percent = parseFloat(v.percent)
			if (percent > maxPercent) {
				maxPercent = percent
			}

			v.formattedCountRequests = teaweb.formatCount(v.countRequests) + "次"
			v.formattedCountAttackRequests = teaweb.formatCount(v.countAttackRequests) + "次"
		})

		if (maxPercent < 100) {
			maxPercent *= 1.2 // 不要让某一项100%
		}

		let screenIsNarrow = window.innerWidth < 512

		return {
			isAttack: isAttack,
			stats: this.vStats,
			chart: null,
			minOpacity: 0.2,
			maxPercent: maxPercent,
			selectedCountryName: "",
			screenIsNarrow: screenIsNarrow
		}
	},
	methods: {
		render: function () {
			if (this.$el.offsetWidth < 300) {
				let that = this
				setTimeout(function () {
					that.render()
				}, 100)
				return
			}

			this.chart = teaweb.initChart(document.getElementById("traffic-map-box"));
			let that = this
			this.chart.setOption({
				backgroundColor: "white",
				grid: {
					top: 0,
					bottom: 0,
					left: 0,
					right: 0
				},
				roam: false,
				tooltip: {
					trigger: "item"
				},
				series: [{
					type: "map",
					map: "world",
					zoom: 1.3,
					selectedMode: false,
					itemStyle: {
						areaColor: "#E9F0F9",
						borderColor: "#DDD"
					},
					label: {
						show: false,
						fontSize: "10px",
						color: "#fff",
						backgroundColor: "#8B9BD3",
						padding: [2, 2, 2, 2]
					},
					emphasis: {
						itemStyle: {
							areaColor: "#8B9BD3",
							opacity: 1.0
						},
						label: {
							show: true,
							fontSize: "10px",
							color: "#fff",
							backgroundColor: "#8B9BD3",
							padding: [2, 2, 2, 2]
						}
					},
					//select: {itemStyle:{ areaColor: "#8B9BD3", opacity: 0.8 }},
					tooltip: {
						formatter: function (args) {
							let name = args.name
							let stat = null
							that.stats.forEach(function (v) {
								if (v.name == name) {
									stat = v
								}
							})

							if (stat != null) {
								return name + "<br/>流量：" + stat.formattedBytes + "<br/>流量占比：" + stat.percent + "%<br/>请求数：" + stat.formattedCountRequests + "<br/>攻击数：" + stat.formattedCountAttackRequests
							}
							return name
						}
					},
					data: this.stats.map(function (v) {
						let opacity = parseFloat(v.percent) / that.maxPercent
						if (opacity < that.minOpacity) {
							opacity = that.minOpacity
						}
						let fullOpacity = opacity * 3
						if (fullOpacity > 1) {
							fullOpacity = 1
						}
						let isAttack = that.vIsAttack
						let bgColor = "#276AC6"
						if (isAttack) {
							bgColor = "#B03A5B"
						}

						return {
							name: v.name,
							value: v.bytes,
							percent: parseFloat(v.percent),
							itemStyle: {
								areaColor: bgColor,
								opacity: opacity
							},
							emphasis: {
								itemStyle: {
									areaColor: bgColor,
									opacity: fullOpacity
								},
								label: {
									show: true,
									formatter: function (args) {
										return args.name
									}
								}
							},
							label: {
								show: false,
								formatter: function (args) {
									if (args.name == that.selectedCountryName) {
										return args.name
									}
									return ""
								},
								fontSize: "10px",
								color: "#fff",
								backgroundColor: "#8B9BD3",
								padding: [2, 2, 2, 2]
							}
						}
					}),
					nameMap: window.WorldCountriesMap
				}]
			})
			this.chart.resize()
		},
		selectCountry: function (countryName) {
			if (this.chart == null) {
				return
			}
			let option = this.chart.getOption()
			let that = this
			option.series[0].data.forEach(function (v) {
				let opacity = v.percent / that.maxPercent
				if (opacity < that.minOpacity) {
					opacity = that.minOpacity
				}

				if (v.name == countryName) {
					if (v.isSelected) {
						v.itemStyle.opacity = opacity
						v.isSelected = false
						v.label.show = false
						that.selectedCountryName = ""
						return
					}
					v.isSelected = true
					that.selectedCountryName = countryName
					opacity *= 3
					if (opacity > 1) {
						opacity = 1
					}

					// 至少是0.5，让用户能够看清
					if (opacity < 0.5) {
						opacity = 0.5
					}
					v.itemStyle.opacity = opacity
					v.label.show = true
				} else {
					v.itemStyle.opacity = opacity
					v.isSelected = false
					v.label.show = false
				}
			})
			this.chart.setOption(option)
		},
		select: function (args) {
			this.selectCountry(args.countryName)
		}
	},
	template: `<div>
	<table style="width: 100%; border: 0; padding: 0; margin: 0">
		<tbody>
       	<tr>
           <td>
               <div class="traffic-map-box" id="traffic-map-box"></div>
           </td>
           <td style="width: 14em" v-if="!screenIsNarrow">
           		<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" @select="select"></traffic-map-box-table>
           </td>
       </tr>
       </tbody>
       <tbody v-if="screenIsNarrow">
		   <tr>
				<td colspan="2">
					<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" :v-screen-is-narrow="true" @select="select"></traffic-map-box-table>
				</td>
			</tr>
		</tbody>
   </table>
</div>`
})

Vue.component("traffic-map-box-table", {
	props: ["v-stats", "v-is-attack", "v-screen-is-narrow"],
	data: function () {
		return {
			stats: this.vStats,
			isAttack: this.vIsAttack
		}
	},
	methods: {
		select: function (countryName) {
			this.$emit("select", {countryName: countryName})
		}
	},
	template: `<div style="overflow-y: auto" :style="{'max-height':vScreenIsNarrow ? 'auto' : '16em'}" class="narrow-scrollbar">
	   <table class="ui table selectable">
		  <thead>
			<tr>
				<th colspan="2">国家/地区排行&nbsp; <tip-icon content="只有开启了统计的网站才会有记录。"></tip-icon></th>
			</tr>
		  </thead>
		   <tbody v-if="stats.length == 0">
			   <tr>
				   <td colspan="2">暂无数据</td>
			   </tr>
		   </tbody>
		   <tbody>
			   <tr v-for="(stat, index) in stats.slice(0, 10)">
				   <td @click.prevent="select(stat.name)" style="cursor: pointer" colspan="2">
					   <div class="ui progress bar" :class="{red: vIsAttack, blue:!vIsAttack}" style="margin-bottom: 0.3em">
						   <div class="bar" style="min-width: 0; height: 4px;" :style="{width: stat.percent + '%'}"></div>
					   </div>
					  <div>{{stat.name}}</div> 
					   <div><span class="grey">{{stat.percent}}% </span>
					   <span class="small grey" v-if="isAttack">{{stat.formattedCountAttackRequests}}</span>
					   <span class="small grey" v-if="!isAttack">（{{stat.formattedBytes}}）</span></div>
				   </td>
			   </tr>
		   </tbody>
	   </table>
   </div>`
})