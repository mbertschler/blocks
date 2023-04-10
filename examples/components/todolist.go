package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
)

// started work based on the template at 13:45
// done implementing the static template at 14:05
// done with NewTodo and first DB integration at 14:50
// done with item toggling and deletion at 15:05
// done with all active and completed filter pages 15:30
// ---- long break ----
// start again at 19:10
// done with ToggleAll at 19:17
// done with item editing 20:35 (half hour was spent debugging the css)

// 1.75 hours in first session (13:45 - 15:30)
// 1.25 hours in second session (19:10 - 20:35)
// in total 3 hours

func todoLayout(todoApp html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		html.Html(html.Attr("lang", "en"),
			html.Head(nil,
				html.Meta(html.Charset("utf-8")),
				html.Meta(html.Name("viewport").Content("width=device-width, initial-scale=1")),
				html.Title(nil, html.Text("Guiapi • TodoMVC")),
				html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/todomvc-app-css@2.4.2/index.min.css")),
				html.Link(html.Rel("stylesheet").Href("/css/main.css")),
			),
			html.Body(nil,
				todoApp,
				html.Elem("footer", html.Class("info"),
					html.P(nil, html.Text("Double-click to edit a todo")),
					html.P(nil, html.Text("Template by "), html.A(html.Href("http://sindresorhus.com"), html.Text("Sindre Sorhus"))),
					html.P(nil, html.Text("Created by "), html.A(html.Href("https://github.com/mbertschler"), html.Text("Martin Bertschler"))),
					html.P(nil, html.Text("Part of "), html.A(html.Href("http://todomvc.com"), html.Text("TodoMVC"))),
				),
				html.Script(html.Src("/js/guiapi.js")),
				html.Script(html.Src("/js/todo.js")),
			),
		),
	}
}

type TodoList struct {
	*App
}

func (t *TodoList) Component() *ComponentConfig {
	return &ComponentConfig{
		Name: "TodoList",
		Actions: map[string]ActionFunc{
			"NewTodo":        t.NewTodo,
			"ToggleItem":     t.ToggleItem,
			"ToggleAll":      t.ToggleAll,
			"DeleteItem":     t.DeleteItem,
			"ClearCompleted": t.ClearCompleted,
			"EditItem":       t.EditItem,
			"UpdateItem":     t.UpdateItem,
		},
	}
}

const (
	TodoListPageAll       = "all"
	TodoListPageActive    = "active"
	TodoListPageCompleted = "completed"
)

func (t *TodoList) RenderAll(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageAll, -1)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) RenderActive(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageActive, -1)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) RenderCompleted(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageCompleted, -1)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) renderAppBlock(ctx *gin.Context, page string, editItemID int) (html.Block, error) {
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}

	var main, footer html.Block
	if len(todos.Items) > 0 {
		main, err = t.renderMainBlock(todos, page, editItemID)
		if err != nil {
			return nil, err
		}
		footer, err = t.renderFooterBlock(todos, page)
		if err != nil {
			return nil, err
		}
	}

	block := html.Elem("section", html.Class("todoapp"),
		html.Elem("header", html.Class("header"),
			html.H1(nil, html.Text("todos")),
			html.Input(html.Class("new-todo ga").Attr("placeholder", "What needs to be done?").
				Attr("autofocus", "").Attr("ga-on", "keydown").Attr("ga-func", "newTodoKeydown")),
		),
		main,
		footer,
	)
	return block, nil
}

func (t *TodoList) renderMainBlock(todos *StoredTodo, page string, editItemID int) (html.Block, error) {
	items := html.Blocks{}
	for _, item := range todos.Items {
		if page == TodoListPageActive && item.Done {
			continue
		}
		if page == TodoListPageCompleted && !item.Done {
			continue
		}
		items.Add(t.renderItem(&item, page, editItemID))
	}
	main := html.Elem("section", html.Class("main"),
		html.Input(html.Class("toggle-all").Attr("type", "checkbox")),
		html.Label(html.Class("ga").Attr("for", "toggle-all").Attr("ga-on", "click").Attr("ga-action", "TodoList.ToggleAll").
			Attr("ga-args", fmt.Sprintf(`{"page":%q}`, page)), html.Text("Mark all as complete")),
		html.Ul(html.Class("todo-list"),
			items,
		),
	)
	return main, nil
}

