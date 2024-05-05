Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-is-location", "v-is-group", "v-cache-policy", "v-web-id"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isPrior: false,
				isOn: false,
				addStatusHeader: true,
				addAgeHeader: false,
				enableCacheControlMaxAge: false,
				cacheRefs: [],
				purgeIsOn: false,
				purgeKey: "",
				disablePolicyRefs: false
			}
		}
		if (cacheConfig.cacheRefs == null) {
			cacheConfig.cacheRefs = []
		}

		let maxBytes = null
		if (this.vCachePolicy != null && this.vCachePolicy.maxBytes != null) {
			maxBytes = this.vCachePolicy.maxBytes
		}

		// key
		if (cacheConfig.key == null) {
			// use Vue.set to activate vue events
			Vue.set(cacheConfig, "key", {
				isOn: false,
				scheme: "https",
				host: ""
			})
		}

		return {
			cacheConfig: cacheConfig,
			moreOptionsVisible: false,
			enablePolicyRefs: !cacheConfig.disablePolicyRefs,
			maxBytes: maxBytes,

			searchBoxVisible: false,
			searchKeyword: "",

			keyOptionsVisible: false
		}
	},
	watch: {
		enablePolicyRefs: function (v) {
			this.cacheConfig.disablePolicyRefs = !v
		},
		searchKeyword: function (v) {
			this.$refs.cacheRefsConfigBoxRef.search(v)
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.cacheConfig.isPrior) && this.cacheConfig.isOn
		},
		isPlus: function () {
			return Tea.Vue.teaIsPlus
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
		},
		changeStale: function (stale) {
			this.cacheConfig.stale = stale
		},

		showSearchBox: function () {
			this.searchBoxVisible = !this.searchBoxVisible
			if (this.searchBoxVisible) {
				let that = this
				setTimeout(function () {
					that.$refs.searchBox.focus()
				})
			} else {
				this.searchKeyword = ""
			}
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	
	<table class="ui table definition selectable" v-show="!vIsGroup">
		<tr>
			<td class="title">全局缓存策略</td>
			<td>
				<div v-if="vCachePolicy != null">{{vCachePolicy.name}} <link-icon :href="'/servers/components/cache/policy?cachePolicyId=' + vCachePolicy.id"></link-icon>
					<p class="comment">使用当前网站所在集群的设置。</p>
				</div>
				<span v-else class="red">当前集群没有设置缓存策略，当前配置无法生效。</span>
			</td>
		</tr>
	</table>
	
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || cacheConfig.isPrior">
			<tr>
				<td class="title">启用缓存</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="cacheConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !vIsGroup">
			<tr>
				<td>缓存主域名</td>
				<td>
					<div v-show="!cacheConfig.key.isOn">默认 &nbsp; <a href="" @click.prevent="keyOptionsVisible = !keyOptionsVisible"><span class="small">[修改]</span></a></div>
					<div v-show="cacheConfig.key.isOn">使用主域名：{{cacheConfig.key.scheme}}://{{cacheConfig.key.host}} &nbsp;  <a href="" @click.prevent="keyOptionsVisible = !keyOptionsVisible"><span class="small">[修改]</span></a></div>
					<div v-show="keyOptionsVisible" style="margin-top: 1em">
						<div class="ui divider"></div>
						<table class="ui table definition">
							<tr>
								<td class="title">启用主域名</td>
								<td><checkbox v-model="cacheConfig.key.isOn"></checkbox>
									<p class="comment">启用主域名后，所有缓存键值中的协议和域名部分都会修改为主域名，用来实现缓存不区分域名。</p>
								</td>
							</tr>	
							<tr v-show="cacheConfig.key.isOn">
								<td>主域名 *</td>
								<td>
									<div class="ui fields inline">
										<div class="ui field">
											<select class="ui dropdown" v-model="cacheConfig.key.scheme">
												<option value="https">https://</option>
												<option value="http">http://</option>
											</select>
										</div>
										<div class="ui field">
											<input type="text" v-model="cacheConfig.key.host" placeholder="example.com" @keyup.enter="keyOptionsVisible = false" @keypress.enter.prevent="1"/>
										</div>
									</div>
									<p class="comment">此域名<strong>必须</strong>是当前网站已绑定域名，在刷新缓存时也需要使用此域名。</p>
								</td>
							</tr>
						</table>
						<button class="ui button tiny" type="button" @click.prevent="keyOptionsVisible = false">完成</button>
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
				<td>使用默认缓存条件</td>
				<td>	
					<checkbox v-model="enablePolicyRefs"></checkbox>
					<p class="comment">选中后使用系统全局缓存策略中已经定义的默认缓存条件。</p>
				</td>
			</tr>
			<tr>
				<td>添加X-Cache报头</td>
				<td>
					<checkbox v-model="cacheConfig.addStatusHeader"></checkbox>
					<p class="comment">选中后自动在响应报头中增加<code-label>X-Cache: BYPASS|MISS|HIT|PURGE</code-label>；在浏览器端查看X-Cache值时请先禁用浏览器缓存，避免影响观察。</p>
				</td>
			</tr>
			<tr>
				<td>添加Age Header</td>
				<td>
					<checkbox v-model="cacheConfig.addAgeHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>Age: [存活时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>支持源站控制有效时间</td>
				<td>
					<checkbox v-model="cacheConfig.enableCacheControlMaxAge"></checkbox>
					<p class="comment">选中后表示支持源站在Header中设置的<code-label>Cache-Control: max-age=[有效时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td class="color-border">允许PURGE</td>
				<td>
					<checkbox v-model="cacheConfig.purgeIsOn"></checkbox>
					<p class="comment">允许使用PURGE方法清除某个URL缓存。</p>
				</td>
			</tr>
			<tr v-show="cacheConfig.purgeIsOn">
				<td class="color-border">PURGE Key *</td>
				<td>
					<input type="text" maxlength="200" v-model="cacheConfig.purgeKey"/>
					<p class="comment"><a href="" @click.prevent="generatePurgeKey">[随机生成]</a>。需要在PURGE方法调用时加入<code-label>X-Edge-Purge-Key: {{cacheConfig.purgeKey}}</code-label> Header。只能包含字符、数字、下划线。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-if="isOn() && moreOptionsVisible && isPlus()">
		<h4>过时缓存策略</h4>
		<http-cache-stale-config :v-cache-stale-config="cacheConfig.stale" @change="changeStale"></http-cache-stale-config>
	</div>
	
	<div v-show="isOn()">
		<submit-btn></submit-btn>
		<div class="ui divider"></div>
	</div>
	
	<div v-show="isOn()" style="margin-top: 1em">
		<h4 style="position: relative">缓存条件 &nbsp; <a href="" style="font-size: 0.8em" @click.prevent="$refs.cacheRefsConfigBoxRef.addRef(false)">[添加]</a> &nbsp; <a href="" style="font-size: 0.8em" @click.prevent="showSearchBox" v-show="!searchBoxVisible">[搜索]</a> 
			<div class="ui input small right labeled" style="position: absolute; top: -0.4em; margin-left: 0.5em; zoom: 0.9" v-show="searchBoxVisible">
				<input type="text" placeholder="搜索..." ref="searchBox"  @keypress.enter.prevent="1" @keydown.esc="showSearchBox" v-model="searchKeyword" size="20"/>
				<a href="" class="ui label blue" @click.prevent="showSearchBox"><i class="icon remove small"></i></a>
			</div>
		</h4>
		<http-cache-refs-config-box ref="cacheRefsConfigBoxRef" :v-cache-config="cacheConfig" :v-cache-refs="cacheConfig.cacheRefs" :v-web-id="vWebId" :v-max-bytes="maxBytes"></http-cache-refs-config-box>
	</div>
	<div class="margin"></div>
</div>`
})