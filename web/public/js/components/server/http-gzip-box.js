Vue.component("http-gzip-box", {
	props: ["v-gzip-config", "v-gzip-ref", "v-is-location"],
	data: function () {
		let gzip = this.vGzipConfig
		if (gzip == null) {
			gzip = {
				isOn: true,
				level: 0,
				minLength: null,
				maxLength: null,
				conds: null
			}
		}

		return {
			gzip: gzip,
			advancedVisible: false
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.vGzipRef.isPrior) && this.vGzipRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
<input type="hidden" name="gzipRefJSON" :value="JSON.stringify(vGzipRef)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="vGzipRef" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || vGzipRef.isPrior">
		<tr>
			<td class="title">启用Gzip压缩</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="vGzipRef.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
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
	</tbody>
	<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
	<tbody v-show="isOn() && advancedVisible">
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
		<tr>
			<td>匹配条件</td>
			<td>
				<http-request-conds-box :v-conds="gzip.conds"></http-request-conds-box>
</td>
		</tr>
	</tbody>
</table>
</div>`
})