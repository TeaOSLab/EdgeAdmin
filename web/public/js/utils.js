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
	datepicker: function (element, callback) {
		// 加载
		if (typeof Pikaday == "undefined") {
			let that = this
			this.loadJS("/js/moment.min.js")
			this.loadJS("/js/pikaday.js", function () {
				that.datepicker(element, callback)
			})
			this.loadCSS("/js/pikaday.css")
			this.loadCSS("/js/pikaday.theme.css")
			this.loadCSS("/js/pikaday.triangle.css")

			return
		}


		if (typeof (element) == "string") {
			element = document.getElementById(element);
		}
		var year = new Date().getFullYear();
		var picker = new Pikaday({
			field: element,
			firstDay: 1,
			minDate: new Date(year - 1, 0, 1),
			maxDate: new Date(year + 10, 11, 31),
			yearRange: [year - 1, year + 10],
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
			}
		});
	},

	formatBytes: function (bytes) {
		bytes = Math.ceil(bytes);
		if (bytes < 1024) {
			return bytes + "B";
		}
		if (bytes < 1024 * 1024) {
			return (Math.round(bytes * 100 / 1024) / 100) + "K";
		}
		if (bytes < 1024 * 1024 * 1024) {
			return (Math.round(bytes * 100 / 1024 / 1024) / 100) + "M";
		}
		if (bytes < 1024 * 1024 * 1024 * 1024) {
			return (Math.round(bytes * 100 / 1024 / 1024 / 1024) / 100) + "G";
		}
		if (bytes < 1024 * 1024 * 1024 * 1024 * 1024) {
			return (Math.round(bytes * 100 / 1024 / 1024 / 1024 / 1024) / 100) + "T";
		}
		return (Math.round(bytes * 100 / 1024 / 1024 / 1024 / 1024 / 1024) / 100) + "P";
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
		let unit = ""
		if (max >= 1024 * 1024 * 1024 * 1024) {
			unit = "T"
			divider = 1024 * 1024 * 1024 * 1024
		} else if (max >= 1024 * 1024 * 1024) {
			unit = "G"
			divider = 1024 * 1024 * 1024
		} else if (max >= 1024 * 1024) {
			unit = "M"
			divider = 1024 * 1024
		} else if (max >= 1024) {
			unit = "K"
			divider = 1024
		}
		return {
			unit: unit,
			divider: divider
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
	popupFinish: function () {
		if (window.POPUP_CALLBACK != null) {
			window.POPUP_CALLBACK.apply(window, arguments);
		}
	},
	popupTip: function (html) {
		Swal.fire({
			html: '<div style="line-height: 1.7;text-align: left "><i class="icon question circle"></i>' + html + "</div>",
			width: "30em",
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
		var width = "20em";
		if (message.length > 30) {
			width = "30em";
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
						callback();
					});
				}
			}
		});
	},
	confirm: function (message, callback) {
		let width = "20em";
		if (message.length > 30) {
			width = "30em";
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
					callback.call(Tea.Vue);
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
						color: "#9DD3E8"
					},
					barWidth: "20em"
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
		let tooltipFunc = options.tooltip
		let axis = options.axis
		let valueFunc = options.value
		let max = options.max
		let interval = options.interval

		let chartBox = document.getElementById(chartId)
		if (chartBox == null) {
			return
		}
		let chart = this.initChart(chartBox)
		let option = {
			xAxis: {
				data: values.map(xFunc),
				axisLabel: {
					interval: interval
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
					return tooltipFunc.apply(this, [args, values])
				}
			},
			grid: {
				left: 40,
				top: 10,
				right: 20,
				bottom: 20
			},
			series: [
				{
					name: name,
					type: "line",
					data: values.map(valueFunc),
					itemStyle: {
						color: "#9DD3E8"
					},
					areaStyle: {}
				}
			],
			animation: true
		}
		chart.setOption(option)
		chart.resize()
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
	initChart: function (dom) {
		let instance = echarts.init(dom)
		window.addEventListener("resize", function () {
			instance.resize()
		})
		return instance
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
