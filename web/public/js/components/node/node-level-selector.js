// 节点级别选择器
Vue.component("node-level-selector", {
	props: ["v-node-level"],
	data: function () {
		let levelCode = this.vNodeLevel
		if (levelCode == null || levelCode < 1) {
			levelCode = 1
		}
		return {
			levels: [
				{
					name: "边缘节点",
					code: 1,
					description: "普通的边缘节点。"
				},
				{
					name: "L2节点",
					code: 2,
					description: "特殊的边缘节点，同时负责同组上一级节点的回源。"
				}
			],
			levelCode: levelCode
		}
	},
	watch: {
		levelCode: function (code) {
			this.$emit("change", code)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="level" v-model="levelCode">
	<option v-for="level in levels" :value="level.code">{{level.name}}</option>
</select>
<p class="comment" v-if="typeof(levels[levelCode - 1]) != null"><plus-label
></plus-label>{{levels[levelCode - 1].description}}</p>
</div>`
})