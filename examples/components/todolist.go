package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
	"github.com/mbertschler/blocks/html/attr"
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

// Code stats before refactoring (excluding comments and imports):
// todolist.go: 473 lines
// todolist.js: 38 lines
// total: 511 lines

// Code stats after refactoring (excluding comments and imports):
// todolist.go: 405 lines -15%
// todolist.js: 38 lines
// total: 443 lines -14%

func registerTodoList(server *Server, db *DB) {
	tl := &TodoList{DB: db}
	server.RegisterComponent(tl)
	server.RegisterPage("/", tl.RenderFullPage(TodoListPageAll))
	server.RegisterPage("/active", tl.RenderFullPage(TodoListPageActive))
	server.RegisterPage("/completed", tl.RenderFullPage(TodoListPageCompleted))
}

type Context struct {
	gin   *gin.Context
	Sess  *Session
	State TodoListState
}

type TypedContextCallable[T any] func(c *Context, args *T) (*Response, error)

func ContextCallable[T any](fn TypedContextCallable[T]) Callable {
	return func(c *gin.Context, raw json.RawMessage) (*Response, error) {
		var input T
		err := json.Unmarshal(raw, &input)
		if err != nil {
			return nil, err
		}

		ctx := &Context{gin: c,
			Sess: sessionFromContext(c)}

		stateJSON, ok := c.Keys["rawState"].([]byte)
		if ok {
			err = json.Unmarshal(stateJSON, &ctx.State)
			if err != nil {
				return nil, err
			}
		}
		return fn(ctx, &input)
	}
}

type TypedContextPage func(c *Context) (html.Block, error)

func ContextPage(fn TypedContextPage) PageFunc {
	return func(c *gin.Context) (html.Block, error) {
		sess := sessionFromContext(c)
		return fn(&Context{gin: c, Sess: sess})
	}
}

func todoLayout(todoApp html.Block, state TodoListState) (html.Block, error) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}
	return html.Blocks{
		html.Doctype("html"),
		html.Html(attr.Attr("lang", "en"),
			html.Head(nil,
				html.Meta(attr.Charset("utf-8")),
				html.Meta(attr.Name("viewport").Content("width=device-width, initial-scale=1")),
				html.Title(nil, html.Text("Guiapi • TodoMVC")),
				html.Link(attr.Rel("stylesheet").Href("https://cdn.jsdelivr.net/npm/todomvc-app-css@2.4.2/index.min.css")),
				html.Link(attr.Rel("stylesheet").Href("/dist/bundle.css")),
			),
			html.Body(nil,
				todoApp,
				html.Elem("footer", attr.Class("info"),
					html.P(nil, html.Text("Double-click to edit a todo")),
					html.P(nil, html.Text("Template by "), html.A(attr.Href("http://sindresorhus.com"), html.Text("Sindre Sorhus"))),
					html.P(nil, html.Text("Created by "), html.A(attr.Href("https://github.com/mbertschler"), html.Text("Martin Bertschler"))),
					html.P(nil, html.Text("Part of "), html.A(attr.Href("http://todomvc.com"), html.Text("TodoMVC"))),
					html.A(attr.Href("/counter"), html.Text("Counter Example")),
				),
				html.Script(nil, html.JS("var state = "+string(stateJSON)+";")),
				html.Script(attr.Src("/dist/bundle.js")),
			),
		),
	}, nil
}

type TodoList struct {
	*DB
}

func (t *TodoList) Component() *ComponentConfig {
	return &ComponentConfig{
		Name: "TodoList",
		Actions: map[string]Callable{
			"NewTodo":        ContextCallable(t.NewTodo),
			"ToggleItem":     ContextCallable(t.ToggleItem),
			"ToggleAll":      ContextCallable(t.ToggleAll),
			"DeleteItem":     ContextCallable(t.DeleteItem),
			"ClearCompleted": ContextCallable(t.ClearCompleted),
			"EditItem":       ContextCallable(t.EditItem),
			"UpdateItem":     ContextCallable(t.UpdateItem),
			"Page":           ContextCallable(t.Page),
		},
	}
}

