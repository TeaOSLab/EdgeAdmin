/**
 * 更多选项
 */
Vue.component("more-options-indicator", {
	props:[],
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		changeVisible: function () {
			this.visible = !this.visible
			if (Tea.Vue != null) {
				Tea.Vue.moreOptionsVisible = this.visible
			}
			this.$emit("change", this.visible)
			this.$emit("input", this.visible)
		}
	},
	template: '<a href="" style="font-weight: normal" @click.prevent="changeVisible()"><slot><span v-if="!visible">更多选项</span><span v-if="visible">收起选项</span></slot> <i class="icon angle" :class="{down:!visible, up:visible}"></i> </a>'
});