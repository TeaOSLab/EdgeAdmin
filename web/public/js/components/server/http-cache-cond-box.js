Vue.component("http-cache-cond-box", {
	template: `<div>
<table class="ui table definition selectable">
	<tr>
		<td>匹配条件</td>
		<td><http-request-conds-box></http-request-conds-box></td>
	</tr>
	<tr>
		<td>缓存有效期</td>
		<td>
			<time-duration-box :name="'lifeJSON'" :v-count="3600" :v-unit="'second'"></time-duration-box>
		</td>
	</tr>
	<tr>
		<td>状态码列表</td>
		<td>
			<values-box name="statusList" size="3" maxlength="3" :values="['200']"></values-box>
			<p class="comment">允许缓存的HTTP状态码列表。</p>
		</td>
	</tr>
	<tr>
		<td>跳过的Cache-Control值</td>
		<td>
			<values-box name="skipResponseCacheControlValues" size="10" maxlength="100" :values="['private', 'no-cache', 'no-store']"></values-box>
			<p class="comment">当响应的Cache-Control为这些值时不缓存响应内容，而且不区分大小写。</p>
		</td>
	</tr>
	<tr>
		<td>跳过Set-Cookie</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="skipResponseSetCookie" value="1" checked="checked"/>
				<label></label>
			</div>
			<p class="comment">选中后，当响应的Header中有Set-Cookie时不缓存响应内容。</p>
		</td>
	</tr>
	<tr>
		<td>支持请求no-cache刷新</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="enableRequestCachePragma" value="1"/>
				<label></label>
			</div>
			<p class="comment">选中后，当请求的Header中含有Pragma: no-cache或Cache-Control: no-cache时，会跳过缓存直接读取源内容。</p>
		</td>
	</tr>	
</table>
</div>`
})