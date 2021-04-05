Tea.context(function () {
    this.userDescription = ""

    this.changeInstance = function (instance) {
        if (instance != null) {
            this.userDescription = instance.media.userDescription
        } else {
            this.userDescription = ""
        }
    }

    this.success = function () {
        let that = this
        teaweb.success("保存成功", function () {
            window.location = Tea.url(".recipient", {
                recipientId: that.recipient.id
            })
        })
    }
})