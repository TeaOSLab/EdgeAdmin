Vue.component("dns-resolver-config-box", {
	props:["v-dns-resolver-config"],
	data: function () {
		let config = this.vDnsResolverConfig
		if (config == null) {
			config = {
				type: "default"
			}
		}
		return {
			config: config,
			types: [
				{
					name: "默认",
					code: "default"
				},
				{
					name: "CGO",
					code: "cgo"
				},
				{
					name: "Go原生",
					code: "goNative"
				},
			]
		}
	},
	template: `<div>
	<input type="hidden" name="dnsResolverJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">使用的DNS解析库</td>
			<td>
				<select class="ui dropdown auto-width" v-model="config.type">
					<option v-for="t in types" :value="t.code">{{t.name}}</option>
				</select>
				<p class="comment">边缘节点使用的DNS解析库。修改此项配置后，需要重启节点进程才会生效。<pro-warning-label></pro-warning-label></p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})