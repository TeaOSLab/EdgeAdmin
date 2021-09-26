// 缓存条件列表
Vue.component("http-cache-refs-box", {
	props: ["v-cache-refs"],
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}
		return {
			refs: refs
		}
	},
	methods: {
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
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
	<div v-show="refs.length > 0">
		<table class="ui table selectable celled">
			<thead>
				<tr>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
				</tr>
				<tr v-for="(cacheRef, index) in refs">
					<td :class="{'color-border': cacheRef.conds.connector == 'and'}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds"></http-request-conds-view>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
					</td>
					<td>
						<span v-if="cacheRef.conds.connector == 'and'">和</span>
						<span v-if="cacheRef.conds.connector == 'or'">或</span>
					</td>
					<td>
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
				</tr>
			</thead>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})