const (
	TodoListPageAll       = "all"
	TodoListPageActive    = "active"
	TodoListPageCompleted = "completed"
)

type TodoListState struct {
	Page string
}

type TodoListProps struct {
	Page       string
	Todos      *StoredTodo
	EditItemID int
}

func (t *TodoList) RenderFullPage(page string) PageFunc {
	return ContextPage(func(ctx *Context) (html.Block, error) {
		content, err := t.renderPageContent(ctx, page)
		if err != nil {
			return nil, err
		}
		return todoLayout(content, ctx.State)
	})
}

func (t *TodoList) renderPageContent(ctx *Context, page string) (html.Block, error) {
	ctx.State.Page = page
	props, err := t.todoListProps(ctx)
	if err != nil {
		return nil, err
	}
	return t.renderBlock(props)
}

type PageArgs struct {
	Page string
}

func (t *TodoList) Page(ctx *Context, args *PageArgs) (*Response, error) {
	block, err := t.renderPageContent(ctx, args.Page)
	if err != nil {
		return nil, err
	}
	res, err := ReplaceContent(".todoapp", block)
	if err != nil {
		return nil, err
	}
	res.State = ctx.State
	return res, nil
}

func (t *TodoList) renderBlock(props *TodoListProps) (html.Block, error) {
	var main, footer html.Block
	if len(props.Todos.Items) > 0 {
		var err error
		main, err = t.renderMainBlock(props.Todos, props.Page, props.EditItemID)
		if err != nil {
			return nil, err
		}
		footer, err = t.renderFooterBlock(props.Todos, props.Page)
		if err != nil {
			return nil, err
		}
	}

	block := html.Elem("section", attr.Class("todoapp"),
		html.Elem("header", attr.Class("header"),
			html.H1(nil, html.Text("todos")),
			html.Input(attr.Class("new-todo ga").Attr("placeholder", "What needs to be done?").
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
		items.Add(t.renderItem(&item, editItemID))
	}
	main := html.Elem("section", attr.Class("main"),
		html.Input(attr.Class("toggle-all").Attr("type", "checkbox")),
		html.Label(attr.Class("ga").Attr("for", "toggle-all").Attr("ga-on", "click").Attr("ga-action", "TodoList.ToggleAll"),
			html.Text("Mark all as complete")),
		html.Ul(attr.Class("todo-list"),
			items,
		),
	)
	return main, nil
}

func (t *TodoList) renderItem(item *StoredTodoItem, editItemID int) html.Block {
	if item.ID == editItemID {
		return t.renderItemEdit(item, editItemID)
	}

	liAttrs := attr.Attr("ga-on", "dblclick").
		Attr("ga-action", "TodoList.EditItem").
		Attr("ga-args", fmt.Sprintf(`{"id":%d}`, item.ID))
	inputAttrs := attr.Class("toggle ga").Attr("type", "checkbox").
		Attr("ga-on", "click").Attr("ga-action", "TodoList.ToggleItem").
		Attr("ga-args", fmt.Sprintf(`{"id":%d}`, item.ID))
	if item.Done {
		liAttrs = liAttrs.Class("completed ga")
		inputAttrs = inputAttrs.Attr("checked", "")
	} else {
		liAttrs = liAttrs.Class("active ga")
	}

	li := html.Li(liAttrs,
		html.Div(attr.Class("view"),
			html.Input(inputAttrs),
			html.Label(nil, html.Text(item.Text)),
			html.Button(attr.Class("destroy ga").
				Attr("ga-on", "click").Attr("ga-action", "TodoList.DeleteItem").
				Attr("ga-args", fmt.Sprintf(`{"id":%d}`, item.ID))),
		),
	)
	return li
}

func (t *TodoList) renderItemEdit(item *StoredTodoItem, editItemID int) html.Block {
	li := html.Li(attr.Class("editing"),
		html.Div(attr.Class("view"),
			html.Input(attr.Class("edit ga").Attr("ga-init", "initEdit").
				Attr("ga-args", fmt.Sprintf(`{"id":%d}`, item.ID)).Attr("value", item.Text)),
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
		clearCompletedButton = html.Button(attr.Class("clear-completed ga").Attr("ga-on", "click").Attr("ga-action", "TodoList.ClearCompleted"),
			html.Text("Clear completed"))
	}

	footer := html.Elem("footer", attr.Class("footer"),
		html.Span(attr.Class("todo-count"),
			html.Strong(nil, html.Text(fmt.Sprint(leftCount))),
			html.Text(itemsLeftText),
		),
		html.Ul(attr.Class("filters"),
			html.Li(nil,
				html.A(attr.Class(allClass+" ga").Href("/").Attr("ga-link", "TodoList.Page").Attr("ga-state", `{"Page":"all"}`), html.Text("All")),
			),
			html.Li(nil,
				html.A(attr.Class(activeClass+" ga").Href("/active").Attr("ga-link", "TodoList.Page").Attr("ga-state", `{"Page":"active"}`), html.Text("Active")),
			),
			html.Li(nil,
				html.A(attr.Class(completedClass+" ga").Href("/completed").Attr("ga-link", "TodoList.Page").Attr("ga-state", `{"Page":"completed"}`), html.Text("Completed")),
			),
		),
		clearCompletedButton,
	)
	return footer, nil
}

type NewTodoArgs struct {
	Text string `json:"text"`
}

func (t *TodoList) NewTodo(ctx *Context, input *NewTodoArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
		var highestID int
		for _, item := range todos.Items {
			if item.ID > highestID {
				highestID = item.ID
			}
		}
		input.Text = strings.TrimSpace(input.Text)
		todos.Items = append(todos.Items, StoredTodoItem{ID: highestID + 1, Text: input.Text})
		return t.DB.SetTodo(todos)
	})
}

func (t *TodoList) ToggleItem(ctx *Context, args *IDArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
		for i, item := range todos.Items {
			if item.ID == args.ID {
				todos.Items[i].Done = !todos.Items[i].Done
			}
		}
		return t.DB.SetTodo(todos)
	})
}

