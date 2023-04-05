"use strict";

const EnterKey = "Enter"

function newTodoKeydown(event) {
    if (event.key != EnterKey) {
        return
    }
    guiapi("TodoList.NewTodo", { text: event.target.value })
}
callableFunctions.newTodoKeydown = newTodoKeydown

setupGuiapi()