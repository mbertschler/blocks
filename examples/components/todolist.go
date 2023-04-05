package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
)

// started work based on the template at 13:45
// done implementing the static template at 14:05
// done with NewTodo and first DB integration at 14:50
// done with item toggling and deletion at 15:05
// done with all active and completed filter pages 15:30

func todoLayout(todoApp html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		html.Html(html.Attr("lang", "en"),
			html.Head(nil,
				html.Meta(html.Charset("utf-8")),
				html.Meta(html.Name("viewport").Content("width=device-width, initial-scale=1")),
				html.Title(nil, html.Text("Guiapi • TodoMVC")),
				html.Link(html.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/todomvc-app-css@2.4.2/index.min.css")),
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
			"NewTodo":    t.NewTodo,
			"ToggleItem": t.ToggleItem,
			"DeleteItem": t.DeleteItem,
		},
	}
}

const (
	TodoListPageAll       = "all"
	TodoListPageActive    = "active"
	TodoListPageCompleted = "completed"
)

func (t *TodoList) RenderAll(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageAll)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) RenderActive(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageActive)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) RenderCompleted(ctx *gin.Context) (html.Block, error) {
	appBlock, err := t.renderAppBlock(ctx, TodoListPageCompleted)
	if err != nil {
		return nil, err
	}
	return todoLayout(appBlock), nil
}

func (t *TodoList) renderAppBlock(ctx *gin.Context, page string) (html.Block, error) {
	sess := sessionFromContext(ctx)
	todos, err := t.App.DB.GetTodo(sess.ID)
	if err != nil {
		return nil, err
	}

	var main, footer html.Block
	if len(todos.Items) > 0 {
		main, err = t.renderMainBlock(todos, page)
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

func (t *TodoList) renderMainBlock(todos *StoredTodo, page string) (html.Block, error) {
	items := html.Blocks{}
	for _, item := range todos.Items {
		if page == TodoListPageActive && item.Done {
			continue
		}
		if page == TodoListPageCompleted && !item.Done {
			continue
		}
		items.Add(t.renderItem(&item, page))
	}
	main := html.Elem("section", html.Class("main"),
		html.Input(html.Class("toggle-all").Attr("type", "checkbox")),
		html.Label(html.Attr("for", "toggle-all"), html.Text("Mark all as complete")),
		html.Ul(html.Class("todo-list"),
			items,
			// html.Li(html.Class("completed"),
			// 	html.Div(html.Class("view"),
			// 		html.Input(html.Class("toggle").Attr("type", "checkbox").Attr("checked", "")),
			// 		html.Label(nil, html.Text("Taste JavaScript")),
			// 		html.Button(html.Class("destroy")),
			// 	),
			// 	html.Input(html.Class("edit").Attr("value", "Create a TodoMVC template")),
			// ),
			// html.Li(nil,
			// 	html.Div(html.Class("view"),
			// 		html.Input(html.Class("toggle").Attr("type", "checkbox")),
			// 		html.Label(nil, html.Text("Buy a unicorn")),
			// 		html.Elem("button", html.Class("destroy")),
			// 	),
			// 	html.Input(html.Class("edit").Attr("value", "Rule the web")),
			// ),
		),
	)
	return main, nil
}

func (t *TodoList) renderItem(item *StoredTodoItem, page string) html.Block {
	liAttrs := html.Attributes{}
	inputAttrs := html.Class("toggle ga").Attr("type", "checkbox").
		Attr("ga-on", "click").Attr("ga-action", "TodoList.ToggleItem").
		Attr("ga-args", fmt.Sprintf(`{"id":%d,"page":%q}`, item.ID, page))
	if item.Done {
		liAttrs = html.Class("completed")
		inputAttrs = inputAttrs.Attr("checked", "")
	}
	li := html.Li(liAttrs,
		html.Div(html.Class("view"),
			html.Input(inputAttrs),
			html.Label(nil, html.Text(item.Text)),
			html.Button(html.Class("destroy ga").
				Attr("ga-on", "click").Attr("ga-action", "TodoList.DeleteItem").
				Attr("ga-args", fmt.Sprintf(`{"id":%d,"page":%q}`, item.ID, page))),
			html.Input(html.Class("edit").Attr("value", item.Text)),
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
	}
	footer := html.Elem("footer", html.Class("footer"),
		html.Span(html.Class("todo-count"),
			html.Strong(nil, html.Text("2")),
			html.Text(" items left"),
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
		html.Button(html.Class("clear-completed"), html.Text("Clear completed")),
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
	todos.Items = append(todos.Items, StoredTodoItem{ID: highestID + 1, Text: input.Text})

	err = t.App.DB.SetTodo(todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderAppBlock(ctx, input.Page)
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

	appBlock, err := t.renderAppBlock(ctx, input.Page)
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

	appBlock, err := t.renderAppBlock(ctx, input.Page)
	if err != nil {
		return nil, err
	}

	return ReplaceContent(".todoapp", appBlock)
}
