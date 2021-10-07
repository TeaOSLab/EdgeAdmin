Vue.component("http-web-root-box", {
	props: ["v-root-config", "v-is-location", "v-is-group"],
	data: function () {
		let rootConfig = this.vRootConfig
		if (rootConfig == null) {
			rootConfig = {
				isPrior: false,
				isOn: true,
				dir: "",
				indexes: [],
				stripPrefix: "",
				decodePath: false,
				isBreak: false
			}
		}
		if (rootConfig.indexes == null) {
			rootConfig.indexes = []
		}
		return {
			rootConfig: rootConfig,
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
					that.rootConfig.indexes.push(resp.data.index)
				}
			})
		},
		removeIndex: function (i) {
			this.rootConfig.indexes.$remove(i)
		},
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.rootConfig.isPrior) && this.rootConfig.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="rootJSON" :value="JSON.stringify(rootConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="rootConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || rootConfig.isPrior">
			<tr>
				<td class="title">是否开启静态资源分发</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="title">静态资源根目录</td>
				<td>
					<input type="text" name="root" v-model="rootConfig.dir" ref="focus" placeholder="类似于 /home/www"/>
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
					<div v-if="rootConfig.indexes.length > 0">
						<div v-for="(index, i) in rootConfig.indexes" class="ui label tiny">
							{{index}} <a href="" title="删除" @click.prevent="removeIndex(i)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="addIndex()">+</button>
					<p class="comment">在URL中只有目录没有文件名时默认查找的首页文件。</p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="rootConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>路径解码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.decodePath"/>
						<label></label>	
					</div>
					<p class="comment">是否对请求路径进行URL解码，比如把 <span class="ui label tiny">/Web+App+Browser.html</span> 解码成 <span class="ui label tiny">/Web App Browser.html</span> 再查找文件。</p>
				</td>
			</tr>
			<tr>
				<td>是否终止请求</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isBreak"/>
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