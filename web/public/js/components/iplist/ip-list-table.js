Vue.component("ip-list-table", {
	props: ["v-items"],
	data: function () {
		return {
			items: this.vItems
		}
	},
	methods: {
		updateItem: function (itemId) {
			this.$emit("update-item", itemId)
		},
		deleteItem: function (itemId) {
			this.$emit("delete-item", itemId)
		},
		viewLogs: function (itemId) {
			teaweb.popup("/servers/iplists/accessLogsPopup?itemId=" + itemId, {
				width: "50em",
				height: "30em"
			})
		}
	},
	template: `<div>
 <table class="ui table selectable celled" v-if="items.length > 0">
        <thead>
            <tr>
                <th style="width:18em">IP</th>
                <th>类型</th>
                <th>级别</th>
                <th>过期时间</th>
                <th>备注</th>
                <th class="three op">操作</th>
            </tr>
        </thead>
        <tr v-for="item in items">
            <td>
                <span v-if="item.type != 'all'">{{item.ipFrom}}<span v-if="item.ipTo.length > 0"> - {{item.ipTo}}</span></span>
                <span v-else class="disabled">*</span>
            </td>
            <td>
                <span v-if="item.type.length == 0">IPv4</span>
                <span v-else-if="item.type == 'ipv4'">IPv4</span>
                <span v-else-if="item.type == 'ipv6'">IPv6</span>
                <span v-else-if="item.type == 'all'"><strong>所有IP</strong></span>
            </td>
            <td>
                <span v-if="item.eventLevelName != null && item.eventLevelName.length > 0">{{item.eventLevelName}}</span>
                <span v-else class="disabled">-</span>
            </td>
            <td>
                <div v-if="item.expiredTime.length > 0">
                    {{item.expiredTime}}
                    <div v-if="item.isExpired" style="margin-top: 0.5em">
                        <span class="ui label tiny basic red">已过期</span>
                    </div>
                </div>
                <span v-else class="disabled">不过期</span>
            </td>
            <td>
                <span v-if="item.reason.length > 0">{{item.reason}}</span>
                <span v-else class="disabled">-</span>
            </td>
            <td>
            	<a href="" @click.prevent="viewLogs(item.id)">日志</a> &nbsp;
                <a href="" @click.prevent="updateItem(item.id)">修改</a> &nbsp;
                <a href="" @click.prevent="deleteItem(item.id)">删除</a>
            </td>
        </tr>
    </table>
</div>`
})