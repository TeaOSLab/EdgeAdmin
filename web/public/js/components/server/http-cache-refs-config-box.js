Vue.component("http-cache-refs-config-box", {
	props: ["v-cache-refs", "v-cache-config", "v-cache-policy-id", "v-web-id", "v-max-bytes"],
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

		let maxBytes = this.vMaxBytes

		let id = 0
		refs.forEach(function (ref) {
			// preset variables
			id++
			ref.id = id
			ref.visible = true

			// check max size
			if (ref.maxSize != null && maxBytes != null && maxBytes.count > 0 && teaweb.compareSizeCapacity(ref.maxSize, maxBytes) > 0) {
				ref.overMaxSize = maxBytes
			}
		})
		return {
			refs: refs,
			id: id // 用来对条件进行排序
		}
	},
	methods: {
		addRef: function (isReverse) {
			window.UPDATING_CACHE_REF = null

			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup?isReverse=" + (isReverse ? 1 : 0), {
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

					// move to bottom
					var afterChangeCallback = function () {
						setTimeout(function () {
							let rightBox = document.querySelector(".right-box")
							if (rightBox != null) {
								rightBox.scrollTo(0, isReverse ? 0 : 100000)
							}
						}, 100)
					}

					that.change(afterChangeCallback)
				}
			})
		},
		updateRef: function (index, cacheRef) {
			window.UPDATING_CACHE_REF = teaweb.clone(cacheRef)

			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup", {
				height: height + "px",
				callback: function (resp) {
					resp.data.cacheRef.id = that.refs[index].id
					Vue.set(that.refs, index, resp.data.cacheRef)
					that.change()
					that.$refs.cacheRef[index].updateConds(resp.data.cacheRef.conds, resp.data.cacheRef.simpleCond)
					that.$refs.cacheRef[index].notifyChange()
				}
			})
		},
		disableRef: function (ref) {
			ref.isOn = false
			this.change()
		},
		enableRef: function (ref) {
			ref.isOn = true
			this.change()
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
		change: function (callback) {
			this.$forceUpdate()

			// 自动保存
			if (this.vCachePolicyId != null && this.vCachePolicyId > 0) { // 缓存策略
				Tea.action("/servers/components/cache/updateRefs")
					.params({
						cachePolicyId: this.vCachePolicyId,
						refsJSON: JSON.stringify(this.refs)
					})
					.post()
			} else if (this.vWebId != null && this.vWebId > 0) { // Server Web or Group Web
				Tea.action("/servers/server/settings/cache/updateRefs")
					.params({
						webId: this.vWebId,
						refsJSON: JSON.stringify(this.refs)
					})
					.success(function (resp) {
						if (resp.data.isUpdated) {
							teaweb.successToast("保存成功", null, function () {
								if (typeof callback == "function") {
									callback()
								}
							})
						}
					})
					.post()
			}
		},
		search: function (keyword) {
			if (typeof keyword != "string") {
				keyword = ""
			}

			this.refs.forEach(function (ref) {
				if (keyword.length == 0) {
					ref.visible = true
					return
				}
				ref.visible = false

				// simple cond
				if (ref.simpleCond != null && typeof ref.simpleCond.value == "string" && teaweb.match(ref.simpleCond.value, keyword)) {
					ref.visible = true
					return
				}

				// composed conds
				if (ref.conds == null || ref.conds.groups == null || ref.conds.groups.length == 0) {
					return
				}

				ref.conds.groups.forEach(function (group) {
					if (group.conds != null) {
						group.conds.forEach(function (cond) {
							if (typeof cond.value == "string" && teaweb.match(cond.value, keyword)) {
								ref.visible = true
							}
						})
					}
				})
			})
			this.$forceUpdate()
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
					<th style="width: 7em">缓存时间</th>
					<th class="three op">操作</th>
				</tr>
			</thead>	
			<tbody v-for="(cacheRef, index) in refs" :key="cacheRef.id" :v-id="cacheRef.id" v-show="cacheRef.visible !== false">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td :class="{'color-border': cacheRef.conds != null && cacheRef.conds.connector == 'and', disabled: !cacheRef.isOn}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" ref="cacheRef" :class="{disabled: !cacheRef.isOn}" v-if="cacheRef.conds != null && cacheRef.conds.groups != null"></http-request-conds-view>
						<http-request-cond-view :v-cond="cacheRef.simpleCond" ref="cacheRef" v-if="cacheRef.simpleCond != null"></http-request-cond-view>
						
						<!-- 特殊参数 -->
						<grey-label v-if="cacheRef.key != null && cacheRef.key.indexOf('\${args}') < 0">忽略URI参数</grey-label>
						
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit.toUpperCase()}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit.toUpperCase()}}</grey-label>
						
						<grey-label v-if="cacheRef.overMaxSize != null"><span class="red">系统限制{{cacheRef.overMaxSize.count}}{{cacheRef.overMaxSize.unit.toUpperCase()}}</span> </grey-label>
						
						<grey-label v-if="cacheRef.methods != null && cacheRef.methods.length > 0">{{cacheRef.methods.join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.expiresTime != null && cacheRef.expiresTime.isPrior && cacheRef.expiresTime.isOn">Expires</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.allowPartialContent">分片缓存</grey-label>
						<grey-label v-if="cacheRef.alwaysForwardRangeRequest">Range回源</grey-label>
						<grey-label v-if="cacheRef.enableIfNoneMatch">If-None-Match</grey-label>
						<grey-label v-if="cacheRef.enableIfModifiedSince">If-Modified-Since</grey-label>
						<grey-label v-if="cacheRef.enableReadingOriginAsync">支持异步</grey-label>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" v-if="cacheRef.isOn" @click.prevent="disableRef(cacheRef)">暂停</a><a href="" v-if="!cacheRef.isOn" @click.prevent="enableRef(cacheRef)"><span class="red">恢复</span></a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="refs.length > 1">所有条件匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。服务设置的优先级比全局缓存策略设置的优先级要高。</p>
		
		<button class="ui button tiny" @click.prevent="addRef(false)" type="button">+添加缓存条件</button> &nbsp; &nbsp; <a href="" @click.prevent="addRef(true)" style="font-size: 0.8em">+添加不缓存条件</a>
	</div>
	<div class="margin"></div>
</div>`
})