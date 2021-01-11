Tea.context(function () {
    this.isRequesting = true
    this.results = []

    this.$delay(function () {
        this.reload()
    }, 2000)

    this.reload = function () {
        this.isRequesting = true
        this.$post("$")
            .params({
                clusterId: this.clusterId
            })
            .success(function (resp) {
               this.results = resp.data.results
            })
            .done(function () {
                this.isRequesting = false
            })
    }
})