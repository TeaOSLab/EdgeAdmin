Vue.component("http-access-log-box", {
	props: ["v-access-log"],
	data: function () {
		return {
			accessLog: this.vAccessLog
		}
	},
	methods: {
		formatCost: function (seconds) {
			var s = (seconds * 1000).toString();
			var pieces = s.split(".");
			if (pieces.length < 2) {
				return s;
			}

			return pieces[0] + "." + pieces[1].substr(0, 3);
		},
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()
			teaweb.popup("/servers/server/log/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "24em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		}
	},
	template: `<div :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	{{accessLog.remoteAddr}} [{{accessLog.timeLocal}}] <em>&quot;{{accessLog.requestMethod}} {{accessLog.scheme}}://{{accessLog.host}}{{accessLog.requestURI}} <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> {{accessLog.status}} <span v-if="accessLog.attrs != null && accessLog.attrs['cache_cached'] == '1'">[cached]</span> <span v-if="accessLog.attrs != null && accessLog.attrs['waf.action'] != null && accessLog.attrs['waf.action'].length > 0">[waf {{accessLog.attrs['waf.action']}}]</span> - 耗时:{{formatCost(accessLog.requestTime)}} ms
	&nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>
</div>`
})