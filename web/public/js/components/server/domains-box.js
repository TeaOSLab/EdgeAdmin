// 域名列表
Vue.component("domains-box", {
	props: ["v-domains", "name", "v-support-wildcard"],
	data: function () {
		let domains = this.vDomains
		if (domains == null) {
			domains = []
		}

		let realName = "domainsJSON"
		if (this.name != null && typeof this.name == "string") {
			realName = this.name
		}

		let supportWildcard = true
		if (typeof this.vSupportWildcard == "boolean") {
			supportWildcard = this.vSupportWildcard
		}

		return {
			domains: domains,

			mode: "single", // single | batch
			batchDomains: "",

			isAdding: false,
			addingDomain: "",

			isEditing: false,
			editingIndex: -1,

			realName: realName,
			supportWildcard: supportWildcard
		}
	},
	watch: {
		vSupportWildcard: function (v) {
			if (typeof v == "boolean") {
				this.supportWildcard = v
			}
		},
		mode: function (mode) {
			let that = this
			setTimeout(function () {
				if (mode == "single") {
					if (that.$refs.addingDomain != null) {
						that.$refs.addingDomain.focus()
					}
				} else if (mode == "batch") {
					if (that.$refs.batchDomains != null) {
						that.$refs.batchDomains.focus()
					}
				}
			}, 100)
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
			if (this.mode == "batch") {
				this.confirmBatch()
				return
			}

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
			if (this.supportWildcard) {
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
			} else {
				if (/[*~^]/.test(this.addingDomain)) {
					teaweb.warn("当前只支持添加普通域名，域名中不能含有特殊符号", function () {
						that.$refs.addingDomain.focus()
					})
					return
				}
			}

			if (this.isEditing && this.editingIndex >= 0) {
				this.domains[this.editingIndex] = this.addingDomain
			} else {
				// 分割逗号（，）、顿号（、）
				if (this.addingDomain.match("[，、,;]")) {
					let domainList = this.addingDomain.split(new RegExp("[，、,;]"))
					domainList.forEach(function (v) {
						if (v.length > 0) {
							that.domains.push(v)
						}
					})
				} else {
					this.domains.push(this.addingDomain)
				}
			}
			this.cancel()
			this.change()
		},
		confirmBatch: function () {
			let domains = this.batchDomains.split("\n")
			let realDomains = []
			let that = this
			let hasProblems = false
			domains.forEach(function (domain) {
				if (hasProblems) {
					return
				}
				if (domain.length == 0) {
					return
				}
				if (that.supportWildcard) {
					if (domain == "~") {
						let expr = domain.substring(1)
						try {
							new RegExp(expr)
						} catch (e) {
							hasProblems = true
							teaweb.warn("正则表达式错误：" + e.message, function () {
								that.$refs.batchDomains.focus()
							})
							return
						}
					}
				} else {
					if (/[*~^]/.test(domain)) {
						hasProblems = true
						teaweb.warn("当前只支持添加普通域名，域名中不能含有特殊符号", function () {
							that.$refs.batchDomains.focus()
						})
						return
					}
				}
				realDomains.push(domain)
			})
			if (hasProblems) {
				return
			}
			if (realDomains.length == 0) {
				teaweb.warn("请输入要添加的域名", function () {
					that.$refs.batchDomains.focus()
				})
				return
			}

			realDomains.forEach(function (domain) {
				that.domains.push(domain)
			})
			this.cancel()
			this.change()
		},
		edit: function (index) {
			this.addingDomain = this.domains[index]
			this.isEditing = true
			this.editingIndex = index
			let that = this
			setTimeout(function () {
				that.$refs.addingDomain.focus()
			}, 50)
		},
		remove: function (index) {
			this.domains.$remove(index)
			this.change()
		},
		cancel: function () {
			this.isAdding = false
			this.mode = "single"
			this.batchDomains = ""
			this.isEditing = false
			this.editingIndex = -1
			this.addingDomain = ""
		},
		change: function () {
			this.$emit("change", this.domains)
		}
	},
	template: `<div>
	<input type="hidden" :name="realName" :value="JSON.stringify(domains)"/>
	<div v-if="domains.length > 0">
		<span class="ui label small basic" v-for="(domain, index) in domains" :class="{blue: index == editingIndex}">
			<span v-if="domain.length > 0 && domain[0] == '~'" class="grey" style="font-style: normal">[正则]</span>
			<span v-if="domain.length > 0 && domain[0] == '.'" class="grey" style="font-style: normal">[后缀]</span>
			<span v-if="domain.length > 0 && domain[0] == '*'" class="grey" style="font-style: normal">[泛域名]</span>
			{{domain}}
			<span v-if="!isAdding && !isEditing">
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>
				&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
			</span>
			<span v-if="isAdding || isEditing">
				&nbsp; <a class="disabled"><i class="icon pencil small"></i></a>
				&nbsp; <a class="disabled"><i class="icon remove small"></i></a>
			</span>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding || isEditing">
		<div class="ui fields">
			<div class="ui field" v-if="isAdding">
				<select class="ui dropdown" v-model="mode">
					<option value="single">单个</option>
					<option value="batch">批量</option>
				</select>
			</div>
			<div class="ui field">
				<div v-show="mode == 'single'">
					<input type="text" v-model="addingDomain" @keyup.enter="confirm()" @keypress.enter.prevent="1" @keydown.esc="cancel()" ref="addingDomain" :placeholder="supportWildcard ? 'example.com、*.example.com' : 'example.com、www.example.com'" size="30" maxlength="100"/>
				</div>
				<div v-show="mode == 'batch'">
					<textarea cols="30" v-model="batchDomains" placeholder="example1.com\nexample2.com\n每行一个域名" ref="batchDomains"></textarea>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="supportWildcard">支持普通域名（<code-label>example.com</code-label>）、泛域名（<code-label>*.example.com</code-label>）<span v-if="vSupportWildcard == undefined">、域名后缀（以点号开头，如<code-label>.example.com</code-label>）和正则表达式（以波浪号开头，如<code-label>~.*.example.com</code-label>）</span>；如果域名后有端口，请加上端口号。</p>
		<p class="comment" v-if="!supportWildcard">只支持普通域名（<code-label>example.com</code-label>、<code-label>www.example.com</code-label>）。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})