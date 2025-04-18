
function process(ss, curIndex, env, host) {
    let obj = JSON.parse(ss);

    let versionDiv = document.getElementById("version")
    let version = document.createElement('pre')
    version.textContent = `当前最新配置版本编号：${curIndex}      `
    versionDiv.appendChild(version)
    let refreshButton = document.createElement("button")
    versionDiv.appendChild(refreshButton)
    refreshButton.textContent = "刷新"
    refreshButton.addEventListener("click", () => {
        window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}`
    })
    refreshButton.style.height = "50%"
    refreshButton.style.alignSelf = "center"


    let versionSelectorDiv = document.getElementById("versionseletor")
    let versionSelectorLabel = document.createElement("label")
    versionSelectorLabel.textContent = "输入版本编号"
    versionSelectorDiv.appendChild(versionSelectorLabel)
    let versionInput = document.createElement('input')
    versionInput.value = obj['index']
    versionSelectorDiv.appendChild(versionInput)
    let buttonView = document.createElement('button')
    buttonView.textContent = "查看"
    buttonView.onclick = function () {
        window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}&index=${versionInput.value}`
    }
    versionSelectorDiv.appendChild(buttonView)
    let buttonRollback = document.createElement('button')
    buttonRollback.textContent = "回滚"
    versionSelectorDiv.appendChild(buttonRollback)
    versionInput.onchange = function () {
        buttonRollback.hidden = true
    }
    buttonRollback.onclick = function () {
        console.log("buttonRollback", env, versionInput.value)
        url = `http://${host}/config/proxy?env=${env}&action=rollback`

        fetch(url + `&index=${versionInput.value}`, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }
        })
            .then(resp => resp.json())
            .then(json => {
                if (json.success === true) {
                    window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}`
                } else {
                    alert(json.message)
                }
            })
    }

    // create the editor
    const container = document.getElementById("jsoneditor")
    const options = {
        mode: 'text',
        modes: ['text', 'code'],
        onEditable: function (node) {
            if (!node.path) {
                // In modes code and text, node is empty: no path, field, or value
                // returning false makes the text area read-only
                return false;
            }
        },
        onModeChange: function (newMode, oldMode) {
            console.log('Mode switched from', oldMode, 'to', newMode)
        }}
    const editor = new JSONEditor(container, options)

    // set json
    editor.set(obj)

}

function processRight(ss, curIndex, env, host) {
    let obj = JSON.parse(ss);

    let keySelectorDiv = document.getElementById("keyselector")
    let keySelectorLabel = document.createElement("label")
    keySelectorLabel.textContent = "输入key"
    keySelectorDiv.appendChild(keySelectorLabel)
    let keyInput = document.createElement("input")
    keySelectorDiv.appendChild(keyInput)

    let buttonView = document.createElement("button")
    buttonView.textContent = "原始数据"
    keySelectorDiv.appendChild(buttonView)
    let buttonUpdate = document.createElement("button")
    buttonUpdate.textContent = "更新"
    keySelectorDiv.appendChild(buttonUpdate)
    let buttonDelete = document.createElement("button")
    buttonDelete.textContent = "删除key"
    keySelectorDiv.appendChild(buttonDelete)
    if (parseInt(curIndex) !== obj['index']) {
        buttonUpdate.disabled = true
        buttonDelete.disabled = true
    }

    // create the editor
    const container = document.getElementById("jsoneditorright")
    const options = {
        mode: 'text',
        modes: ['text', 'code'],
        onModeChange: function (newMode, oldMode) {
            console.log('Mode switched from', oldMode, 'to', newMode)
        }}
    const editor = new JSONEditor(container, options)

    buttonUpdate.onclick = function () {
        if (!keyInput.value.length) {
            alert("invalid key")
            return
        }
        let key = keyInput.value.trim()
        console.log("buttonUpdate", env, curIndex, key, editor.get(), JSON.stringify(editor.get()))
        url = `http://${host}/config/proxy?env=${env}&action=update`

        fetch(url, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"key": key, "value": editor.get()})
        })
            .then(resp => resp.json())
            .then(json => {
                if (json.success === true) {
                    window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}`
                } else {
                    alert(json.message)
                }
            })
    }

    buttonDelete.onclick = function () {
        if (!keyInput.value.length) {
            alert("invalid key")
            return
        }
        let key = keyInput.value.trim()
        console.log("buttonDelete", env, curIndex, key)
        url = `http://${host}/config/proxy?env=${env}&action=delete`

        fetch(url, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"key": key, "value": null})
        })
            .then(resp => resp.json())
            .then(json => {
                if (json.success === true) {
                    window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}`
                } else {
                    alert(json.message)
                }
            })
    }

    buttonView.onclick = function () {
        let key = keyInput.value.trim()
        editor.set(obj['data'][key])
    }

}


