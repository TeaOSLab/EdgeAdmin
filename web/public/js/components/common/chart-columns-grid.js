Vue.component("chart-columns-grid", {
	props: [],
	mounted: function () {
		this.columns = this.calculateColumns()

		let that = this
		window.addEventListener("resize", function () {
			that.columns = that.calculateColumns()
		})
	},
	updated: function () {
		let totalElements = this.$el.getElementsByClassName("column").length
		if (totalElements == this.totalElements) {
			return
		}
		this.totalElements = totalElements
		this.calculateColumns()
	},
	data: function () {
		return {
			columns: "four",
			totalElements: 0
		}
	},
	methods: {
		calculateColumns: function () {
			let w = window.innerWidth
			let columns = Math.floor(w / 500)
			if (columns == 0) {
				columns = 1
			}

			let columnElements = this.$el.getElementsByClassName("column")
			if (columnElements.length == 0) {
				return "one"
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
	template: `<div :class="'ui ' + columns + ' columns grid chart-grid'">
	<slot></slot>
</div>`
})