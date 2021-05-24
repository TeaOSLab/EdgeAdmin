Vue.component("http-cache-refs-config-box", {
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
					that.refs.push(resp.data.cacheRef)
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
					Vue.set(that.refs, index, resp.data.cacheRef)
				}
			})
		},
		removeRef: function (index) {
			let that = this
			teaweb.confirm("确定要删除此缓存设置吗？", function () {
				that.refs.$remove(index)
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
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<div>
		<table class="ui table selectable celled" v-show="refs.length > 0">
			<thead>
				<tr>
					<th>缓存条件</th>
					<th class="width10">缓存时间</th>
					<th class="two op">操作</th>
				</tr>
				<tr v-for="(cacheRef, index) in refs">
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