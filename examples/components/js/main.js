import "../css/main.css"

import TodoList from "./todolist.js"
import { registerFunctions, setupGuiapi } from "./guiapi"

registerFunctions(TodoList)
setupGuiapi({ debug: true })
