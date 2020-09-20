Vue.component("http-gzip-box", {
	props: ["v-gzip-config"],
	data: function () {
		let gzip = this.vGzipConfig
		if (gzip == null) {
			gzip = {
				isOn: true,
				level: 0,
				minLength: null,
				maxLength: null
			}
		}

		return {
			gzip: gzip,
		}
	},
	template: `<div>
<table class="ui table selectable definition">
	<tr>
		<td class="title">压缩级别</td>
		<td>
			<select class="dropdown auto-width" name="level" v-model="gzip.level">
				<option value="0">不压缩</option>
				<option v-for="i in 9" :value="i">{{i}}</option>
			</select>
			<p class="comment">级别越高，压缩比例越大。</p>
		</td>
	</tr>	
	<tr>
		<td>Gzip内容最小长度</td>
		<td>
			<size-capacity-box :v-name="'minLength'" :v-value="gzip.minLength" :v-unit="'kb'"></size-capacity-box>
			<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
		</td>
	</tr>
	<tr>
		<td>Gzip内容最大长度</td>
		<td>
			<size-capacity-box :v-name="'maxLength'" :v-value="gzip.maxLength" :v-unit="'mb'"></size-capacity-box>
			<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
		</td>
	</tr>
</table>
</div>`
})