Vue.component("message-media-instance-selector", {
    props: ["v-instance-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/instances/options")
            .post()
            .success(function (resp) {
                that.instances = resp.data.instances

                // 初始化简介
                if (that.instanceId > 0) {
                    let instance = that.instances.$find(function (_, instance) {
                        return instance.id == that.instanceId
                    })
                    if (instance != null) {
                        that.description = instance.description
                        that.update(instance.id)
                    }
                }
            })
    },
    data: function () {
        let instanceId = this.vInstanceId
        if (instanceId == null) {
            instanceId = 0
        }
        return {
            instances: [],
            description: "",
            instanceId: instanceId
        }
    },
    watch: {
        instanceId: function (v) {
            this.update(v)
        }
    },
    methods: {
        update: function (v) {
            let instance = this.instances.$find(function (_, instance) {
                return instance.id == v
            })
            if (instance == null) {
                this.description = ""
            } else {
                this.description = instance.description
            }
            this.$emit("change", instance)
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="instanceId" v-model="instanceId">
        <option value="0">[选择媒介]</option>
        <option v-for="instance in instances" :value="instance.id">{{instance.name}} ({{instance.media.name}})</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})