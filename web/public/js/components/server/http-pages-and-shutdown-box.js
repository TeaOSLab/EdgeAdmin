Vue.component("http-pages-and-shutdown-box", {
	props: ["v-pages", "v-shutdown-config", "v-is-location"],
	data: function () {
		let pages = []
		if (this.vPages != null) {
			pages = this.vPages
		}
		let shutdownConfig = {
			isPrior: false,
			isOn: false,
			bodyType: "url",
			url: "",
			body: "",
			status: 0
		}
		if (this.vShutdownConfig != null) {
			if (this.vShutdownConfig.body == null) {
				this.vShutdownConfig.body = ""
			}
			if (this.vShutdownConfig.bodyType == null) {
				this.vShutdownConfig.bodyType = "url"
			}
			shutdownConfig = this.vShutdownConfig
		}

		let shutdownStatus = ""
		if (shutdownConfig.status > 0) {
			shutdownStatus = shutdownConfig.status.toString()
		}

		return {
			pages: pages,
			shutdownConfig: shutdownConfig,
			shutdownStatus: shutdownStatus
		}
	},
	watch: {
		shutdownStatus: function (status) {
			let statusInt = parseInt(status)
			if (!isNaN(statusInt) && statusInt > 0 && statusInt < 1000) {
				this.shutdownConfig.status = statusInt
			} else {
				this.shutdownConfig.status = 0
			}
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
		},
		addShutdownHTMLTemplate: function () {
			this.shutdownConfig.body  = `<!DOCTYPE html>
<html>
<head>
\t<title>升级中</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
</head>
<body>

<h3>网站升级中</h3>
<p>为了给您提供更好的服务，我们正在升级网站，请稍后重新访问。</p>

<footer>Powered by GoEdge.</footer>

</body>
</html>`
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>
<table class="ui table selectable definition">
	<tr>
		<td class="title">特殊页面</td>
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
			<p class="comment">根据响应状态码返回一些特殊页面，比如404，500等错误页面。</p>
		</td>
	</tr>	
	<tr>
		<td>临时关闭页面</td>
		<td>
			<div>
				<table class="ui table selectable definition">
					<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
					<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
						<tr>
							<td class="title">是否开启</td>
							<td>
								<div class="ui checkbox">
									<input type="checkbox" value="1" v-model="shutdownConfig.isOn" />
									<label></label>
								</div>
							</td>
						</tr>
					</tbody>
					<tbody v-show="(!vIsLocation || shutdownConfig.isPrior) && shutdownConfig.isOn">
						<tr>
							<td>内容类型 *</td>
							<td>
								<select class="ui dropdown auto-width" v-model="shutdownConfig.bodyType">
									<option value="url">读取URL</option>
									<option value="html">HTML</option>
								</select>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'url'">
							<td class="title">页面URL *</td>
							<td>
								<input type="text" v-model="shutdownConfig.url" placeholder="页面文件路径或一个完整URL"/>
								<p class="comment">页面文件是相对于节点安装目录的页面文件比如pages/40x.html，或者一个完整的URL。</p>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'html'">
							<td>HTML *</td>
							<td>
								<textarea name="body" ref="shutdownHTMLBody" v-model="shutdownConfig.body"></textarea>
								<p class="comment"><a href="" @click.prevent="addShutdownHTMLTemplate">[使用模板]</a>。填写页面的HTML内容，支持请求变量。</p>
							</td>
						</tr>
						<tr>
							<td>状态码</td>
							<td><input type="text" size="3" maxlength="3" name="shutdownStatus" style="width:5.2em" placeholder="状态码" v-model="shutdownStatus"/></td>
						</tr>
					</tbody>
				</table>
				<p class="comment">开启临时关闭页面时，所有请求的响应都会显示此页面。可用于临时升级网站使用。</p>
			</div>
		</td>
	</tr>
</table>
<div class="ui margin"></div>
</div>`
})