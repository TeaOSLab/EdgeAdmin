Vue.component("http-web-root-box", {
	props: ["v-root-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRootConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				dir: "",
				indexes: [],
				stripPrefix: "",
				decodePath: false,
				isBreak: false,
				exceptHiddenFiles: true,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
		}
		if (config.indexes == null) {
			config.indexes = []
		}

		if (config.onlyURLPatterns == null) {
			config.onlyURLPatterns = []
		}
		if (config.exceptURLPatterns == null) {
			config.exceptURLPatterns = []
		}

		return {
			config: config,
			advancedVisible: false
		}
	},
	methods: {
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		addIndex: function () {
			let that = this
			teaweb.popup("/servers/server/settings/web/createIndex", {
				height: "10em",
				callback: function (resp) {
					that.config.indexes.push(resp.data.index)
				}
			})
		},
		removeIndex: function (i) {
			this.config.indexes.$remove(i)
		},
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="rootJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用静态资源分发</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="title">静态资源根目录</td>
				<td>
					<input type="text" name="root" v-model="config.dir" ref="focus" placeholder="类似于 /home/www"/>
					<p class="comment">可以访问此根目录下的静态资源。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>

		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td>首页文件</td>
				<td>
					<!-- TODO 支持排序 -->
					<div v-if="config.indexes.length > 0">
						<div v-for="(index, i) in config.indexes" class="ui label small basic">
							{{index}} <a href="" title="删除" @click.prevent="removeIndex(i)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="addIndex()">+</button>
					<p class="comment">在URL中只有目录没有文件名时默认查找的首页文件。</p>
				</td>
			</tr>
			<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示不支持通过这些URL访问。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示仅支持通过这些URL访问。</p>
				</td>
			</tr>	
			<tr>
				<td>排除隐藏文件</td>
				<td>
					<checkbox v-model="config.exceptHiddenFiles"></checkbox>
					<p class="comment">排除以点（.）符号开头的隐藏目录或文件，比如<code-label>/.git/logs/HEAD</code-label></p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="config.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>路径解码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.decodePath"/>
						<label></label>	
					</div>
					<p class="comment">是否对请求路径进行URL解码，比如把 <span class="ui label tiny">/Web+App+Browser.html</span> 解码成 <span class="ui label tiny">/Web App Browser.html</span> 再查找文件。</p>
				</td>
			</tr>
			<tr>
				<td>终止请求</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.isBreak"/>
						<label></label>	
					</div>
					<p class="comment">在找不到要访问的文件的情况下是否终止请求并返回404，如果选择终止请求，则不再尝试反向代理等设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})