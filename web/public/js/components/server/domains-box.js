// 域名列表
Vue.component("domains-box", {
	props: ["v-domains"],
	data: function () {
		let domains = this.vDomains
		if (domains == null) {
			domains = []
		}
		return {
			domains: domains,
			isAdding: false,
			addingDomain: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingDomain.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingDomain = this.addingDomain.replace(/\s/g, "")

			if (this.addingDomain.length == 0) {
				teaweb.warn("请输入要添加的域名", function () {
					that.$refs.addingDomain.focus()
				})
				return
			}


			// 基本校验
			if (this.addingDomain[0] == "~") {
				let expr = this.addingDomain.substring(1)
				try {
					new RegExp(expr)
				} catch (e) {
					teaweb.warn("正则表达式错误：" + e.message, function () {
						that.$refs.addingDomain.focus()
					})
					return
				}
			}

			this.domains.push(this.addingDomain)
			this.cancel()
		},
		remove: function (index) {
			this.domains.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingDomain = ""
		}
	},
	template: `<div>
	<input type="hidden" name="domainsJSON" :value="JSON.stringify(domains)"/>
	<div v-if="domains.length > 0">
		<span class="ui label small basic" v-for="(domain, index) in domains">
			<span v-if="domain.length > 0 && domain[0] == '~'" class="grey" style="font-style: normal">[正则]</span>
			<span v-if="domain.length > 0 && domain[0] == '.'" class="grey" style="font-style: normal">[后缀]</span>
			<span v-if="domain.length > 0 && domain[0] == '*'" class="grey" style="font-style: normal">[泛域名]</span>
			{{domain}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingDomain" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingDomain" placeholder="*.xxx.com" size="30"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">支持普通域名（<code-label>example.com</code-label>）、泛域名（<code-label>*.example.com</code-label>）、域名后缀（以点号开头，如<code-label>.example.com</code-label>）和正则表达式（以波浪号开头，如<code-label>~.*.example.com</code-label>）。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})