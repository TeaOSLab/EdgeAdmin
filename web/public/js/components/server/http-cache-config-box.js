Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-cache-policies"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isOn: false,
				cachePolicyId: 0
			}
		}
		return {
			cacheConfig: cacheConfig
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	<table class="ui table definition selectable">
		<tbody>
		<tr>
			<td class="title">是否开启缓存</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="cacheConfig.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
		</tbody>
		<tbody v-show="cacheConfig.isOn">
		<tr>
			<td class="title">选择缓存策略</td>
			<td>
				<span class="disabled" v-if="vCachePolicies.length == 0">暂时没有可选的缓存策略</span>
				<div v-if="vCachePolicies.length > 0">
					<select class="ui dropdown auto-width" v-model="cacheConfig.cachePolicyId">
						<option value="0">[不使用缓存策略]</option>
						<option v-for="policy in vCachePolicies" :value="policy.id">{{policy.name}}</option>
					</select>
				</div>
			</td>
		</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})