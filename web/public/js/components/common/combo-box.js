Vue.component("combo-box", {
	// data-url 和 data-key 成对出现
	props: [
		"name", "title", "placeholder", "size", "v-items", "v-value",
		"data-url", // 数据源URL
		"data-key", // 数据源中数据的键名
		"data-search", // 是否启用动态搜索，如果值为on或true，则表示启用
		"width"
	],
	mounted: function () {
		if (this.dataURL.length > 0) {
			this.search("")
		}

		// 设定菜单宽度
		let searchBox = this.$refs.searchBox
		if (searchBox != null) {
			let inputWidth = searchBox.offsetWidth
			if (inputWidth != null && inputWidth > 0) {
				this.$refs.menu.style.width = inputWidth + "px"
			} else if (this.styleWidth.length > 0) {
				this.$refs.menu.style.width = this.styleWidth
			}
		}
	},
	data: function () {
		let items = this.vItems
		if (items == null || !(items instanceof Array)) {
			items = []
		}
		items = this.formatItems(items)

		// 当前选中项
		let selectedItem = null
		if (this.vValue != null) {
			let that = this
			items.forEach(function (v) {
				if (v.value == that.vValue) {
					selectedItem = v
				}
			})
		}

		let width = this.width
		if (width == null || width.length == 0) {
			width = "11em"
		} else {
			if (/\d+$/.test(width)) {
				width += "em"
			}
		}

		// data url
		let dataURL = ""
		if (typeof this.dataUrl == "string" && this.dataUrl.length > 0) {
			dataURL = this.dataUrl
		}

		return {
			allItems: items, // 原始的所有的items
			items: items.$copy(), // 候选的items
			selectedItem: selectedItem, // 选中的item
			keyword: "",
			visible: false,
			hideTimer: null,
			hoverIndex: 0,
			styleWidth: width,

			isInitial: true,
			dataURL: dataURL,
			urlRequestId: 0 // 记录URL请求ID，防止并行冲突
		}
	},
	methods: {
		search: function (keyword) {
			// 从URL中获取选项数据
			let dataUrl = this.dataURL
			let dataKey = this.dataKey
			let that = this

			let requestId = Math.random()
			this.urlRequestId = requestId

			Tea.action(dataUrl)
				.params({
					keyword: (keyword == null) ? "" : keyword
				})
				.post()
				.success(function (resp) {
					if (requestId != that.urlRequestId) {
						return
					}

					if (resp.data != null) {
						if (typeof (resp.data[dataKey]) == "object") {
							let items = that.formatItems(resp.data[dataKey])
							that.allItems = items
							that.items = items.$copy()

							if (that.isInitial) {
								that.isInitial = false
								if (that.vValue != null) {
									items.forEach(function (v) {
										if (v.value == that.vValue) {
											that.selectedItem = v
										}
									})
								}
							}
						}
					}
				})
		},
		formatItems: function (items) {
			items.forEach(function (v) {
				if (v.value == null) {
					v.value = v.id
				}
			})
			return items
		},
		reset: function () {
			this.selectedItem = null
			this.change()
			this.hoverIndex = 0

			let that = this
			setTimeout(function () {
				if (that.$refs.searchBox) {
					that.$refs.searchBox.focus()
				}
			})
		},
		clear: function () {
			this.selectedItem = null
			this.change()
			this.hoverIndex = 0
		},
		changeKeyword: function () {
			let shouldSearch = this.dataURL.length > 0 && (this.dataSearch == "on" || this.dataSearch == "true")

			this.hoverIndex = 0
			let keyword = this.keyword
			if (keyword.length == 0) {
				if (shouldSearch) {
					this.search(keyword)
				} else {
					this.items = this.allItems.$copy()
				}
				return
			}


			if (shouldSearch) {
				this.search(keyword)
			} else {
				this.items = this.allItems.$copy().filter(function (v) {
					if (v.fullname != null && v.fullname.length > 0 && teaweb.match(v.fullname, keyword)) {
						return true
					}
					return teaweb.match(v.name, keyword)
				})
			}
		},
		selectItem: function (item) {
			this.selectedItem = item
			this.change()
			this.hoverIndex = 0
			this.keyword = ""
			this.changeKeyword()
		},
		confirm: function () {
			if (this.items.length > this.hoverIndex) {
				this.selectItem(this.items[this.hoverIndex])
			}
		},
		show: function () {
			this.visible = true

			// 不要重置hoverIndex，以便焦点可以在输入框和可选项之间切换
		},
		hide: function () {
			let that = this
			this.hideTimer = setTimeout(function () {
				that.visible = false
			}, 500)
		},
		downItem: function () {
			this.hoverIndex++
			if (this.hoverIndex > this.items.length - 1) {
				this.hoverIndex = 0
			}
			this.focusItem()
		},
		upItem: function () {
			this.hoverIndex--
			if (this.hoverIndex < 0) {
				this.hoverIndex = 0
			}
			this.focusItem()
		},
		focusItem: function () {
			if (this.hoverIndex < this.items.length) {
				this.$refs.itemRef[this.hoverIndex].focus()
				let that = this
				setTimeout(function () {
					that.$refs.searchBox.focus()
					if (that.hideTimer != null) {
						clearTimeout(that.hideTimer)
						that.hideTimer = null
					}
				})
			}
		},
		change: function () {
			this.$emit("change", this.selectedItem)

			let that = this
			setTimeout(function () {
				if (that.$refs.selectedLabel != null) {
					that.$refs.selectedLabel.focus()
				}
			})
		},
		submitForm: function (event) {
			if (event.target.tagName != "A") {
				return
			}
			let parentBox = this.$refs.selectedLabel.parentNode
			while (true) {
				parentBox = parentBox.parentNode
				if (parentBox == null || parentBox.tagName == "BODY") {
					return
				}
				if (parentBox.tagName == "FORM") {
					parentBox.submit()
					break
				}
			}
		},

		setDataURL: function (dataURL) {
			this.dataURL = dataURL
		},
		reloadData: function () {
			this.search("")
		}
	},
	template: `<div style="display: inline; z-index: 10; background: white" class="combo-box">
	<!-- 搜索框 -->
	<div v-if="selectedItem == null">
		<input type="text" v-model="keyword" :placeholder="placeholder" :size="size" :style="{'width': styleWidth}"  @input="changeKeyword" @focus="show" @blur="hide" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="searchBox" @keydown.down.prevent="downItem" @keydown.up.prevent="upItem"/>
	</div>
	
	<!-- 当前选中 -->
	<div v-if="selectedItem != null">
		<input type="hidden" :name="name" :value="selectedItem.value"/>
		<span class="ui label basic" style="line-height: 1.4; font-weight: normal; font-size: 1em" ref="selectedLabel"><span><span v-if="title != null && title.length > 0">{{title}}：</span>{{selectedItem.name}}</span>
			<a href="" title="清除" @click.prevent="reset"><i class="icon remove small"></i></a>
		</span>
	</div>
	
	<!-- 菜单 -->
	<div v-show="selectedItem == null && items.length > 0 && visible">
		<div class="ui menu vertical small narrow-scrollbar" ref="menu">
			<a href="" v-for="(item, index) in items" ref="itemRef" class="item" :class="{active: index == hoverIndex, blue: index == hoverIndex}" @click.prevent="selectItem(item)" style="line-height: 1.4">
				<span v-if="item.fullname != null && item.fullname.length > 0">{{item.fullname}}</span>
				<span v-else>{{item.name}}</span>
			</a>
		</div>
	</div>
</div>`
})