Vue.component("columns-grid", {
	props: [],
	mounted: function () {
		this.columns = this.calculateColumns()

		let that = this
		window.addEventListener("resize", function () {
			that.columns = that.calculateColumns()
		})
	},
	data: function () {
		return {
			columns: "four"
		}
	},
	methods: {
		calculateColumns: function () {
			let w = window.innerWidth
			let columns = Math.floor(w / 250)
			if (columns == 0) {
				columns = 1
			}

			let columnElements = this.$el.getElementsByClassName("column")
			if (columnElements.length == 0) {
				return
			}
			let maxColumns = columnElements.length
			if (columns > maxColumns) {
				columns = maxColumns
			}

			// 添加右侧边框
			for (let index = 0; index < columnElements.length; index++) {
				let el = columnElements[index]
				el.className = el.className.replace("with-border", "")
				if (index % columns == columns - 1 || index == columnElements.length - 1 /** 最后一个 **/) {
					el.className += " with-border"
				}
			}

			switch (columns) {
				case 1:
					return "one"
				case 2:
					return "two"
				case 3:
					return "three"
				case 4:
					return "four"
				case 5:
					return "five"
				case 6:
					return "six"
				case 7:
					return "seven"
				case 8:
					return "eight"
				case 9:
					return "nine"
				case 10:
					return "ten"
				default:
					return "ten"
			}
		}
	},
	template: `<div :class="'ui ' + columns + ' columns grid counter-chart'">
	<slot></slot>
</div>`
})