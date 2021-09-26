Vue.component("http-cache-refs-config-box", {
	props: ["v-cache-refs", "v-cache-config", "v-cache-policy-id"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRefs = []
			ids.forEach(function (id) {
				that.refs.forEach(function (ref) {
					if (ref.id == id) {
						newRefs.push(ref)
					}
				})
			})
			that.updateRefs(newRefs)
			that.change()
		})
	},
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}

		let id = 0
		refs.forEach(function (ref) {
			id++
			ref.id = id
		})
		return {
			refs: refs,
			id: id // 用来对条件进行排序
		}
	},
	methods: {
		addRef: function (isReverse) {
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
			teaweb.popup("/servers/server/settings/cache/createPopup?isReverse=" + (isReverse ? 1 : 0), {
				width: width + "px",
				height: height + "px",
				callback: function (resp) {
					let newRef = resp.data.cacheRef
					if (newRef.conds == null) {
						return
					}

					that.id++
					newRef.id = that.id

					if (newRef.isReverse) {
						let newRefs = []
						let isAdded = false
						that.refs.forEach(function (v) {
							if (!v.isReverse && !isAdded) {
								newRefs.push(newRef)
								isAdded = true
							}
							newRefs.push(v)
						})
						if (!isAdded) {
							newRefs.push(newRef)
						}

						that.updateRefs(newRefs)
					} else {
						that.refs.push(newRef)
					}

					that.change()
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
					resp.data.cacheRef.id = that.refs[index].id
					Vue.set(that.refs, index, resp.data.cacheRef)

					// 通知子组件更新
					that.$refs.cacheRef[index].notifyChange()

					that.change()
				}
			})
		},
		removeRef: function (index) {
			let that = this
			teaweb.confirm("确定要删除此缓存设置吗？", function () {
				that.refs.$remove(index)
				that.change()
			})
		},
		updateRefs: function (newRefs) {
			this.refs = newRefs
			if (this.vCacheConfig != null) {
				this.vCacheConfig.cacheRefs = newRefs
			}
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
		},
		change: function () {
			// 自动保存
			if (this.vCachePolicyId != null && this.vCachePolicyId > 0) {
				Tea.action("/servers/components/cache/updateRefs")
					.params({
						cachePolicyId: this.vCachePolicyId,
						refsJSON: JSON.stringify(this.refs)
					})
					.post()
			}
		}
	},
	template: `<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<div>
		<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
		<table class="ui table selectable celled" v-show="refs.length > 0" id="sortable-table">
			<thead>
				<tr>
					<th style="width:1em"></th>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
					<th class="two op">操作</th>
				</tr>
			</thead>	
			<tbody v-for="(cacheRef, index) in refs" :key="cacheRef.id" :v-id="cacheRef.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td :class="{'color-border': cacheRef.conds.connector == 'and'}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" ref="cacheRef"></http-request-conds-view>
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
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="refs.length > 1">所有条件匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
		
		<button class="ui button tiny" @click.prevent="addRef(false)">+添加缓存设置</button> &nbsp; &nbsp; <a href="" @click.prevent="addRef(true)">+添加不缓存设置</a>
	</div>
	<div class="margin"></div>
</div>`
})