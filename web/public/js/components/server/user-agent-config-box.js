Vue.component("user-agent-config-box", {
	props: ["v-is-location", "v-is-group", "value"],
	data: function () {
		let config = this.value
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				filters: []
			}
		}
		if (config.filters == null) {
			config.filters = []
		}
		return {
			config: config,
			isAdding: false,
			addingFilter: {
				keywords: [],
				action: "deny"
			},
			moreOptionsVisible: false,
			batchKeywords: ""
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除此名单吗？", function () {
				that.config.filters.$remove(index)
			})
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.batchKeywords.focus()
			})
		},
		confirm: function () {
			if (this.addingFilter.action == "deny") {
				this.config.filters.push(this.addingFilter)
			} else {
				let index = -1
				this.config.filters.forEach(function (filter, filterIndex) {
					if (filter.action == "allow") {
						index = filterIndex
					}
				})

				if (index < 0) {
					this.config.filters.unshift(this.addingFilter)
				} else {
					this.config.filters.$insert(index + 1, this.addingFilter)
				}
			}

			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingFilter = {
				keywords: [],
				action: "deny"
			}
			this.batchKeywords = ""
		},
		changeKeywords: function (keywords) {
			let arr = keywords.split(/\n/)
			let resultKeywords = []
			arr.forEach(function (keyword){
				keyword = keyword.trim()
				if (!resultKeywords.$contains(keyword)) {
					resultKeywords.push(keyword)
				}
			})
			this.addingFilter.keywords = resultKeywords
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
	<input type="hidden" name="userAgentJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用UA名单</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启UserAgent名单。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>UA名单</td>
				<td>
					<div v-if="config.filters.length > 0">
						<table class="ui table celled">
							<thead class="full-width">
								<tr>
									<th>UA关键词</th>
									<th class="two wide">动作</th>
									<th class="one op">操作</th>
								</tr>
							</thead>
							<tbody v-for="(filter, index) in config.filters">
								<tr>
									<td style="background: white">
										<span v-for="keyword in filter.keywords" class="ui label basic tiny">
											<span v-if="keyword.length > 0">{{keyword}}</span>
											<span v-if="keyword.length == 0" class="disabled">[空]</span>
										</span>
									</td>
									<td>
										<span v-if="filter.action == 'allow'" class="green">允许</span><span v-if="filter.action == 'deny'" class="red">不允许</span>
									</td>
									<td><a href="" @click.prevent="remove(index)">删除</a></td>
								</tr>
							</tbody>
						</table>
					</div>
					<div v-if="isAdding" style="margin-top: 0.5em">
						<table class="ui table definition">
							<tr>
								<td class="title">UA关键词</td>
								<td>
									<textarea v-model="batchKeywords" @input="changeKeywords(batchKeywords)" ref="batchKeywords" style="width: 20em" placeholder="*浏览器标识*"></textarea>
									<p class="comment">每行一个关键词；不区分大小写，比如<code-label>Chrome</code-label>；支持<code-label>*</code-label>通配符，比如<code-label>*Firefox*</code-label>；也支持空行关键词，表示空UserAgent。</p>
								</td>
							</tr>
							<tr>
								<td>动作</td>
								<td><select class="ui dropdown auto-width" v-model="addingFilter.action">
										<option value="deny">不允许</option>
										<option value="allow">允许</option>
									</select>
								</td>
							</tr>
						</table>
						<button type="button" class="ui button tiny" @click.prevent="confirm">保存</button> &nbsp; <a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
					</div>
					<div v-show="!isAdding" style="margin-top: 0.5em">
						<button class="ui button tiny" type="button" @click.prevent="add">+</button>
					</div>
				</td>
			</tr>
			<tr>
				<td colspan="2"><more-options-indicator @change="showMoreOptions"></more-options-indicator></td>
			</tr>
		</tbody>
		<tbody v-show="moreOptionsVisible && isOn()">
			<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示只对这些URL进行处理；如果不填则表示支持所有的URL。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})