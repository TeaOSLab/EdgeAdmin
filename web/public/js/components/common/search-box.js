Vue.component("search-box", {
	props: ["placeholder", "width"],
	data: function () {
		let width = this.width
		if (width == null) {
			width = "10em"
		}
		return {
			realWidth: width,
			realValue: ""
		}
	},
	methods: {
		onInput: function () {
			this.$emit("input", { value: this.realValue})
			this.$emit("change", { value: this.realValue})
		},
		clearValue: function () {
			this.realValue = ""
			this.focus()
			this.onInput()
		},
		focus: function () {
			this.$refs.valueRef.focus()
		}
	},
	template: `<div>
	<div class="ui input small" :class="{'right labeled': realValue.length > 0}">
		<input type="text" :placeholder="placeholder" :style="{width: realWidth}" @input="onInput" v-model="realValue" ref="valueRef"/>
		<a href="" class="ui label blue" v-if="realValue.length > 0" @click.prevent="clearValue" style="padding-right: 0"><i class="icon remove"></i></a>
	</div>
</div>`
})