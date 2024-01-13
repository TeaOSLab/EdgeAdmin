window.teaweb = {
	set: function (key, value) {
		localStorage.setItem(key, JSON.stringify(value));
	},
	get: function (key) {
		var item = localStorage.getItem(key);
		if (item == null || item.length == 0) {
			return null;
		}

		return JSON.parse(item);
	},
	getString: function (key) {
		var value = this.get(key);
		if (typeof (value) == "string") {
			return value;
		}
		return "";
	},
	getBool: function (key) {
		return Boolean(this.get(key));
	},
	remove: function (key) {
		localStorage.removeItem(key)
	},
	match: function (source, keyword) {
		if (source == null) {
			return false;
		}
		if (keyword == null) {
			return true;
		}
		source = source.trim();
		keyword = keyword.trim();
		if (keyword.length == 0) {
			return true;
		}
		if (source.length == 0) {
			return false;
		}
		var pieces = keyword.split(/\s+/);
		for (var i = 0; i < pieces.length; i++) {
			var pattern = pieces[i];
			pattern = pattern.replace(/(\+|\*|\?|[|]|{|}|\||\\|\(|\)|\.)/g, "\\$1");
			var reg = new RegExp(pattern, "i");
			if (!reg.test(source)) {
				return false;
			}
		}
		return true;
	},
	clone: function (source) {
		let s = JSON.stringify(source)
		return JSON.parse(s)
	},

	loadJS: function (file, callback) {
		let element = document.createElement("script")
		element.setAttribute("type", "text/javascript")
		element.setAttribute("src", file)
		if (typeof callback == "function") {
			element.addEventListener("load", callback)
		}
		document.head.append(element)
	},
	loadCSS: function (file, callback) {
		let element = document.createElement("link")
		element.setAttribute("rel", "stylesheet")
		element.setAttribute("type", "text/css")
		element.setAttribute("href", file)
		if (typeof callback == "function") {
			element.addEventListener("load", callback)
		}
		document.head.append(element)
	},
	datepicker: function (element, callback, bottomLeft) {
		// 加载
		if (typeof Pikaday == "undefined") {
			let that = this
			this.loadJS("/js/moment.min.js", function () {
				that.loadJS("/js/pikaday.js", function () {
					that.datepicker(element, callback, bottomLeft)
				})
			})
			this.loadCSS("/js/pikaday.css")
			this.loadCSS("/js/pikaday.theme.css")
			this.loadCSS("/js/pikaday.triangle.css")

			return
		}

		if (typeof (element) == "string") {
			element = document.getElementById(element);
		}
		let year = new Date().getFullYear();
		let picker = new Pikaday({
			field: element,
			firstDay: 1,
			minDate: new Date(year - 1, 0, 1),
			maxDate: new Date(year + 20, 11, 31),
			yearRange: [year - 1, year + 20],
			format: "YYYY-MM-DD",
			i18n: {
				previousMonth: '上月',
				nextMonth: '下月',
				months: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
				weekdays: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
				weekdaysShort: ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
			},
			theme: 'triangle-theme',
			onSelect: function () {
				if (typeof (callback) == "function") {
					callback.call(Tea.Vue, picker.toString());
				}
			},
			reposition: !bottomLeft
		})

		element.picker = picker
	},
	formatBytes: function (bytes) {
		bytes = Math.ceil(bytes);
		if (bytes < Math.pow(1024, 1)) {
			return bytes + "B";
		}
		if (bytes < Math.pow(1024, 2)) {
			return (Math.round(bytes * 100 / Math.pow(1024, 1)) / 100) + "KiB";
		}
		if (bytes < Math.pow(1024, 3)) {
			return (Math.round(bytes * 100 / Math.pow(1024, 2)) / 100) + "MiB";
		}
		if (bytes < Math.pow(1024, 4)) {
			return (Math.round(bytes * 100 / Math.pow(1024, 3)) / 100) + "GiB";
		}
		if (bytes < Math.pow(1024, 5)) {
			return (Math.round(bytes * 100 / Math.pow(1024, 4)) / 100) + "TiB";
		}
		if (bytes < Math.pow(1024, 6)) {
			return (Math.round(bytes * 100 / Math.pow(1024, 5)) / 100) + "PiB";
		}
		return (Math.round(bytes * 100 / Math.pow(1024, 6)) / 100) + "EiB";
	},
	formatBits: function (bits, decimal) {
		bits = Math.ceil(bits);
		let div = 10000
		switch (decimal) {
			case 1:
				div = 10
				break
			case 2:
				div = 100
				break
			case 3:
				div = 1000
				break
			case 4:
				div = 10000
				break
		}
		if (bits < Math.pow(1024, 1)) {
			return bits + "bps";
		}
		if (bits < Math.pow(1024, 2)) {
			return (Math.round(bits * div / Math.pow(1024, 1)) / div) + "Kbps";
		}
		if (bits < Math.pow(1024, 3)) {
			return (Math.round(bits * div / Math.pow(1024, 2)) / div) + "Mbps";
		}
		if (bits < Math.pow(1024, 4)) {
			return (Math.round(bits * div / Math.pow(1024, 3)) / div) + "Gbps";
		}
		if (bits < Math.pow(1024, 5)) {
			return (Math.round(bits * div / Math.pow(1024, 4)) / div) + "Tbps";
		}
		if (bits < Math.pow(1024, 6)) {
			return (Math.round(bits * div / Math.pow(1024, 5)) / div) + "Pbps";
		}
		return (Math.round(bits * div / Math.pow(1024, 6)) / div) + "Ebps";
	},
	formatNumber: function (x) {
		if (x == null) {
			return "null"
		}
		let s = x.toString()
		let dotIndex = s.indexOf(".")
		if (dotIndex >= 0) {
			return this.formatNumber(s.substring(0, dotIndex)) + "." + s.substring(dotIndex + 1)
		}

		if (s.length <= 3) {
			return s;
		}
		let result = []
		for (let i = 0; i < Math.floor(s.length / 3); i++) {
			let start = s.length - (i + 1) * 3
			result.push(s.substring(start, start + 3))
		}
		if (s.length % 3 != 0) {
			result.push(s.substring(0, s.length % 3))
		}
		return result.reverse().join(", ")
	},
	formatCount: function (x) {
		let unit = ""
		let divider = ""
		if (x >= 1000 * 1000 * 1000) {
			unit = "B"
			divider = 1000 * 1000 * 1000
		} else if (x >= 1000 * 1000) {
			unit = "M"
			divider = 1000 * 1000
		} else if (x >= 1000) {
			unit = "K"
			divider = 1000
		}
		if (unit.length == 0) {
			return x.toString()
		}
		return (Math.round(x * 100 / divider) / 100) + unit
	},
	bytesAxis: function (stats, countFunc) {
		let max = Math.max.apply(this, stats.map(countFunc))
		let divider = 1
		let unit = "B"
		if (max >= Math.pow(1024, 6)) {
			unit = "E"
			divider = Math.pow(1024, 6)
		} else if (max >= Math.pow(1024, 5)) {
			unit = "P"
			divider = Math.pow(1024, 5)
		} else if (max >= Math.pow(1024, 4)) {
			unit = "T"
			divider = Math.pow(1024, 4)
		} else if (max >= Math.pow(1024, 3)) {
			unit = "G"
			divider = Math.pow(1024, 3)
		} else if (max >= Math.pow(1024, 2)) {
			unit = "M"
			divider = Math.pow(1024, 2)
		} else if (max >= Math.pow(1024, 1)) {
			unit = "K"
			divider = Math.pow(1024, 1)
		}
		return {
			unit: unit,
			divider: divider
		}
	},
	bitsAxis: function (stats, countFunc) {
		let axis = this.bytesAxis(stats, countFunc)
		let unit = axis.unit
		if (unit == "B") {
			unit = "bps"
		} else {
			unit += "bps"
		}
		return {
			unit: unit,
			divider: axis.divider
		}
	},
	countAxis: function (stats, countFunc) {
		let max = Math.max.apply(this, stats.map(countFunc))
		let divider = 1
		let unit = ""
		if (max >= 1000 * 1000 * 1000) {
			unit = "B"
			divider = 1000 * 1000 * 1000
		} else if (max >= 1000 * 1000) {
			unit = "M"
			divider = 1000 * 1000
		} else if (max >= 1000) {
			unit = "K"
			divider = 1000
		}
		return {
			unit: unit,
			divider: divider,
			max: max
		}
	},
	splitFormat: function (s) {
		let matchResult = s.match(/^([0-9.]+)([a-zA-Z]+)$/)
		let size = parseFloat(matchResult[1])
		let unit = matchResult[2]
		return [size, unit]
	},
	popup: function (url, options) {
		if (url != null && url.length > 0 && url.substring(0, 1) == '.') {
			url = Tea.url(url)
		}

		if (options == null) {
			options = {};
		}
		var width = "40em";
		var height = "20em";
		window.POPUP_CALLBACK = function () {
			Swal.close();
		};

		if (options["width"] != null) {
			width = options["width"];
		}
		if (options["height"] != null) {
			height = options["height"];
		}
		if (typeof (options["callback"]) == "function") {
			window.POPUP_CALLBACK = function () {
				Swal.close();
				options["callback"].apply(Tea.Vue, arguments);
			};
		}

		Swal.fire({
			html: '<iframe src="' + url + '#popup-' + width + '" style="border:0; width: 100%; height:' + height + '"></iframe>',
			width: width,
			padding: "0.5em",
			showConfirmButton: false,
			showCloseButton: true,
			focusConfirm: false,
			onClose: function (popup) {
				if (typeof (options["onClose"]) == "function") {
					options["onClose"].apply(Tea.Vue, arguments)
				}
			}
		});
	},
	popupSuccess: function (url, width, height) {
		let options = {}
		if (width != null) {
			options["width"] = width
		}
		if (height != null) {
			options["height"] = height
		}
		options["callback"] = function () {
			teaweb.success("保存成功", function () {
				teaweb.reload()
			})
		}
		this.popup(url, options)
	},
	popupFinish: function () {
		if (window.POPUP_CALLBACK != null) {
			window.POPUP_CALLBACK.apply(window, arguments);
		}
	},
	popupTip: function (html) {
		Swal.fire({
			html: '<div style="line-height: 1.7;text-align: left "><i class="icon question circle"></i>' + html + "</div>",
			width: "34em",
			padding: "4em",
			showConfirmButton: false,
			showCloseButton: true,
			focusConfirm: false
		});
	},
	isPopup: function () {
		var hash = window.location.hash;
		return hash != null && hash.startsWith("#popup");
	},
	closePopup: function () {
		if (this.isPopup()) {
			window.parent.Swal.close();
		}
	},
	Swal: function () {
		return this.isPopup() ? window.parent.Swal : window.Swal;
	},
	hasPopup: function () {
		return document.getElementsByClassName("swal2-container").length > 0
	},
	success: function (message, callback) {
		var width = "20em";
		if (message.length > 30) {
			width = "30em";
		}

		let config = {
			confirmButtonText: "确定",
			buttonsStyling: false,
			icon: "success",
			customClass: {
				closeButton: "ui button",
				cancelButton: "ui button",
				confirmButton: "ui button primary"
			},
			width: width,
			onAfterClose: function () {
				if (typeof (callback) == "function") {
					setTimeout(function () {
						callback();
					});
				} else if (typeof (callback) == "string") {
					window.location = callback
				}
			}
		}

		if (message.startsWith("html:")) {
			config.html = message.substring(5)
		} else {
			config.text = message
		}

		Swal.fire(config);
	},
	toast: function (message, timeout, callback) {
		if (timeout == null) {
			timeout = 2000
		}
		var width = "20em";
		if (message.length > 30) {
			width = "30em";
		}
		Swal.fire({
			text: message,
			icon: "info",
			width: width,
			timer: timeout,
			showConfirmButton: false,
			onAfterClose: function () {
				if (typeof callback == "function") {
					callback()
				}
			}
		});
	},
	successToast: function (message, timeout, callback) {
		if (timeout == null) {
			timeout = 2000
		}
		var width = "20em";
		if (message.length > 30) {
			width = "30em";
		}
		Swal.fire({
			text: message,
			icon: "success",
			width: width,
			timer: timeout,
			showConfirmButton: false,
			onAfterClose: function () {
				if (typeof callback == "function") {
					callback()
				}
			}
		});
	},
	successRefresh: function (message) {
		teaweb.success(message, function () {
			teaweb.reload()
		})
	},
	warn: function (message, callback) {
		var width = "20em"
		if (message.length > 30) {
			width = "30em"
		}
		Swal.fire({
			text: message,
			confirmButtonText: "确定",
			buttonsStyling: false,
			customClass: {
				closeButton: "ui button",
				cancelButton: "ui button",
				confirmButton: "ui button primary"
			},
			icon: "warning",
			width: width,
			onAfterClose: function () {
				if (typeof (callback) == "function") {
					setTimeout(function () {
						callback()
					})
				}
			}
		})
	},
	confirm: function (message, callback) {
		let width = "20em"
		if (message.length > 30) {
			width = "30em"
		}
		let config = {
			confirmButtonText: "确定",
			cancelButtonText: "取消",
			showCancelButton: true,
			showCloseButton: false,
			buttonsStyling: false,
			customClass: {
				closeButton: "ui button",
				cancelButton: "ui button",
				confirmButton: "ui button primary"
			},
			icon: "warning",
			width: width,
			preConfirm: function () {
				if (typeof (callback) == "function") {
					callback.call(Tea.Vue)
				}
			}
		}
		if (message.startsWith("html:")) {
			config.html = message.substring(5)
		} else {
			config.text = message
		}
		Swal.fire(config);
	},
	reload: function () {
		window.location.reload()
	},
	renderBarChart: function (options) {
		let chartId = options.id
		if (chartId == null || chartId.length == 0) {
			throw new Error("'options.id' should not be empty")
		}

		let name = options.name
		let values = options.values
		if (values == null || !(values instanceof Array)) {
			throw new Error("'options.values' should be array")
		}

		let xFunc = options.x
		if (typeof (xFunc) != "function") {
			throw new Error("'options.x' should be a function")
		}

		let tooltipFunc = options.tooltip
		if (typeof (tooltipFunc) != "function") {
			throw new Error("'options.tooltip' should be a function")
		}

		let axis = options.axis
		if (axis == null) {
			axis = {unit: "", count: 1}
		}
		let valueFunc = options.value
		if (typeof (valueFunc) != "function") {
			throw new Error("'options.value' should be a function")
		}
		let click = options.click

		let bottom = 24
		let rotate = 0
		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = this.initChart(chartBox)
		let result = this.xRotation(chart, values.map(xFunc))
		if (result != null) {
			bottom = result[0]
			rotate = result[1]
		}
		let option = {
			xAxis: {
				data: values.map(xFunc),
				axisLabel: {
					interval: 0,
					rotate: rotate
				}
			},
			yAxis: {
				axisLabel: {
					formatter: function (value) {
						return value + axis.unit
					}
				}
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: function (args) {
					return tooltipFunc.apply(this, [args, values])
				}
			},
			grid: {
				left: 40,
				top: 10,
				right: 20,
				bottom: bottom
			},
			series: [
				{
					name: name,
					type: "bar",
					data: values.map(valueFunc),
					itemStyle: {
						color: this.DefaultChartColor
					},
					barWidth: "10em",
					areaStyle: {}
				}
			],
			animation: true,
		}
		chart.setOption(option)
		if (click != null) {
			chart.on("click", function (args) {
				click.call(this, args, values)
			})
		}
		chart.resize()
	},
	renderLineChart: function (options) {
		let chartId = options.id
		let name = options.name
		let values = options.values
		let xFunc = options.x
		let xColorFunc = options.xColor
		let tooltipFunc = options.tooltip
		let axis = options.axis
		let valueFunc = options.value
		let max = options.max
		let interval = options.interval
		let left = options.left
		if (typeof left != "number") {
			left = 0
		}

		let right = options.right
		if (typeof right != "number") {
			right = 0
		}

		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			console.error("chart id '" + chartId + "' not found")
			return
		}
		let chart = this.initChart(chartBox, options.cache)
		let option = {
			xAxis: {
				data: values.map(xFunc),
				axisLabel: {
					interval: interval,
					textStyle: {
						color: xColorFunc
					}
				}
			},
			yAxis: {
				axisLabel: {
					formatter: function (value) {
						return value + axis.unit
					}
				},
				max: max
			},
			tooltip: {
				show: true,
				trigger: "item",
				formatter: function (args) {
					if (tooltipFunc != null) {
						return tooltipFunc.apply(this, [args, values])
					}
					return null
				}
			},
			grid: {
				left: 40 + left,
				top: 10,
				right: 20 + right,
				bottom: 20
			},
			series: [
				{
					name: name,
					type: "line",
					data: values.map(valueFunc),
					itemStyle: {
						color: this.DefaultChartColor
					},
					areaStyle: {},
					smooth: true,
					markLine: options.markLine
				}
			],
			animation: true,
			smooth: true
		}
		chart.setOption(option)
		chart.resize()
	},
	renderGaugeChart: function (options) {
		let chartId = options.id
		let name = options.name // 标题
		let min = options.min // 最小值
		let max = options.max // 最大值
		let value = options.value // 当前值
		let unit = options.unit // 单位
		let detail = options.detail // 说明文字
		let color = options.color // 颜色
		let startAngle = options.startAngle
		if (startAngle == null) {
			startAngle = 225
		}
		let endAngle = options.endAngle
		if (endAngle == null) {
			endAngle = -45
		}

		color = this.chartColor(color)

		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = this.initChart(chartBox)

		let option = {
			textStyle: {
				fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
			},
			color: color,
			title: (name != null && name.length > 0) ? {
				text: name,
				top: 1,
				bottom: 0,
				x: "center",
				textStyle: {
					fontSize: 12,
					fontWeight: "bold",
					fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
				}
			} : null,
			legend: {
				data: [""]
			},
			xAxis: {
				data: []
			},
			yAxis: {},
			series: [{
				name: "",
				type: "gauge",
				min: min,
				max: max,
				startAngle: startAngle,
				endAngle: endAngle,

				data: [
					{
						"name": "",// 不显示名称
						"value": Math.round(value * 100) / 100
					}
				],
				radius: "100%",
				center: ["50%", (name != null && name.length > 0) ? "60%" : "50%"],

				splitNumber: 5,
				splitLine: {
					length: 4
				},

				axisLine: {
					lineStyle: {
						width: 4
					}
				},
				axisTick: {
					show: true,
					length: 2
				},
				axisLabel: {
					formatter: function (v) {
						return v;
					},
					textStyle: {
						fontSize: 8
					}
				},
				progress: {
					show: true,
					width: 4
				},
				detail: {
					formatter: function (v) {
						return unit;
					},
					textStyle: {
						fontSize: 12,
						fontWeight: "normal",
						fontFamily: "Arial,Helvetica,sans-serif",
						color: "grey"
						//lineHeight: 16
					},
					valueAnimation: true
				},

				pointer: {
					width: 2
				}
			}],

			grid: {
				left: -2,
				right: 0,
				bottom: 0,
				top: 0
			},
			axisPointer: {
				show: false
			},
			tooltip: {
				formatter: 'X:{b0} Y:{c0}',
				show: false
			},
			animation: true
		};

		chart.setOption(option)
	},

	renderPercentChart: function (options) {
		let chartId = options.id
		let color = this.chartColor(options.color)
		let value = options.value
		let name = options.name
		let total = options.total
		if (total == null) {
			total = 100
		}
		let unit = options.unit
		if (unit == null) {
			unit = ""
		}

		let max = options.max
		if (max != null && max <= value) {
			max = null
		}
		let maxColor = this.chartColor(options.maxColor)
		let maxName = options.maxName

		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = this.initChart(chartBox)

		let option = {
			tooltip: {
				formatter: "{a} <br/>{b} : {c}" + unit
			},
			series: [
				{
					name: name,
					max: total,
					type: "gauge",
					radius: "100%",
					detail: {
						formatter: "{value}",
						show: false,
						valueAnimation: true
					},
					data: [
						{
							value: value,
							name: name
						}
					],
					pointer: {
						show: false
					},
					splitLine: {
						show: false
					},
					axisTick: {
						show: false
					},
					axisLine: {
						show: true,
						lineStyle: {
							width: 4
						}
					},
					progress: {
						show: true,
						width: 4,
						itemStyle: {
							color: color
						}
					},
					splitNumber: {
						show: false
					},
					title: {
						show: false
					},
					startAngle: 270,
					endAngle: -90
				}
			]
		}

		if (max != null) {
			option.series.push({
				name: maxName,
				max: total,
				type: "gauge",
				radius: "100%",
				detail: {
					formatter: "{value}",
					show: false,
					valueAnimation: true
				},
				data: [
					{
						value: max,
						name: maxName
					}
				],
				pointer: {
					show: false
				},
				splitLine: {
					show: false
				},
				axisTick: {
					show: false
				},
				axisLine: {
					show: true,
					lineStyle: {
						width: 4
					}
				},
				progress: {
					show: true,
					width: 4,
					itemStyle: {
						color: maxColor,
						opacity: 0.3
					}
				},
				splitNumber: {
					show: false
				},
				title: {
					show: false
				},
				startAngle: 270,
				endAngle: -90
			})
		}

		chart.setOption(option)
	},

	xRotation: function (chart, names) {
		let chartWidth = chart.getWidth()
		let width = 0
		names.forEach(function (name) {
			width += name.length * 10
		})
		if (width <= chartWidth) {
			return null
		}

		return [40, -20]
	},
	chartMap: {}, // dom id => chart
	initChart: function (dom, cache) {
		if (typeof (cache) != "boolean") {
			cache = true
		}

		let domId = dom.getAttribute("id")
		if (domId != null && domId.length > 0 && typeof (this.chartMap[domId]) == "object") {
			return this.chartMap[domId]
		}
		let instance = echarts.init(dom)
		window.addEventListener("resize", function () {
			instance.resize()
		})
		if (cache) {
			this.chartMap[domId] = instance
		}
		return instance
	},
	encodeHTML: function (s) {
		s = s.replace(/&/g, "&amp;")
		s = s.replace(/</g, "&lt;")
		s = s.replace(/>/g, "&gt;")
		s = s.replace(/"/g, "&quot;")
		return s
	},
	chartColor: function (color) {
		// old blue: #5470c6
		if (color == null || color.length == 0) {
			color = "#5470c6"
		}

		if (color == "red") {
			color = "#ee6666"
		}
		if (color == "yellow") {
			color = "#fac858"
		}
		if (color == "blue") {
			color = "#5470c6"
		}
		if (color == "green") {
			color = "#3ba272"
		}
		return color
	},
	DefaultChartColor: "#9DD3E8",
	validateIP: function (ip) {
		if (typeof ip != "string") {
			return false
		}

		if (ip.length == 0) {
			return false
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
	playAlert: function () {
		let audioBox = document.createElement("AUDIO")
		audioBox.setAttribute("control", "")
		audioBox.setAttribute("autoplay", "")
		audioBox.innerHTML = "<source src=\"/audios/alert.ogg\" type=\"audio/ogg\"/>";
		document.body.appendChild(audioBox);
		audioBox.play().then(function () {
			setTimeout(function () {
				document.body.removeChild(audioBox);
			}, 2000);
		}).catch(function (e) {
			console.log(e.message);
		})
	},
	convertSizeCapacityToBytes: function (c) {
		if (c == null) {
			return 0
		}
		switch (c.unit) {
			case "byte":
				return c.count
			case "kb":
				return c.count * 1024
			case "mb":
				return c.count * Math.pow(1024, 2)
			case "gb":
				return c.count * Math.pow(1024, 3)
			case "tb":
				return c.count * Math.pow(1024, 4)
			case "pb":
				return c.count * Math.pow(1024, 5)
			case "eb":
				return c.count * Math.pow(1024, 6)
		}
		return 0
	},
	compareSizeCapacity: function (c1, c2) {
		let b1 = this.convertSizeCapacityToBytes(c1)
		let b2 = this.convertSizeCapacityToBytes(c2)
		if (b1 == b2) {
			return 0
		}
		if (b1 > b2) {
			return 1
		}
		return -1
	},
	convertSizeCapacityToString: function (c) {
		if (c == null || c.count == null || c.unit == null || c.unit.length == 0) {
			return ""
		}
		if (c.unit == "byte") {
			return c.count + "B"
		}
		return c.count + c.unit[0].toUpperCase() + "i" + c.unit.substring(1).toUpperCase()
	}
}

String.prototype.quoteIP = function () {
	let ip = this.toString()
	if (ip.length == 0) {
		return ""
	}
	if (ip.indexOf(":") < 0) {
		return ip
	}
	if (ip.substring(0, 1) == "[") {
		return ip
	}
	return "[" + ip + "]"
}
