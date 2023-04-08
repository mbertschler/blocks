"use strict";

const EnterKey = "Enter"
const EscapeKey = "Escape"

function newTodoKeydown(event) {
    if (event.key != EnterKey) {
        return
    }
    guiapi("TodoList.NewTodo", { text: event.target.value })
}
callableFunctions.newTodoKeydown = newTodoKeydown

function initEdit(element, args) {
    let stoppedEditing = false
    element.addEventListener("blur", function (event) {
        if (stoppedEditing) {
            return
        }
        guiapi("TodoList.UpdateItem", { id: args.id, text: event.target.value, page: args.page });
        return false;
    })
    element.addEventListener("keydown", function (event) {
        if (event.key == EscapeKey) {
            stoppedEditing = true
            guiapi("TodoList.EditItem", { id: -1, page: args.page })
            return false;
        }
        if (event.key != EnterKey) {
            return false;
        }
        guiapi("TodoList.UpdateItem", { id: args.id, text: event.target.value, page: args.page });
    })
    element.focus()
}
callableFunctions.initEdit = initEdit

setupGuiapi()