Vue.component("http-pages-and-shutdown-box", {
	props: ["v-enable-global-pages", "v-pages", "v-shutdown-config", "v-is-location"],
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
			shutdownStatus: shutdownStatus,
			enableGlobalPages: this.vEnableGlobalPages
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
				height: "30em",
				callback: function (resp) {
					that.pages.push(resp.data.page)
					that.notifyChange()
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
					that.notifyChange()
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要删除此自定义页面吗？", function () {
				that.pages.$remove(pageIndex)
				that.notifyChange()
			})
		},
		addShutdownHTMLTemplate: function () {
			this.shutdownConfig.body = `<!DOCTYPE html>
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
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>

<h4 style="margin-bottom: 0.5em">自定义页面</h4>

<p class="comment" style="padding-top: 0; margin-top: 0">根据响应状态码返回一些自定义页面，比如404，500等错误页面。</p>

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

<h4 style="margin-top: 2em;">临时关闭页面</h4>
<p class="comment" style="margin-top: 0; padding-top: 0">开启临时关闭页面时，所有请求都会直接显示此页面。可用于临时升级网站或者禁止用户访问某个网页。</p>	
<div>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
			<tr>
				<td class="title">启用临时关闭网站</td>
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
						<option value="redirectURL">跳转URL</option>
					</select>
				</td>
			</tr>
			<tr v-if="shutdownConfig.bodyType == 'url'">
				<td class="title">显示页面URL *</td>
				<td>
					<input type="text" v-model="shutdownConfig.url" placeholder="类似于 https://example.com/page.html"/>
					<p class="comment">将从此URL中读取内容。</p>
				</td>
			</tr>
			<tr v-if="shutdownConfig.bodyType == 'redirectURL'">
				<td class="title">跳转到URL *</td>
				<td>
					<input type="text" v-model="shutdownConfig.url" placeholder="类似于 https://example.com/page.html"/>
					 <p class="comment">将会跳转到此URL。</p>
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
</div>

<h4 style="margin-top: 2em;">其他设置</h4>
<table class="ui table definition selectable">
	<tr>
		<td class="title">启用系统自定义页面</td>
		<td>
			<checkbox name="enableGlobalPages" v-model="enableGlobalPages"></checkbox>
			<p class="comment">选中后，表示如果当前网站没有自定义页面，则尝试使用系统对应的自定义页面。</p>
		</td>
	</tr>
</table>
<div class="ui margin"></div>

</div>`
})