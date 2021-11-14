Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-is-location", "v-is-group", "v-cache-policy"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isPrior: false,
				isOn: false,
				addStatusHeader: true,
				cacheRefs: [],
				purgeIsOn: false,
				purgeKey: ""
			}
		}
		if (cacheConfig.cacheRefs == null) {
			cacheConfig.cacheRefs = []
		}
		return {
			cacheConfig: cacheConfig,
			moreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.cacheConfig.isPrior) && this.cacheConfig.isOn
		},
		generatePurgeKey: function () {
			let r = Math.random().toString() + Math.random().toString()
			let s = r.replace(/0\./g, "")
				.replace(/\./g, "")
			let result = ""
			for (let i = 0; i < s.length; i++) {
				result += String.fromCharCode(parseInt(s.substring(i, i + 1)) + ((Math.random() < 0.5) ? "a" : "A").charCodeAt(0))
			}
			this.cacheConfig.purgeKey = result
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || cacheConfig.isPrior">
			<tr v-show="!vIsGroup">
				<td>缓存策略</td>
				<td>
					<div v-if="vCachePolicy != null">{{vCachePolicy.name}} <link-icon :href="'/servers/components/cache/policy?cachePolicyId=' + vCachePolicy.id"></link-icon>
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置缓存策略，当前配置无法生效。</span>
				</td>
			</tr>
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
		<tbody v-show="isOn()">
			<tr>
				<td colspan="2">
					<a href="" @click.prevent="showMoreOptions"><span v-if="moreOptionsVisible">收起选项</span><span v-else>更多选项</span><i class="icon angle" :class="{up: moreOptionsVisible, down:!moreOptionsVisible}"></i></a>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>自动添加X-Cache Header</td>
				<td>
					<checkbox v-model="cacheConfig.addStatusHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>X-Cache: BYPASS|MISS|HIT</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>允许PURGE</td>
				<td>
					<checkbox v-model="cacheConfig.purgeIsOn"></checkbox>
					<p class="comment">允许使用PURGE方法清除某个URL缓存。</p>
				</td>
			</tr>
			<tr v-show="cacheConfig.purgeIsOn">
				<td>PURGE Key *</td>
				<td>
					<input type="text" maxlength="200" v-model="cacheConfig.purgeKey"/>
					<p class="comment"><a href="" @click.prevent="generatePurgeKey">[随机生成]</a>。需要在PURGE方法调用时加入<code-label>Edge-Purge-Key: {{cacheConfig.purgeKey}}</code-label> Header。只能包含字符、数字、下划线。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="isOn()">
		<h4>缓存条件</h4>
		<http-cache-refs-config-box :v-cache-config="cacheConfig" :v-cache-refs="cacheConfig.cacheRefs" ></http-cache-refs-config-box>
	</div>
	<div class="margin"></div>
</div>`
})