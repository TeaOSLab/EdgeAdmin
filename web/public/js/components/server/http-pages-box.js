Vue.component("http-pages-box", {
	props: ["v-pages"],
	data: function () {
		let pages = []
		if (this.vPages != null) {
			pages = this.vPages
		}

		return {
			pages: pages
		}
	},
	methods: {
		addPage: function () {
			let that = this
			teaweb.popup("/servers/server/settings/pages/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.pages.push(resp.data.page)
					that.notifyChange()
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "26em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
					that.notifyChange()
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要移除此页面吗？", function () {
				that.pages.$remove(pageIndex)
				that.notifyChange()
			})
		},
		notifyChange: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					Tea.runActionOn(parent)
				}, 100)
			}
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>

<div v-if="pages.length > 0">
	<table class="ui table selectable celled">
		<thead>
			<tr>
				<th class="two wide">响应状态码</th>
				<th>页面类型</th>
				<th class="two wide">新状态码</th>
				<th>例外URL</th>
				<th>限制URL</th>
				<th class="two op">操作</th>
			</tr>	
		</thead>
		<tr v-for="(page,index) in pages">
			<td>
				<a href="" @click.prevent="updatePage(index, page.id)">
					<span v-if="page.status != null && page.status.length == 1">{{page.status[0]}}</span>
					<span v-else>{{page.status}}</span>
					
					<i class="icon expand small"></i>
				</a>
			</td>
			<td style="word-break: break-all">
				<div v-if="page.bodyType == 'url'">
					{{page.url}}
					<div>
						<grey-label>读取URL</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'redirectURL'">
					{{page.url}}
					<div>
						<grey-label>跳转URL</grey-label>	
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
				<div v-if="page.bodyType == 'html'">
					[HTML内容]
					<div>
						<grey-label v-if="page.newStatus > 0">{{page.newStatus}}</grey-label>
					</div>
				</div>
			</td>
			<td>
				<span v-if="page.newStatus > 0">{{page.newStatus}}</span>
				<span v-else class="disabled">保持</span>	
			</td>
			<td>
				<div v-if="page.exceptURLPatterns != null && page.exceptURLPatterns">
					<span v-for="urlPattern in page.exceptURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<div v-if="page.onlyURLPatterns != null && page.onlyURLPatterns">
					<span v-for="urlPattern in page.onlyURLPatterns" class="ui basic label small">{{urlPattern.pattern}}</span>
				</div>
				<span v-else class="disabled">-</span>
			</td>
			<td>
				<a href="" title="修改" @click.prevent="updatePage(index, page.id)">修改</a> &nbsp; 
				<a href="" title="删除" @click.prevent="removePage(index)">删除</a>
			</td>
		</tr>
	</table>
</div>
<div style="margin-top: 1em">
	<button class="ui button small" type="button" @click.prevent="addPage()">+添加自定义页面</button>
</div>
<div class="ui margin"></div>
</div>`
})