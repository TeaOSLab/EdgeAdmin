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
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "26em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要移除此页面吗？", function () {
				that.pages.$remove(pageIndex)
			})
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<table class="ui table selectable definition">
	<tr>
		<td class="title">自定义页面</td>
		<td>
			<div v-if="pages.length > 0">
				<div class="ui label small basic" v-for="(page,index) in pages">
					{{page.status}} -&gt; <span v-if="page.bodyType == 'url'">{{page.url}}</span><span v-if="page.bodyType == 'html'">[HTML内容]</span> <a href="" title="修改" @click.prevent="updatePage(index, page.id)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removePage(index)"><i class="icon remove"></i></a>
				</div>
				<div class="ui divider"></div>
			</div>
			<div>
				<button class="ui button small" type="button" @click.prevent="addPage()">+</button>
			</div>
			<p class="comment">根据响应状态码返回一些自定义页面，比如404，500等错误页面。</p>
		</td>
	</tr>
</table>
<div class="ui margin"></div>
</div>`
})