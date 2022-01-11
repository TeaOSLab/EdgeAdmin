Vue.component("combo-box", {
	props: ["name", "title", "placeholder", "size", "v-items", "v-value"],
	data: function () {
		let items = this.vItems
		if (items == null || !(items instanceof Array)) {
			items = []
		}

		// 自动使用ID作为值
		items.forEach(function (v) {
			if (v.value == null) {
				v.value = v.id
			}
		})

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

		return {
			allItems: items,
			items: items.$copy(),
			selectedItem: selectedItem,
			keyword: "",
			visible: false,
			hideTimer: null,
			hoverIndex: 0
		}
	},
	methods: {
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
		changeKeyword: function () {
			this.hoverIndex = 0
			let keyword = this.keyword
			if (keyword.length == 0) {
				this.items = this.allItems.$copy()
				return
			}
			this.items = this.allItems.$copy().filter(function (v) {
				return teaweb.match(v.name, keyword)
			})
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
		}
	},
	template: `<div style="display: inline; z-index: 10; background: white">
	<!-- 搜索框 -->
	<div v-if="selectedItem == null">
		<input type="text" v-model="keyword" :placeholder="placeholder" :size="size" style="width: 11em" @input="changeKeyword" @focus="show" @blur="hide" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="searchBox" @keyup.down="downItem" @keyup.up="upItem"/>
	</div>
	
	<!-- 当前选中 -->
	<div v-if="selectedItem != null">
		<input type="hidden" :name="name" :value="selectedItem.value"/>
		<span class="ui label basic">{{title}}：{{selectedItem.name}}
			<a href="" title="清除" @click.prevent="reset"><i class="icon remove small"></i></a>
		</span>
	</div>
	
	<!-- 菜单 -->
	<div v-if="selectedItem == null && items.length > 0 && visible">
		<div class="ui menu vertical small narrow-scrollbar" style="width: 11em; max-height: 17em; overflow-y: auto; position: absolute; border: rgba(129, 177, 210, 0.81) 1px solid; border-top: 0; z-index: 100">
			<a href="" v-for="(item, index) in items" ref="itemRef" class="item" :class="{active: index == hoverIndex, blue: index == hoverIndex}" @click.prevent="selectItem(item)" style="line-height: 1.4">{{item.name}}</a>
		</div>
	</div>
</div>`
})