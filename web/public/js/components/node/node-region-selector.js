Vue.component("node-region-selector", {
	props: ["v-region"],
	data: function () {
		return {
			selectedRegion: this.vRegion
		}
	},
	methods: {
		selectRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		addRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		removeRegion: function () {
			this.selectedRegion = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedRegion != null">
		<input type="hidden" name="regionId" :value="selectedRegion.id"/>
		{{selectedRegion.name}} &nbsp;<a href="" title="删除" @click.prevent="removeRegion()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedRegion == null">
		<a href="" @click.prevent="selectRegion()">[选择区域]</a> &nbsp; <a href="" @click.prevent="addRegion()">[添加区域]</a>
	</div>
</div>`
})