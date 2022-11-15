Vue.component("node-cache-disk-dirs-box", {
	props: ["value", "name"],
	data: function () {
		let dirs = this.value
		if (dirs == null) {
			dirs = []
		}
		return {
			dirs: dirs,

			isEditing: false,
			isAdding: false,

			addingPath: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingPath.focus()
			}, 100)
		},
		confirm: function () {
			let addingPath = this.addingPath.trim()
			if (addingPath.length == 0) {
				let that = this
				teaweb.warn("请输入要添加的缓存目录", function () {
					that.$refs.addingPath.focus()
				})
				return
			}
			if (addingPath[0] != "/") {
				addingPath = "/" + addingPath
			}
			this.dirs.push({
				path: addingPath
			})
			this.cancel()
		},
		cancel: function () {
			this.addingPath = ""
			this.isAdding = false
			this.isEditing = false
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除此目录吗？", function () {
				that.dirs.$remove(index)
			})
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(dirs)"/>
	<div style="margin-bottom: 0.3em">
		<span class="ui label small basic" v-for="(dir, index) in dirs">
			<i class="icon folder"></i>{{dir.path}}  &nbsp;  <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" style="width: 30em" v-model="addingPath" @keyup.enter="confirm()" @keypress.enter.prevent="1" @keydown.esc="cancel()" ref="addingPath" placeholder="新的缓存目录，比如 /mnt/cache"/>
			</div>
			<div class="ui field">
				<button class="ui button small" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})