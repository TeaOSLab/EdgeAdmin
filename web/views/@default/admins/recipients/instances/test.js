Tea.context(function () {
    this.isRunning = false
    this.isFinished = false
    this.response = ""
    this.error = ""

    this.submitBefore = function () {
        this.isRunning = true
        this.isFinished = false
        this.response = ""
        this.error = ""
    }

    this.submitSuccess = function (resp) {
        this.reloadStatus(resp.data.taskId)
    }

    this.submitFail = function (resp) {
        this.isRunning = false
        this.isFinished = true
        this.response = ""
        this.error = resp.errors[0].messages[0]
        this.errorLines = []
    }

    this.submitError = function () {
        this.isRunning = false
        this.isFinished = true
        this.response = ""
        this.errorLines = []
        this.error = "请求超时"
    }

    // 更新状态
    this.reloadStatus = function (taskId) {
        let isDone = false
        this.$post("/admins/recipients/tasks/taskInfo")
            .params({
                taskId: taskId
            })
            .success(function (resp) {
                if (resp.data.status == 2 || resp.data.status == 3) {
                    isDone = true
                    this.updateStatus(resp.data.result)
                }
            })
            .done(function () {
                this.$delay(function () {
                    if (isDone) {
                        return
                    }
                    this.reloadStatus(taskId)
                }, 3000)
            })
    }

    this.updateStatus = function (result) {
        this.isRunning = false
        this.isFinished = true
        this.response = result.response
        this.responseLines = []
        if (this.response != null) {
            this.responseLines = this.response.split("\n")
        }
        this.error = result.error
        this.errorLines = []
        if (this.error.length > 0) {
            this.errorLines = this.error.split("\n")
        }
    }
})