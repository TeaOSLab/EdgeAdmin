// TODO 支持关键词搜索
// TODO 改成弹窗选择
Vue.component("admin-selector", {
    props: ["v-admin-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/options")
            .post()
            .success(function (resp) {
                that.admins = resp.data.admins
            })
    },
    data: function () {
        let adminId = this.vAdminId
        if (adminId == null) {
            adminId = 0
        }
        return {
            admins: [],
            adminId: adminId
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="adminId" v-model="adminId">
        <option value="0">[选择系统用户]</option>
        <option v-for="admin in admins" :value="admin.id">{{admin.name}}（{{admin.username}}）</option>
    </select>
</div>`
})