func (t *TodoList) ToggleAll(ctx *Context, args *NoArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
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
		return t.DB.SetTodo(todos)
	})
}

func (t *TodoList) DeleteItem(ctx *Context, args *IDArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
		var newItems []StoredTodoItem
		for _, item := range todos.Items {
			if item.ID == args.ID {
				continue
			}
			newItems = append(newItems, item)
		}
		todos.Items = newItems
		return t.DB.SetTodo(todos)
	})
}

type NoArgs struct{}

func (t *TodoList) ClearCompleted(ctx *Context, _ *NoArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
		var newItems []StoredTodoItem
		for _, item := range todos.Items {
			if item.Done {
				continue
			}
			newItems = append(newItems, item)
		}
		todos.Items = newItems
		return t.DB.SetTodo(todos)
	})
}

type IDArgs struct {
	ID int `json:"id"`
}

func (t *TodoList) EditItem(ctx *Context, args *IDArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, _ *StoredTodo) error {
		props.EditItemID = args.ID
		return nil
	})
}

type UpdateItemArgs struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (t *TodoList) UpdateItem(ctx *Context, args *UpdateItemArgs) (*Response, error) {
	return t.updateTodoList(ctx, func(props *TodoListProps, todos *StoredTodo) error {
		for i, item := range todos.Items {
			if item.ID == args.ID {
				todos.Items[i].Text = strings.TrimSpace(args.Text)
			}
		}
		return t.DB.SetTodo(todos)
	})
}

func (t *TodoList) updateTodoList(ctx *Context, fn func(*TodoListProps, *StoredTodo) error) (*Response, error) {
	props, err := t.todoListProps(ctx)
	if err != nil {
		return nil, err
	}

	err = fn(props, props.Todos)
	if err != nil {
		return nil, err
	}

	appBlock, err := t.renderBlock(props)
	if err != nil {
		return nil, err
	}
	return ReplaceContent(".todoapp", appBlock)
}

func (t *TodoList) todoListProps(ctx *Context) (*TodoListProps, error) {
	todos, err := t.DB.GetTodo(ctx.Sess.ID)
	if err != nil {
		return nil, err
	}

	return &TodoListProps{
		Page:  ctx.State.Page,
		Todos: todos,
	}, nil
}
