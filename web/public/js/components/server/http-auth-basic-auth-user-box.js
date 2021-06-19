// 基本认证用户配置
Vue.component("http-auth-basic-auth-user-box", {
	props: ["v-users"],
	data: function () {
		let users = this.vUsers
		if (users == null) {
			users = []
		}
		return {
			users: users,
			isAdding: false,
			updatingIndex: -1,

			username: "",
			password: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.username = ""
			this.password = ""

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		cancel: function () {
			this.isAdding = false
			this.updatingIndex = -1
		},
		confirm: function () {
			let that = this
			if (this.username.length == 0) {
				teaweb.warn("请输入用户名", function () {
					that.$refs.username.focus()
				})
				return
			}
			if (this.password.length == 0) {
				teaweb.warn("请输入密码", function () {
					that.$refs.password.focus()
				})
				return
			}
			if (this.updatingIndex < 0) {
				this.users.push({
					username: this.username,
					password: this.password
				})
			} else {
				this.users[this.updatingIndex].username = this.username
				this.users[this.updatingIndex].password = this.password
			}
			this.cancel()
		},
		update: function (index, user) {
			this.updatingIndex = index

			this.isAdding = true
			this.username = user.username
			this.password = user.password

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		remove: function (index) {
			this.users.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="httpAuthBasicAuthUsersJSON" :value="JSON.stringify(users)"/>
	<div v-if="users.length > 0">
		<div class="ui label small basic" v-for="(user, index) in users">
			{{user.username}} <a href="" title="修改" @click.prevent="update(index, user)"><i class="icon pencil tiny"></i></a>
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="用户名" v-model="username" size="15" ref="username"/>
			</div>
			<div class="ui field">
				<input type="password" placeholder="密码" v-model="password" size="15" ref="password"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>&nbsp;
				<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})