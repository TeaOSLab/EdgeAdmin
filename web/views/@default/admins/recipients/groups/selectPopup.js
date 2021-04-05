Tea.context(function () {
    this.selectGroup = function (group) {
        NotifyPopup({
            code: 200,
            data: {
                group: group
            }
        })
    }
})