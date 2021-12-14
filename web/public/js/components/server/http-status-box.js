// 请求方法列表
Vue.component("http-status-box", {
	props: ["v-status-list"],
	data: function () {
		let statusList = this.vStatusList
		if (statusList == null) {
			statusList = []
		}
		return {
			statusList: statusList,
			isAdding: false,
			addingStatus: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingStatus.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingStatus = this.addingStatus.replace(/\s/g, "").toUpperCase()

			if (this.addingStatus.length == 0) {
				teaweb.warn("请输入要添加的状态码", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			// 是否已经存在
			if (this.statusList.$contains(this.addingStatus)) {
				teaweb.warn("此状态码已经存在，无需重复添加", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			// 格式
			if (!this.addingStatus.match(/^\d{3}$/)) {
				teaweb.warn("请输入正确的状态码", function () {
					that.$refs.addingStatus.focus()
				})
				return
			}

			this.statusList.push(parseInt(this.addingStatus, 10))
			this.cancel()
		},
		remove: function (index) {
			this.statusList.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingStatus = ""
		}
	},
	template: `<div>
	<input type="hidden" name="statusListJSON" :value="JSON.stringify(statusList)"/>
	<div v-if="statusList.length > 0">
		<span class="ui label small basic" v-for="(status, index) in statusList">
			{{status}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingStatus" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingStatus" placeholder="如200" size="3" maxlength="3" style="width: 5em"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为三位数字，比如<code-label>200</code-label>、<code-label>404</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})