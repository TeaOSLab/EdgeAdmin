Vue.component("traffic-map-box", {
	props: ["v-stats", "v-is-attack"],
	mounted: function () {
		this.render()
	},
	data: function () {
		let maxPercent = 0
		this.vStats.forEach(function (v) {
			let percent = parseFloat(v.percent)
			if (percent > maxPercent) {
				maxPercent = percent
			}
		})

		if (maxPercent < 100) {
			maxPercent *= 1.2 // 不要让某一项100%
		}

		return {
			stats: this.vStats,
			chart: null,
			minOpacity: 0.2,
			maxPercent: maxPercent
		}
	},
	methods: {
		render: function () {
			this.chart = teaweb.initChart(document.getElementById("traffic-map-box"));
			let that = this
			this.chart.setOption({
				backgroundColor: "white",
				grid: {top: 0, bottom: 0, left: 0, right: 0},
				roam: true,
				tooltip: {
					trigger: "item"
				},
				series: [{
					type: "map",
					map: "world",
					zoom: 1.2,
					selectedMode: false,
					itemStyle: {
						areaColor: "#E9F0F9",
						borderColor: "#999"
					},
					emphasis: {
						itemStyle: {
							areaColor: "#8B9BD3",
							opacity: 1.0
						}
					},
					//select: {itemStyle:{ areaColor: "#8B9BD3", opacity: 0.8 }},
					label: {
						show: true,
						formatter: function (args) {
							return ""
						}
					},
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
								return name + "<br/>流量：" + stat.formattedBytes + "<br/>请求数：" + teaweb.formatNumber(stat.countRequests) + "<br/>流量占比：" + stat.percent + "%"
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
						let isAttack = this.vIsAttack
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
								}
							}
						}
					}),
					nameMap: window.WorldCountriesMap
				}]
			})
			this.chart.resize()
		},
		select: function (countryName) {
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
						return
					}
					v.isSelected = true
					opacity *= 3
					if (opacity > 1) {
						opacity = 1
					}

					// 至少是0.5，让用户能够看清
					if (opacity < 0.5) {
						opacity = 0.5
					}
					v.itemStyle.opacity = opacity
				} else {
					v.itemStyle.opacity = opacity
					v.isSelected = false
				}
			})
			this.chart.setOption(option)
		}
	},
	template: `<div>
<table style="width: 100%; border: 0; padding: 0; margin: 0">
       <tr>
           <td>
               <div class="traffic-map-box" id="traffic-map-box" ></div>
           </td>
           <td style="width: 14em">
           		<div style="overflow-y: auto; height: 16em">
				   <table class="ui table selectable">
					  <thead>
						<tr>
							<th colspan="2">国家/地区排行</th>
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
								   <div class="ui progress bar blue" style="margin-bottom: 0.3em">
									   <div class="bar" style="min-width: 0; height: 4px;" :style="{width: stat.percent + '%'}"></div>
								   </div>
								  <div>{{stat.name}}</div> 
								   <div><span class="grey">{{stat.percent}}% </span><span class="small grey">（{{stat.formattedBytes}}）</span></div>
							   </td>
						   </tr>
					   </tbody>
				   </table>
               </div>
           </td>
       </tr>
   </table>
</div>`
})