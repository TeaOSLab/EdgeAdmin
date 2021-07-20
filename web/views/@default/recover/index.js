Tea.context(function () {
	this.STEP_INTRO = "intro"
	this.STEP_NEW_API = "newAPI"
	this.STEP_API_LIST = "apiList"
	this.STEP_FINISH = "finish"

	this.step = this.STEP_INTRO

	// 介绍
	this.goIntroNext = function () {
		this.step = this.STEP_NEW_API
	}

	// 新API地址
	this.apiNodeInfo = {}
	this.apiRequesting = false
	this.goBackIntro = function () {
		this.step = this.STEP_INTRO
	}

	this.apiSubmit = function () {
		this.apiRequesting = true
	}

	this.apiDone = function () {
		this.apiRequesting = false
	}

	this.apiSuccess = function (resp) {
		this.step = this.STEP_API_LIST
		this.apiNodeInfo = resp.data.apiNode
	}

	// 修改API主机地址
	this.goBackAPI = function () {
		this.step = this.STEP_NEW_API
	}

	this.updateHostsSuccess = function () {
		this.step = this.STEP_FINISH
	}
})