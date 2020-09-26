Vue.component("more-options-tbody", {
	data: function () {
		return {
			isVisible: false
		}
	},
	methods: {
		show: function () {
			this.isVisible = !this.isVisible
			this.$emit("change", this.isVisible)
		}
	},
	template: `<tbody>
	<tr>
		<td colspan="2"><a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a></td>
	</tr>
</tbody>`
})