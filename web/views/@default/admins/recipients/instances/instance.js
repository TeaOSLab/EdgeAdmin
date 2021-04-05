Tea.context(function () {
    let scriptEditor = null

    this.from = encodeURIComponent(window.location.toString())

    if (this.instance.media.type == "script" && this.instance.params.scriptType == "code") {
        this.$delay(function () {
            this.loadEditor()
        })
    }

    this.loadEditor = function () {
        if (scriptEditor == null) {
            scriptEditor = CodeMirror(document.getElementById("script-code-editor"), {
                theme: "idea",
                lineNumbers: false,
                value: "",
                readOnly: true,
                showCursorWhenSelecting: true,
                height: "auto",
                //scrollbarStyle: null,
                viewportMargin: Infinity,
                lineWrapping: true,
                highlightFormatting: false,
                indentUnit: 4,
                indentWithTabs: true
            })
        }
        scriptEditor.setValue(this.instance.params.script)

        let lang = "shell"
        if (this.instance.params.scriptLang != null && this.instance.params.scriptLang.length > 0) {
            lang = this.instance.params.scriptLang
        }
        let mimeType = "text/x-" + lang
        if (lang == "nodejs") {
            mimeType = "text/javascript"
        } else if (lang == "shell") {
            mimeType = "text/x-sh"
        }
        let info = CodeMirror.findModeByMIME(mimeType)
        if (info != null) {
            scriptEditor.setOption("mode", info.mode)
            CodeMirror.modeURL = "/codemirror/mode/%N/%N.js"
            CodeMirror.autoLoadMode(scriptEditor, info.mode)
        }
    }
})