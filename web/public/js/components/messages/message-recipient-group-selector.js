Vue.component("message-recipient-group-selector", {
    props: ["v-groups"],
    data: function () {
        let groups = this.vGroups
        if (groups == null) {
            groups = []
        }
        let groupIds = []
        if (groups.length > 0) {
            groupIds = groups.map(function (v) {
                return v.id.toString()
            }).join(",")
        }

        return {
            groups: groups,
            groupIds: groupIds
        }
    },
    methods: {
        addGroup: function () {
            let that = this
            teaweb.popup("/admins/recipients/groups/selectPopup?groupIds=" + this.groupIds, {
                callback: function (resp) {
                    that.groups.push(resp.data.group)
                    that.update()
                }
            })
        },
        removeGroup: function (index) {
            this.groups.$remove(index)
            this.update()
        },
        update: function () {
            let groupIds = []
            if (this.groups.length > 0) {
                this.groups.forEach(function (v) {
                    groupIds.push(v.id)
                })
            }
            this.groupIds = groupIds.join(",")
        }
    },
    template: `<div>
    <input type="hidden" name="groupIds" :value="groupIds"/>
    <div v-if="groups.length > 0">
        <div>
            <div v-for="(group, index) in groups" class="ui label small basic">
                {{group.name}} &nbsp; <a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
            </div>
        </div>
        <div class="ui divider"></div>
    </div>   
    <button class="ui button tiny" type="button" @click.prevent="addGroup()">+</button>
</div>`
})