func (t *TodoList) renderItem(item *StoredTodoItem, page string, editItemID int) html.Block {
	if item.ID == editItemID {
		return t.renderItemEdit(item, page, editItemID)
	}

	liAttrs := html.Attr("ga-on", "dblclick").
		Attr("ga-action", "TodoList.EditItem").
		Attr("ga-args", fmt.Sprintf(`{"id":%d,"page":%q}`, item.ID, page))
	inputAttrs := html.Class("toggle ga").Attr("type", "checkbox").
		Attr("ga-on", "click").Attr("ga-action", "TodoList.ToggleItem").
		Attr("ga-args", fmt.Sprintf(`{"id":%d,"page":%q}`, item.ID, page))
	if item.Done {
		liAttrs = liAttrs.Class("completed ga")
		inputAttrs = inputAttrs.Attr("checked", "")
	} else {
		liAttrs = liAttrs.Class("active ga")
	}

	li := html.Li(liAttrs,
		html.Div(html.Class("view"),
			html.Input(inputAttrs),
			html.Label(nil, html.Text(item.Text)),
			html.Button(html.Class("destroy ga").
				Attr("ga-on", "click").Attr("ga-action", "TodoList.DeleteItem").
				Attr("ga-args", fmt.Sprintf(`{"id":%d,"page":%q}`, item.ID, page))),
		),
	)
	return li
}

func (t *TodoList) renderItemEdit(item *StoredTodoItem, page string, editItemID int) html.Block {
	li := html.Li(html.Class("editing"),
		html.Div(html.Class("view"),
			html.Input(html.Class("edit ga").Attr("ga-init", "initEdit").
				Attr("ga-args", fmt.Sprintf(`{"id":%d, "page":%q}`, item.ID, page)).Attr("value", item.Text)),
		),
	)
	return li
}

func (t *TodoList) renderFooterBlock(todos *StoredTodo, page string) (html.Block, error) {
	var allClass, activeClass, completedClass string
	switch page {
	case TodoListPageAll:
		allClass = "selected"
	case TodoListPageActive:
		activeClass = "selected"
	case TodoListPageCompleted:
		completedClass = "selected"
	default:
		allClass = "selected"
	}

	leftCount := 0
	someDone := false
	for _, item := range todos.Items {
		if !item.Done {
			leftCount++
		} else {
			someDone = true
		}
	}
	itemsLeftText := " items left"
	if leftCount == 1 {
		itemsLeftText = " item left"
	}

	var clearCompletedButton html.Block
	if someDone {
		clearCompletedButton = html.Button(html.Class("clear-completed ga").Attr("ga-on", "click").Attr("ga-action", "TodoList.ClearCompleted").
			Attr("ga-args", fmt.Sprintf(`{"page":%q}`, page)), html.Text("Clear completed"))
	}

	footer := html.Elem("footer", html.Class("footer"),
		html.Span(html.Class("todo-count"),
			html.Strong(nil, html.Text(fmt.Sprint(leftCount))),
			html.Text(itemsLeftText),
		),
		html.Ul(html.Class("filters"),
			html.Li(nil,
				html.A(html.Class(allClass).Href("/"), html.Text("All")),
			),
			html.Li(nil,
				html.A(html.Class(activeClass).Href("/active"), html.Text("Active")),
			),
			html.Li(nil,
				html.A(html.Class(completedClass).Href("/completed"), html.Text("Completed")),
			),
		),
		clearCompletedButton,
	)
	return footer, nil
}

func (t *TodoList) NewTodo(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		Text string `json:"text"`
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		return nil, err
	}
	var highestID int
	for _, item := range todos.Items {
		if item.ID > highestID {
			highestID = item.ID
		}
	}
	input.Text = strings.TrimSpace(input.Text)
	todos.Items = append(todos.Items, StoredTodoItem{ID: highestID + 1, Text: input.Text})

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) ToggleItem(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		ID   int    `json:"id"`
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		log.Println("error unmarshaling args", string(args))
		return nil, err
	}

	for i, item := range todos.Items {
		if item.ID == input.ID {
			todos.Items[i].Done = !todos.Items[i].Done
		}
	}

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) ToggleAll(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		log.Println("error unmarshaling args", string(args))
		return nil, err
	}

	allDone := true
	for _, item := range todos.Items {
		if !item.Done {
			allDone = false
			break
		}
	}

	for i := range todos.Items {
		todos.Items[i].Done = !allDone
	}

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) DeleteItem(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		ID   int    `json:"id"`
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		return nil, err
	}

	var newItems []StoredTodoItem
	for _, item := range todos.Items {
		if item.ID == input.ID {
			continue
		}
		newItems = append(newItems, item)
	}
	todos.Items = newItems

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) ClearCompleted(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		return nil, err
	}

	var newItems []StoredTodoItem
	for _, item := range todos.Items {
		if item.Done {
			continue
		}
		newItems = append(newItems, item)
	}
	todos.Items = newItems

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) EditItem(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		ID   int    `json:"id"`
		Page string `json:"page"`
	}
	var input In
	err := json.Unmarshal(args, &input)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, input.ID)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) UpdateItem(ctx *gin.Context, args json.RawMessage) (*Response, error) {
	type In struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
		Page string `json:"page"`
	}
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}
	var input In
	err = json.Unmarshal(args, &input)
	if err != nil {
		log.Println("error unmarshaling args", string(args))
		return nil, err
	}

	for i, item := range todos.Items {
		if item.ID == input.ID {
			todos.Items[i].Text = strings.TrimSpace(input.Text)
		}
	}

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page, -1)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}
