Tea.context(function () {
	let scriptEditor = null
	let isLoaded = false;

	this.$delay(function () {
		isLoaded = true;

		if (this.instance.media.type == "email") {
			this.changeEmailUsername()
		}
	})

	this.success = function () {
	    let that = this
		teaweb.success("保存成功", function () {
		    window.location = "/admins/recipients/instances/instance?instanceId=" + that.instance.id
        })
	};

	/**
	 * 名称
	 */
	this.rateNoticeVisible = false

	this.changeName = function (name) {
		if (name.indexOf("短信") > -1 || name.indexOf("钉钉") > -1 || name.indexOf("微信") > -1) {
			this.rateNoticeVisible = true
		} else {
			this.rateNoticeVisible = false
		}
	};
	this.changeName(this.instance.media.name)

	/**
	 * 类型
	 */
	this.mediaType = this.instance.media.type;

	this.changeMediaType = function (media) {
	    this.mediaType = media.type
		if (!isLoaded) {
			return;
		}
		if (this.mediaType == "email") {
			this.$delay(function () {
				this.$find("form input[name='emailSmtp']").focus();
			});
		} else if (this.mediaType == "webHook") {
			this.$delay(function () {
				this.$find("form input[name='webHookURL']").focus();
			});
		} else if (this.mediaType == "script") {
			this.$delay(function () {
				this.$find("form input[name='scriptPath']").focus();
			});
		} else if (this.mediaType == "dingTalk") {
			this.$delay(function () {
				this.$find("form textarea[name='dingTalkWebHookURL']").focus();
			});
		} else if (this.mediaType == "qyWeixin") {
			this.$delay(function () {
				this.$find("form input[name='qyWeixinCorporateId']").focus();
			});
		} else if (this.mediaType == "qyWeixinRobot") {
			this.$delay(function () {
				this.$find("form textarea[name='qyWeixinRobotWebHookURL']").focus();
			});
		}
	};
	this.changeMediaType(this.instance.media);

	/**
	 * 邮箱
	 */
	this.emailUsernameHelp = "";

	this.changeEmailUsername = function () {
		this.emailUsernameHelp = "";
		if (this.instance.params.username.indexOf("qq.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://service.mail.qq.com/cgi-bin/help?id=28\" target='_blank'>QQ邮箱相关设置帮助</a>";
		} else if (this.instance.params.username.indexOf("163.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://help.mail.163.com/faqDetail.do?code=d7a5dc8471cd0c0e8b4b8f4f8e49998b374173cfe9171305fa1ce630d7f67ac22dc0e9af8168582a\" target='_blank'>网易邮箱相关设置帮助</a>";
		}
	};

	/**
	 * webHook
	 */
    this.methods = ["GET", "POST"]
	this.webHookMethod = "GET";
	this.webHookHeadersAdding = false;
	this.webHookHeaders = [];
	this.webHookHeadersAddingName = "";
	this.webHookHeadersAddingValue = "";

	this.addWebHookHeader = function () {
		this.webHookHeadersAdding = true;
		this.$delay(function () {
			this.$find("form input[name='webHookHeaderName']").focus();
		});
	};

	this.cancelWebHookHeadersAdding = function () {
		this.webHookHeadersAdding = false;
	};

	this.confirmWebHookHeadersAdding = function () {
		this.webHookHeaders.push({
			"name": this.webHookHeadersAddingName,
			"value": this.webHookHeadersAddingValue
		});
		this.webHookHeadersAddingName = "";
		this.webHookHeadersAddingValue = "";
		this.webHookHeadersAdding = false;
	};

	this.removeWebHookHeader = function (index) {
		if (!window.confirm("确定要删除此Header吗？")) {
			return;
		}
		this.webHookHeaders.$remove(index);
	};

	this.webHookContentType = "params";

	this.selectWebHookContentType = function (contentType) {
		this.webHookContentType = contentType;
		this.$delay(function () {
			if (contentType == "params") {

			} else if (contentType == "body") {
				this.$find("form textarea[name='webHookBody']").focus();
			}
		});
	};

	this.webHookParamsAdding = false;
	this.webHookParams = [];
	this.webHookParamsAddingName = "";
	this.webHookParamsAddingValue = "";

	this.addWebHookParam = function () {
		this.webHookParamsAdding = true;
		this.$delay(function () {
			this.$find("form input[name='webHookParamName']").focus();
		});
	};

	this.cancelWebHookParamsAdding = function () {
		this.webHookParamsAdding = false;
	};

	this.confirmWebHookParamsAdding = function () {
		this.webHookParams.push({
			"name": this.webHookParamsAddingName,
			"value": this.webHookParamsAddingValue
		});
		this.webHookParamsAddingName = "";
		this.webHookParamsAddingValue = "";
		this.webHookParamsAdding = false;
	};

	this.removeWebHookParam = function (index) {
		if (!window.confirm("确定要删除此参数吗？")) {
			return;
		}
		this.webHookParams.$remove(index);
	};

	this.webHookBody = "";

	if (this.instance.media.type == "webHook") {
		this.webHookMethod = this.instance.params.method;
		if (this.instance.params.headers != null) {
			this.webHookHeaders = this.instance.params.headers;
		}

		if (this.instance.params.contentType == "params") {
			this.webHookContentType = "params";
			if (this.instance.params.params != null) {
				this.webHookParams = this.instance.params.params;
			}
		}

		if (this.instance.params.contentType == "body") {
			this.webHookContentType = "body";
			this.webHookBody = this.instance.params.body;
		}
	}

	/**
	 * 脚本
	 */
	this.scriptTab = "path";
	this.scriptLang = "shell";

	if (this.instance.media.type == "script") {
		if (this.instance.params.scriptType == "path") {
			this.scriptTab = "path";
		} else {
			this.scriptTab = "code";
			this.scriptLang = this.instance.params.scriptLang;
			this.$delay(function () {
				this.loadEditor();
			});
		}
	}

	this.scriptLangs = [
		{
			"name": "Shell",
			"code": "shell"
		},
		{
			"name": "批处理(bat)",
			"code": "bat"
		},
		{
			"name": "PHP",
			"code": "php"
		},
		{
			"name": "Python",
			"code": "python"
		},
		{
			"name": "Ruby",
			"code": "ruby"
		},
		{
			"name": "NodeJS",
			"code": "nodejs"
		}
	]

	this.selectScriptTab = function (tab) {
		this.scriptTab = tab

		if (tab == "path") {
			this.$delay(function () {
				this.$find("form input[name='scriptPath']").focus()
			})
		} else if (tab == "code") {
			this.$delay(function () {
				this.loadEditor()
			})
		}
	};

	this.selectScriptLang = function (lang) {
		this.scriptLang = lang;
		switch (lang) {
			case "shell":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "shell") {
					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
				}
				var info = CodeMirror.findModeByMIME("text/x-sh");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "bat":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "bat") {
					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("");
				}
				break;
			case "php":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "php") {
					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env php\n\n<?php\n// your PHP codes here");
				}
				var info = CodeMirror.findModeByMIME("text/x-php");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "python":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "python") {

					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env python\n\n''' your Python codes here '''");
				}
				var info = CodeMirror.findModeByMIME("text/x-python");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "ruby":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "ruby") {
					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env ruby\n\n# your Ruby codes here");
				}
				var info = CodeMirror.findModeByMIME("text/x-ruby");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "nodejs":
				if (this.instance.media.type == "script" && this.instance.params.scriptType == "code" && this.instance.params.scriptLang == "nodejs") {
					scriptEditor.setValue(this.instance.params.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env node\n\n// your javascript codes here");
				}
				var info = CodeMirror.findModeByMIME("text/javascript");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
		}

		scriptEditor.save();
		scriptEditor.focus();
	};

	this.loadEditor = function () {
		if (scriptEditor == null) {
			scriptEditor = CodeMirror.fromTextArea(document.getElementById("script-code-editor"), {
				theme: "idea",
				lineNumbers: true,
				value: "",
				readOnly: false,
				showCursorWhenSelecting: true,
				height: "auto",
				//scrollbarStyle: null,
				viewportMargin: Infinity,
				lineWrapping: true,
				highlightFormatting: false,
				indentUnit: 4,
				indentWithTabs: true
			});
		}
		if (this.instance.params.script != null && this.instance.params.script.length > 0) {
			scriptEditor.setValue(this.instance.params.script);
		} else {
			scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
		}
		scriptEditor.save();
		scriptEditor.focus();

		var info = CodeMirror.findModeByMIME("text/x-sh");
		if (info != null) {
			scriptEditor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(scriptEditor, info.mode);
		}

		scriptEditor.on("change", function () {
			scriptEditor.save();
		});
	};

	/**
	 * 环境变量
	 */
	this.env = [];
	if (this.instance.media.type == "script" && this.instance.params.env != null) {
		this.env = this.instance.params.env;
	}

	this.envAdding = false;
	this.envAddingName = "";
	this.envAddingValue = "";

	this.addEnv = function () {
		this.envAdding = !this.envAdding;
		this.$delay(function () {
			this.$find("form input[name='envAddingName']").focus();
		});
	};

	this.confirmAddEnv = function () {
		if (this.envAddingName.length == 0) {
			alert("请输入变量名");
			this.$find("form input[name='envAddingName']").focus();
			return;
		}
		this.env.push({
			"name": this.envAddingName,
			"value": this.envAddingValue
		});
		this.envAdding = false;
		this.envAddingName = "";
		this.envAddingValue = "";
	};

	this.removeEnv = function (index) {
		this.env.$remove(index);
	};

	this.cancelEnv = function () {
		this.envAdding = false;
	};

	/**
	 * 阿里云短信模板
	 */
	this.aliyunSmsTemplateVars = [];
	if (this.instance.params.variables != null) {
		this.aliyunSmsTemplateVars = this.instance.params.variables;
	}
	this.aliyunSmsTemplateVarAdding = false;
	this.aliyunSmsTemplateVarAddingName = "";
	this.aliyunSmsTemplateVarAddingValue = "";

	this.addAliyunSmsTemplateVar = function () {
		this.aliyunSmsTemplateVarAdding = !this.aliyunSmsTemplateVarAdding;
		this.$delay(function () {
			this.$find("form input[name='aliyunSmsTemplateVarAddingName']").focus();
		});
	};

	this.confirmAddAliyunSmsTemplateVar = function () {
		if (this.aliyunSmsTemplateVarAddingName.length == 0) {
			alert("请输入变量名");
			this.$find("form input[name='aliyunSmsTemplateVarAddingName']").focus();
			return;
		}
		this.aliyunSmsTemplateVars.push({
			"name": this.aliyunSmsTemplateVarAddingName,
			"value": this.aliyunSmsTemplateVarAddingValue
		});
		this.aliyunSmsTemplateVarAdding = false;
		this.aliyunSmsTemplateVarAddingName = "";
		this.aliyunSmsTemplateVarAddingValue = "";
	};

	this.removeAliyunSmsTemplateVar = function (index) {
		this.aliyunSmsTemplateVars.$remove(index);
	};

	this.cancelAliyunSmsTemplateVar = function () {
		this.aliyunSmsTemplateVarAdding = false;
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = true;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});