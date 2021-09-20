// 启用状态标签
Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="ui label tiny green basic">已启用</span><span v-if="!vIsOn" class="ui label tiny red basic">已停用</span></div>'
})

// 文字代码标签
Vue.component("code-label", {
	methods: {
		click: function (args) {
			this.$emit("click", args)
		}
	},
	template: `<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px" @click.prevent="click"><slot></slot></span>`
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
	template: `<span class="ui label basic grey tiny" style="margin-top: 0.4em; font-size: 0.7em; border: 1px solid #ddd!important; font-weight: normal;"><slot></slot></span>`
})
