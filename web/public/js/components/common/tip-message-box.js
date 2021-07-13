// 信息提示窗口
Vue.component("tip-message-box", {
	props: ["code"],
	mounted: function () {
		let that = this
		Tea.action("/ui/showTip")
			.params({
				code: this.code
			})
			.success(function (resp) {
				that.visible = resp.data.visible
			})
			.post()
	},
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		close: function () {
			this.visible = false
			Tea.action("/ui/hideTip")
				.params({
					code: this.code
				})
				.post()
		}
	},
	template: `<div class="ui icon message" v-if="visible">
	<i class="icon info circle"></i>
	<i class="close icon" title="取消" @click.prevent="close" style="margin-top: 1em"></i>
	<div class="content">
		<slot></slot>
	</div>
</div>`
})