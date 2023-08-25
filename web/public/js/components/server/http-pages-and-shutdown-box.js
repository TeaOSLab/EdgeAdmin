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
			bodyType: "html",
			url: "",
			body: "",
			status: 0
		}
		if (this.vShutdownConfig != null) {
			if (this.vShutdownConfig.body == null) {
				this.vShutdownConfig.body = ""
			}
			if (this.vShutdownConfig.bodyType == null) {
				this.vShutdownConfig.bodyType = "html"
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
<html lang="en">
<head>
\t<title>升级中</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
\t<style>
\t\taddress { line-height: 1.8; }
\t</style>
</head>
<body>

<h1>网站升级中</h1>
<p>为了给您提供更好的服务，我们正在升级网站，请稍后重新访问。</p>

<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>

</body>
</html>`
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>
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
	<tr>
		<td>临时关闭页面</td>
		<td>
			<div>
				<table class="ui table selectable definition">
					<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
					<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
						<tr>
							<td class="title">临时关闭网站</td>
							<td>
								<div class="ui checkbox">
									<input type="checkbox" value="1" v-model="shutdownConfig.isOn" />
									<label></label>
								</div>
								<p class="comment">选中后，表示临时关闭当前网站，并显示自定义内容。</p>
							</td>
						</tr>
					</tbody>
					<tbody v-show="(!vIsLocation || shutdownConfig.isPrior) && shutdownConfig.isOn">
						<tr>
							<td>显示内容类型 *</td>
							<td>
								<select class="ui dropdown auto-width" v-model="shutdownConfig.bodyType">
									<option value="html">HTML</option>
									<option value="url">读取URL</option>
								</select>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'url'">
							<td class="title">显示页面URL *</td>
							<td>
								<input type="text" v-model="shutdownConfig.url" placeholder="类似于 https://example.com/page.html"/>
								<p class="comment">将从此URL中读取内容。</p>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'html'">
							<td>显示页面HTML *</td>
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
				<p class="comment">开启临时关闭页面时，所有请求都会直接显示此页面。可用于临时升级网站或者禁止用户访问某个网页。</p>
			</div>
		</td>
	</tr>
</table>
<div class="ui margin"></div>
</div>`
})