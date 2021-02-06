Vue.component("firewall-event-level-options", {
    props: ["v-value"],
    mounted: function () {
        let that = this
        Tea.action("/ui/eventLevelOptions")
            .post()
            .success(function (resp) {
                that.levels = resp.data.eventLevels
                that.change()
            })
    },
    data: function () {
        let value = this.vValue
        if (value == null || value.length == 0) {
            value = "" // 不要给默认值，因为黑白名单等默认值均有不同
        }

        return {
            levels: [],
            description: "",
            level: value
        }
    },
    methods: {
        change: function () {
            this.$emit("change")

            let that = this
            let l = this.levels.$find(function (k, v) {
                return v.code == that.level
            })
            if (l != null) {
                this.description = l.description
            } else {
                this.description = ""
            }
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="eventLevel" v-model="level" @change="change">
        <option v-for="level in levels" :value="level.code">{{level.name}}</option>
    </select>
    <p class="comment">{{description}}</p>
</div>`
})