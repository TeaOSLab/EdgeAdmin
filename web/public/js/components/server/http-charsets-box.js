Vue.component("http-charsets-box", {
	props: ["v-usual-charsets", "v-all-charsets", "v-charset-config", "v-is-location"],
	data: function () {
		let charsetConfig = this.vCharsetConfig
		if (charsetConfig == null) {
			charsetConfig = {
				isPrior: false,
				isOn: false,
				charset: ""
			}
		}
		return {
			charsetConfig: charsetConfig
		}
	},
	template: `<div>
	<input type="hidden" name="charsetJSON" :value="JSON.stringify(charsetConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="charsetConfig" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || charsetConfig.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="(!vIsLocation || charsetConfig.isPrior) && charsetConfig.isOn">	
			<tr>
				<td class="title">选择字符编码</td>
				<td><select class="ui dropdown" style="width:20em" name="charset" v-model="charsetConfig.charset">
						<option value="">[未选择]</option>
						<optgroup label="常用字符编码"></optgroup>
						<option v-for="charset in vUsualCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
						<optgroup label="全部字符编码"></optgroup>
						<option v-for="charset in vAllCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
					</select>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})