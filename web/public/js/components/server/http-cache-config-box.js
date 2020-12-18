Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-is-location", "v-cache-policy"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isPrior: false,
				isOn: false,
				cacheRefs: []
			}
		}
		return {
			cacheConfig: cacheConfig
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.cacheConfig.isPrior) && this.cacheConfig.isOn
		},
		addRef: function () {
			window.UPDATING_CACHE_REF = null

			let width = window.innerWidth
			if (width > 1024) {
				width = 1024
			}
			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup", {
				width: width + "px",
				height: height + "px",
				callback: function (resp) {
					that.cacheConfig.cacheRefs.push(resp.data.cacheRef)
				}
			})
		},
		updateRef: function (index, cacheRef) {
			window.UPDATING_CACHE_REF = cacheRef

			let width = window.innerWidth
			if (width > 1024) {
				width = 1024
			}
			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup", {
				width: width + "px",
				height: height + "px",
				callback: function (resp) {
					Vue.set(that.cacheConfig.cacheRefs, index, resp.data.cacheRef)
				}
			})
		},
		removeRef: function (index) {
			let that = this
			teaweb.confirm("确定要删除此缓存设置吗？", function () {
				that.cacheConfig.cacheRefs.$remove(index)
			})
		},
		timeUnitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
				case "week":
					return "周 "
			}
			return unit
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || cacheConfig.isPrior">
			<tr>
				<td>缓存策略</td>
				<td>
					<div v-if="vCachePolicy != null">{{vCachePolicy.name}}
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
	</table>
	
	<div v-show="isOn()">
		<table class="ui table selectable celled" v-show="cacheConfig.cacheRefs.length > 0">
			<thead>
				<tr>
					<th>条件</th>
					<th class="width10">缓存时间</th>
					<th class="two op">操作</th>
				</tr>
				<tr v-for="(cacheRef, index) in cacheConfig.cacheRefs">
					<td>
						<http-request-conds-view :v-conds="cacheRef.conds"></http-request-conds-view>
					</td>
					<td>{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</td>
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</thead>
		</table>
		
		<button class="ui button tiny" @click.prevent="addRef">+添加缓存设置</button>
	</div>
	<div class="margin"></div>
</div>`
})