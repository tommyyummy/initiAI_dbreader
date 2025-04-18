
function processLeft(ss, curIndex, env, realSS, host) {
    let obj = JSON.parse(ss);
    let realObj = JSON.parse(realSS);

    let versionDiv = document.getElementById("versionleft")
    let version = document.createElement('pre')
    version.textContent = `总配置最新版本编号：${curIndex}`
    versionDiv.appendChild(version)
    let refreshButton = document.createElement("button")
    versionDiv.appendChild(refreshButton)
    refreshButton.textContent = "刷新"
    refreshButton.addEventListener("click", () => {
        console.log("refreshButton")
        window.location.href = window.location.protocol + "//" + window.location.host + window.location.pathname + `?env=${env}`
    })
    refreshButton.style.height = "50%"
    refreshButton.style.alignSelf = "center"
    let checkFullButton = document.createElement("button")
    checkFullButton.textContent = "查看全量配置"
    checkFullButton.isFull = false
    checkFullButton.addEventListener("click", () => {
        console.log("checkFullButton")
        if(checkFullButton.isFull) {
            editor.set(obj)
            checkFullButton.isFull = false
            checkFullButton.textContent = "查看全量配置"
        } else {
            editor.set(realObj)
            checkFullButton.isFull = true
            checkFullButton.textContent = "收起全量配置"
        }
    })
    versionDiv.appendChild(checkFullButton)


    // create the editor
    const container = document.getElementById("jsoneditorleft")
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

let namespaceData = {}


function processMid(ss, curIndex, env, host) {
    let obj = JSON.parse(ss);

    let namespaceSelectorDiv = document.getElementById("namespacemid")
    let namespaceSelectorLabel = document.createElement("namespacelabel")
    namespaceSelectorLabel.textContent = "输入namespace"
    namespaceSelectorDiv.appendChild(namespaceSelectorLabel)
    let namespaceInput = document.createElement("input")
    namespaceSelectorDiv.appendChild(namespaceInput)
    namespaceInput.onchange = function () {
        buttonView.hidden = true
        buttonRollback.hidden = true
    }
    let namespaceCheckButton = document.createElement("button")
    namespaceCheckButton.textContent = "查看"
    namespaceSelectorDiv.appendChild(namespaceCheckButton)
    let version = document.createElement('pre')
    version.textContent = `namespace最新配置版本编号：`
    namespaceSelectorDiv.appendChild(version)
    namespaceCheckButton.onclick = function () {
        console.log("namespaceCheckButton")
        let namespace = namespaceInput.value.trim()
        if (!namespace) {return}
        let url = `http://${host}/config_v2/proxy?env=${env}&action=namespace_new&namespace=${namespace}`

        fetch(url, {
            method: "GET",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }
        })
            .then(resp => resp.json())
            .then(json => {
                let index = json.index;
                version.textContent = `namespace最新配置版本编号：` + index.toString()
                versionInput.value = index
                editor.set(json)
                if (index !== obj.index_data[namespace]) {
                    if (obj.index_data[namespace] === undefined && index === 0) {
                        alert(`namespace ${namespace} 尚未被使用`)
                    } else {
                        alert(`namespace ${namespace} 当前最新配置版本编号 ${index} 与 总配置最新版本编号 ${obj.index_data[namespace]} 不符，可能有未捕捉到的修改或数据源有误。\n\r！！！请刷新页面，若还存在问题，请检查配置数据！！！`)
                    }
                }
                buttonView.hidden = false
                buttonRollback.hidden = false
                namespaceData = json
            })
    }


    let versionSelectorDiv = document.getElementById("versionseletormid")
    let versionSelectorLabel = document.createElement("label")
    versionSelectorLabel.textContent = "输入版本编号"
    versionSelectorDiv.appendChild(versionSelectorLabel)
    let versionInput = document.createElement('input')
    versionSelectorDiv.appendChild(versionInput)
    let buttonView = document.createElement('button')
    buttonView.textContent = "查看"
    buttonView.onclick = function () {
        console.log("buttonView")
        let namespace = namespaceInput.value.trim()
        let url = `http://${host}/config_v2/proxy?env=${env}&action=namespace_index&namespace=${namespace}&index=${versionInput.value}`

        fetch(url, {
            method: "GET",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            }
        })
            .then(resp => resp.json())
            .then(json => {
                versionInput.value = json.index
                editor.set(json)
                buttonRollback.hidden = false
                namespaceData = json
            })
    }
    versionSelectorDiv.appendChild(buttonView)
    let buttonRollback = document.createElement('button')
    buttonRollback.textContent = "回滚"
    versionSelectorDiv.appendChild(buttonRollback)
    versionInput.onchange = function () {
        buttonRollback.hidden = true
    }
    buttonRollback.onclick = function () {
        console.log("buttonRollback")
        let namespace = namespaceInput.value.trim()
        let url = `http://${host}/config_v2/proxy?env=${env}&action=rollback&namespace=${namespace}&index=${versionInput.value}`

        fetch(url, {
            method: "GET",
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
    const container = document.getElementById("jsoneditormid")
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
    // editor.set(obj)

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
        console.log("buttonUpdate")

        let namespace = namespaceData.namespace
        let key = keyInput.value.trim()
        let url = `http://${host}/config_v2/proxy?env=${env}&action=update`

        fetch(url, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"namespace": namespace, "key": key, "value": editor.get()})
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
        console.log("buttonDelete")

        let namespace = namespaceData.namespace
        let key = keyInput.value.trim()
        let url = `http://${host}/config_v2/proxy?env=${env}&action=delete`

        fetch(url, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"namespace": namespace, "key": key, "value": null})
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
        console.log("keyButtonView")
        let key = keyInput.value.trim()
        editor.set(namespaceData['data'][key])
    }

}


