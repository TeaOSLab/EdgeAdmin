Vue.component("url-patterns-box", {
	props: ["value"],
	data: function () {
		let patterns = []
		if (this.value != null) {
			patterns = this.value
		}

		return {
			patterns: patterns,
			isAdding: false,

			addingPattern: {"type": "wildcard", "pattern": ""},
			editingIndex: -1,

			patternIsInvalid: false,

			windowIsSmall: window.innerWidth < 600
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.patternInput.focus()
			})
		},
		edit: function (index) {
			this.isAdding = true
			this.editingIndex = index
			this.addingPattern = {
				type: this.patterns[index].type,
				pattern: this.patterns[index].pattern
			}
		},
		confirm: function () {
			if (this.requireURL(this.addingPattern.type)) {
				let pattern = this.addingPattern.pattern.trim()
				if (pattern.length == 0) {
					let that = this
					teaweb.warn("请输入URL", function () {
						that.$refs.patternInput.focus()
					})
					return
				}
			}
			if (this.editingIndex < 0) {
				this.patterns.push({
					type: this.addingPattern.type,
					pattern: this.addingPattern.pattern
				})
			} else {
				this.patterns[this.editingIndex].type = this.addingPattern.type
				this.patterns[this.editingIndex].pattern = this.addingPattern.pattern
			}
			this.notifyChange()
			this.cancel()
		},
		remove: function (index) {
			this.patterns.$remove(index)
			this.cancel()
			this.notifyChange()
		},
		cancel: function () {
			this.isAdding = false
			this.addingPattern = {"type": "wildcard", "pattern": ""}
			this.editingIndex = -1
		},
		patternTypeName: function (patternType) {
			switch (patternType) {
				case "wildcard":
					return "通配符"
				case "regexp":
					return "正则"
				case "images":
					return "常见图片文件"
				case "audios":
					return "常见音频文件"
				case "videos":
					return "常见视频文件"
			}
			return ""
		},
		notifyChange: function () {
			this.$emit("input", this.patterns)
		},
		changePattern: function () {
			this.patternIsInvalid = false
			let pattern = this.addingPattern.pattern
			switch (this.addingPattern.type) {
				case "wildcard":
					if (pattern.indexOf("?") >= 0) {
						this.patternIsInvalid = true
					}
					break
				case "regexp":
					if (pattern.indexOf("?") >= 0) {
						let pieces = pattern.split("?")
						for (let i = 0; i < pieces.length - 1; i++) {
							if (pieces[i].length == 0 || pieces[i][pieces[i].length - 1] != "\\") {
								this.patternIsInvalid = true
							}
						}
					}
					break
			}
		},
		requireURL: function (patternType) {
			return patternType == "wildcard" || patternType == "regexp"
		}
	},
	template: `<div>
	<div v-show="patterns.length > 0">
		<div v-for="(pattern, index) in patterns" class="ui label basic small" :class="{blue: index == editingIndex, disabled: isAdding && index != editingIndex}" style="margin-bottom: 0.8em">
			<span class="grey" style="font-weight: normal">[{{patternTypeName(pattern.type)}}]</span> <span >{{pattern.pattern}}</span> &nbsp; 
			<a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil tiny"></i></a> 
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-show="isAdding" style="margin-top: 0.5em">
		<div :class="{'ui fields inline': !windowIsSmall}">
			<div class="ui field">
				<select class="ui dropdown auto-width" v-model="addingPattern.type">
					<option value="wildcard">通配符</option>
					<option value="regexp">正则表达式</option>
					<option value="images">常见图片</option>
					<option value="audios">常见音频</option>
					<option value="videos">常见视频</option>
				</select>
			</div>
			<div class="ui field" v-show="addingPattern.type == 'wildcard' || addingPattern.type ==  'regexp'">
				<input type="text" :placeholder="(addingPattern.type == 'wildcard') ? '可以使用星号（*）通配符，不区分大小写' : '可以使用正则表达式，不区分大小写'" v-model="addingPattern.pattern" @input="changePattern" size="36" ref="patternInput" @keyup.enter="confirm()" @keypress.enter.prevent="1" spellcheck="false"/>
				<p class="comment" v-if="patternIsInvalid"><span class="red" style="font-weight: normal"><span v-if="addingPattern.type == 'wildcard'">通配符</span><span v-if="addingPattern.type == 'regexp'">正则表达式</span>中不能包含问号（?）及问号以后的内容。</span></p>
			</div>
			<div class="ui field" style="padding-left: 0"  v-show="addingPattern.type == 'wildcard' || addingPattern.type ==  'regexp'">
				<tip-icon content="通配符示例：<br/>单个路径开头：/hello/world/*<br/>单个路径结尾：*/hello/world<br/>包含某个路径：*/article/*<br/>某个域名下的所有URL：*example.com/*<br/>忽略某个扩展名：*.js" v-if="addingPattern.type == 'wildcard'"></tip-icon>
				<tip-icon content="正则表达式示例：<br/>单个路径开头：^/hello/world<br/>单个路径结尾：/hello/world$<br/>包含某个路径：/article/<br/>匹配某个数字路径：/article/(\\d+)<br/>某个域名下的所有URL：^(http|https)://example.com/" v-if="addingPattern.type == 'regexp'"></tip-icon>
			</div>
			<div class="ui field">
				<button class="ui button tiny" :class="{disabled:this.patternIsInvalid}" type="button" @click.prevent="confirm">确定</button><a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if=!isAdding style="margin-top: 0.5em">
		<button class="ui button tiny basic" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})