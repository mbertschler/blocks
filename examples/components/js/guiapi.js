var callableFunctions = {}

const debugGuiapi = false

function guiapi(name, args, callback) {
    if (debugGuiapi) {
        console.log("guiapi", name, args)
    }
    if (!callback) {
        callback = () => { }
    }
    var req = {
        Name: name,
        Args: args,
    };
    fetch("/guiapi", {
        method: 'POST',
        mode: 'cors',
        credentials: 'same-origin',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(req)
    }).then((response) => {
        response.json().then(r => handleResponse(r, callback))
            .catch(r => console.error('response.json() error:', r))
    }).catch((reason) => {
        console.error('fetch() error:', reason)
        callback(reason)
    })
}

function handleResponse(r, callback) {
    if (r.Error) {
        console.error("[" + r.Error.Code + "]", r.Error.Message, r.Error);
        window.alert("guiapi error, check console");
        callback(r.Error)
        return;
    }
    if (r.HTML) {
        for (var j = 0; j < r.HTML.length; j++) {
            var update = r.HTML[j];
            const el = document.querySelector(update.Selector);
            if (!el) {
                console.warn("update selector not found :(", update.Selector, update);
                continue;
            }

            switch (update.Operation) {
                case 1:
                    el.innerHTML = update.Content;
                    break;
                case 2:
                    el.outerHTML = update.Content;
                    break;
                case 3:
                    el.insertAdjacentHTML('beforebegin', update.Content);
                    break;
                default:
                    console.warn("update type not implemented :(", update);
            }
        }
    }
    if (r.JS) {
        for (var j = 0; j < r.JS.length; j++) {
            var call = r.JS[j];
            var func = callableFunctions[call.Name];
            if (func) {
                func(call.Args);
            } else {
                console.warn("function call not implemented :(", call);
            }
        }
    }
    hydrate()
    callback(null)
}

function hydrate() {
    var elements = document.querySelectorAll(".ga")
    for (var el of elements) {
        if (el.dataset.action) {
            hydrateAction(el)
        }
        if (el.attributes.getNamedItem("ga-click")) {
            hydrateClick(el)
        }
    }
}

function hydrateAction(el) {
    var action = el.dataset.action
    var args = null
    var variable = el.dataset.var
    if (variable) {
        args = window[variable]
    }
    var fn = callableFunctions[action];
    if (fn) {
        fn(el, args);
    } else {
        console.warn("function call not implemented :(", action, el);
        console.log("during hydrating", elements.length, "elements", elements)
    }
    el.classList.remove("ga")
}

function hydrateClick(el) {
    var action = el.attributes.getNamedItem("ga-click").value
    var args = null
    if (el.attributes.getNamedItem("ga-args")) {
        args = el.attributes.getNamedItem("ga-args").value
    }
    el.addEventListener("click", function (e) {
        guiapi(action, args);
        e.preventDefault();
        e.stopPropagation();
        return false;
    })
    el.classList.remove("ga")
}

hydrate()
