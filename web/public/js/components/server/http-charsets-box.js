Vue.component("http-charsets-box", {
	props: ["v-usual-charsets", "v-all-charsets", "v-charset"],
	data: function () {
		let charset = this.vCharset
		if (charset == null) {
			charset = ""
		}
		return {
			charset: charset
		}
	},
	template: `<div>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">选择字符编码</td>
			<td><select class="ui dropdown auto-width" name="charset" v-model="charset">
					<option value="">[未选择]</option>
					<optgroup label="常用字符编码"></optgroup>
					<option v-for="charset in vUsualCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
					<optgroup label="全部字符编码"></optgroup>
					<option v-for="charset in vAllCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
				</select>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`
})