Vue.component("message-media-selector", {
    props: ["v-media-type"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/mediaOptions")
            .post()
            .success(function (resp) {
                that.medias = resp.data.medias

                // 初始化简介
                if (that.mediaType.length > 0) {
                    let media = that.medias.$find(function (_, media) {
                        return media.type == that.mediaType
                    })
                    if (media != null) {
                        that.description = media.description
                    }
                }
            })
    },
    data: function () {
        let mediaType = this.vMediaType
        if (mediaType == null) {
            mediaType = ""
        }
        return {
            medias: [],
            description: "",
            mediaType: mediaType
        }
    },
    watch: {
        mediaType: function (v) {
            let media = this.medias.$find(function (_, media) {
                return media.type == v
            })
            if (media == null) {
                this.description = ""
            } else {
                this.description = media.description
            }
            this.$emit("change", media)
        },
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="mediaType" v-model="mediaType">
        <option value="">[选择媒介类型]</option>
        <option v-for="media in medias" :value="media.type">{{media.name}}</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})