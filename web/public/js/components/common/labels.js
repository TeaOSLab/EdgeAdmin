// 启用状态标签
Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="green">已启用</span><span v-if="!vIsOn" class="red">已停用</span></div>'
})

// 文字代码标签
Vue.component("code-label", {
	methods: {
		click: function (args) {
			this.$emit("click", args)
		}
	},
	template: `<span class="ui label basic small" style="padding: 3px;margin-left:2px;margin-right:2px" @click.prevent="click"><slot></slot></span>`
})

Vue.component("code-label-plain", {
	template: `<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px"><slot></slot></span>`
})


// tiny标签
Vue.component("tiny-label", {
	template: `<span class="ui label tiny" style="margin-bottom: 0.5em"><slot></slot></span>`
})

Vue.component("tiny-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em"><slot></slot></span>`
})

// 更小的标签
Vue.component("micro-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em; font-size: 0.7em; padding: 4px"><slot></slot></span>`
})


// 灰色的Label
Vue.component("grey-label", {
	props: ["color"],
	data: function () {
		let color = "grey"
		if (this.color != null && this.color.length > 0) {
			color = "red"
		}
		return {
			labelColor: color
		}
	},
	template: `<span class="ui label basic tiny" :class="labelColor" style="margin-top: 0.4em; font-size: 0.7em; border: 1px solid #ddd!important; font-weight: normal;"><slot></slot></span>`
})

// 可选标签
Vue.component("optional-label", {
	template: `<em><span class="grey">（可选）</span></em>`
})

// Plus专属
Vue.component("plus-label", {
	template: `<span style="color: #B18701;">Plus专属功能。</span>`
})

// 提醒设置项为专业设置
Vue.component("pro-warning-label", {
	template: `<span><i class="icon warning circle yellow"></i>注意：通常不需要修改；如要修改，请在专家指导下进行。</span>`
})
