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

<div v-if="pages.length > 0" style="max-width: 30em">
	<table class="ui table selectable celled">
		<thead>
			<tr>
				<th class="four wide">响应状态码</th>
				<th>页面类型</th>
				<th style="width: 6.5em">操作</th>
			</tr>	
		</thead>
		<tr v-for="(page,index) in pages">
			<td>
				<span v-if="page.status != null && page.status.length == 1">{{page.status[0]}}</span>
				<span v-else>{{page.status}}</span>
			</td>
			<td style="word-break: break-all"><span v-if="page.bodyType == 'url'">{{page.url}}</span><span v-if="page.bodyType == 'html'">[HTML内容]</span></td>
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