// 消息接收人设置
Vue.component("message-receivers-box", {
	props: ["v-node-cluster-id"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/cluster/settings/message/selectedReceivers")
			.params({
				clusterId: this.clusterId
			})
			.post()
			.success(function (resp) {
				that.receivers = resp.data.receivers
			})
	},
	data: function () {
		let clusterId = this.vNodeClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusterId: clusterId,
			receivers: []
		}
	},
	methods: {
		addReceiver: function () {
			let that = this
			let recipientIdStrings = []
			let groupIdStrings = []
			this.receivers.forEach(function (v) {
				if (v.type == "recipient") {
					recipientIdStrings.push(v.id.toString())
				} else if (v.type == "group") {
					groupIdStrings.push(v.id.toString())
				}
			})

			teaweb.popup("/clusters/cluster/settings/message/selectReceiverPopup?recipientIds=" + recipientIdStrings.join(",") + "&groupIds=" + groupIdStrings.join(","), {
				callback: function (resp) {
					that.receivers.push(resp.data)
				}
			})
		},
		removeReceiver: function (index) {
			this.receivers.$remove(index)
		}
	},
	template: `<div>
        <input type="hidden" name="receiversJSON" :value="JSON.stringify(receivers)"/>           
        <div v-if="receivers.length > 0">
            <div v-for="(receiver, index) in receivers" class="ui label basic small">
               <span v-if="receiver.type == 'group'">分组：</span>{{receiver.name}} <span class="grey small" v-if="receiver.subName != null && receiver.subName.length > 0">({{receiver.subName}})</span> &nbsp; <a href="" title="删除" @click.prevent="removeReceiver(index)"><i class="icon remove"></i></a>
            </div>
             <div class="ui divider"></div>
        </div>
      <button type="button" class="ui button tiny" @click.prevent="addReceiver">+</button>
</div>`
})