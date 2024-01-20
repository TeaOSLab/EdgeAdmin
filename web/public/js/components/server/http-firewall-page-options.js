Vue.component("http-firewall-page-options", {
	props: ["v-page-options"],
	data: function () {
		var defaultPageBody = `<!DOCTYPE html>
<html lang="en">
<head>
	<title>403 Forbidden</title>
	<style>
		address { line-height: 1.8; }
	</style>
</head>
<body>
<h1>403 Forbidden By WAF</h1>
<address>Connection: \${remoteAddr} (Client) -&gt; \${serverAddr} (Server)</address>
<address>Request ID: \${requestId}</address>
</body>
</html>`

		return {
			pageOptions: this.vPageOptions,
			status: this.vPageOptions.status,
			body: this.vPageOptions.body,
			defaultPageBody: defaultPageBody,
			isEditing: false
		}
	},
	watch: {
		status: function (v) {
			if (typeof v === "string" && v.length != 3) {
				return
			}
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.pageOptions.status = 403
			} else {
				this.pageOptions.status = statusCode
			}
		},
		body: function (v) {
			this.pageOptions.body = v
		}
	},
	methods: {
		edit: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<input type="hidden" name="pageOptionsJSON" :value="JSON.stringify(pageOptions)"/>
	<a href="" @click.prevent="edit">状态码：{{status}} / 提示内容：<span v-if="pageOptions.body != null && pageOptions.body.length > 0">[{{pageOptions.body.length}}字符]</span><span v-else class="disabled">[无]</span>
	 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<table class="ui table" v-show="isEditing">
		<tr>
			<td class="title">状态码 *</td>
			<td><input type="text" style="width: 4em" maxlength="3" v-model="status"/></td>
		</tr>
		<tr>
			<td>网页内容</td>
			<td>
				<textarea v-model="body"></textarea>
				<p class="comment"><a href="" @click.prevent="body = defaultPageBody">[使用模板]</a> </p>
			</td>
		</tr>
	</table>
</div>	
`
})