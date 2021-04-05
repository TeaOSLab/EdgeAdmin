Tea.context(function () {
    this.mediaType = ""
    this.advancedOptionsVisible = true

    let that = this
    this.changeMediaType = function (media) {
        that.mediaType = media.type
    }

    /**
     * 邮箱
     */
    this.emailUsername = "";
    this.emailUsernameHelp = "";

    this.changeEmailUsername = function () {
        this.emailUsernameHelp = "";
        if (this.emailUsername.indexOf("qq.com") > 0) {
            this.emailUsernameHelp = "，<a href=\"https://service.mail.qq.com/cgi-bin/help?id=28\" target='_blank'>QQ邮箱相关设置帮助</a>";
        } else if (this.emailUsername.indexOf("163.com") > 0) {
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

    /**
     * 企业微信
     */
    this.qyWeixinTextFormat = "text";

    /**
     * 企业微信群机器人
     */
    this.qyWeixinRobotTextFormat = "text";

    /**
     * 脚本
     */
    let scriptEditor = null
    this.scriptTab = "path";
    this.scriptLang = "shell";
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
    ];

    this.selectScriptTab = function (tab) {
        this.scriptTab = tab;

        if (tab == "path") {
            this.$delay(function () {
                this.$find("form input[name='scriptPath']").focus();
            });
        } else if (tab == "code") {
            this.$delay(function () {
                this.loadEditor();
            });
        }
    };

    this.selectScriptLang = function (lang) {
        this.scriptLang = lang;
        switch (lang) {
            case "shell":
                scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
                var info = CodeMirror.findModeByMIME("text/x-sh");
                if (info != null) {
                    scriptEditor.setOption("mode", info.mode);
                    CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
                    CodeMirror.autoLoadMode(scriptEditor, info.mode);
                }
                break;
            case "bat":
                scriptEditor.setValue("");
                break;
            case "php":
                scriptEditor.setValue("#!/usr/bin/env php\n\n<?php\n// your PHP codes here");
                var info = CodeMirror.findModeByMIME("text/x-php");
                if (info != null) {
                    scriptEditor.setOption("mode", info.mode);
                    CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
                    CodeMirror.autoLoadMode(scriptEditor, info.mode);
                }
                break;
            case "python":
                scriptEditor.setValue("#!/usr/bin/env python\n\n''' your Python codes here '''");
                var info = CodeMirror.findModeByMIME("text/x-python");
                if (info != null) {
                    scriptEditor.setOption("mode", info.mode);
                    CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
                    CodeMirror.autoLoadMode(scriptEditor, info.mode);
                }
                break;
            case "ruby":
                scriptEditor.setValue("#!/usr/bin/env ruby\n\n# your Ruby codes here");
                var info = CodeMirror.findModeByMIME("text/x-ruby");
                if (info != null) {
                    scriptEditor.setOption("mode", info.mode);
                    CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
                    CodeMirror.autoLoadMode(scriptEditor, info.mode);
                }
                break;
            case "nodejs":
                scriptEditor.setValue("#!/usr/bin/env node\n\n// your javascript codes here");
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
        scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
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
})