Tea.context(function () {
    this.userDescription = ""

    this.changeInstance = function (instance) {
        if (instance != null) {
            this.userDescription = instance.media.userDescription
        } else {
            this.userDescription = ""
        }
    }
})