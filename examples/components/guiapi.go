package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbertschler/blocks/html"

	"github.com/gin-gonic/gin"
)

// NewGuiapi returns an empty handler
func NewGuiapi() *Handler {
	return &Handler{
		Functions: map[string]Callable{},
	}
}

// SetFunc sets a callable GUI API function in the handler.
func (h *Handler) SetFunc(name string, fn Callable) {
	h.Functions[name] = fn
}

// ReplaceContent is a helper function that returns a Result that replaces
// the element content chosen by the selector with the passed Block.
func ReplaceContent(selector string, block html.Block) (*Response, error) {
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := &Response{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLReplaceContent,
				Selector:  selector,
				Content:   out,
			},
		},
	}
	return ret, nil
}

// ReplaceElement is a helper function that returns a Result that
// replaces the whole element chosen by the selector with the passed Block.
func ReplaceElement(selector string, block html.Block) (*Response, error) {
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := &Response{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLReplaceElement,
				Selector:  selector,
				Content:   out,
			},
		},
	}
	return ret, nil
}

// InsertBefore is a helper function that returns a Result that
// inserts a block on the same level before the passed selector.
func InsertBefore(selector string, block html.Block) (*Response, error) {
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := &Response{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLInsertBefore,
				Selector:  selector,
				Content:   out,
			},
		},
	}
	return ret, nil
}

// Redirect lets the browser navigate to a given path
func Redirect(path string) (*Response, error) {
	ret := &Response{
		JS: []JSCall{
			{
				Name: "redirect",
				Args: path,
			},
		},
	}
	return ret, nil
}

// Handle handles HTTP requests to the GUI API.
func (h *Handler) Handle(c *gin.Context) {
	var req Request
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	resp := h.process(c, &req)
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) process(c *gin.Context, req *Request) *Response {
	var res = Response{
		ID:   req.ID,
		Name: req.Name,
	}
	fn, ok := h.Functions[req.Name]
	if !ok {
		res.Error = &Error{
			Code:    "undefinedFunction",
			Message: fmt.Sprint(req.Name, " is not defined"),
		}

	} else {
		r, err := fn(c, req.Args)
		if err != nil {
			res.Error = &Error{
				Code:    "error",
				Message: err.Error(),
			}
		}
		if r != nil {
			res.HTML = r.HTML
			res.JS = r.JS
		}
	}
	return &res
}

// Request is the sent body of a GUI API call
type Request struct {
	ID   int    `json:",omitempty"` // ID can be used from the client to identify responses
	Name string // Name of the action that is called
	// Args as object, gets parsed by the called function
	Args json.RawMessage `json:",omitempty"`
}

type Handler struct {
	Functions map[string]Callable
}

type Callable func(c *gin.Context, args json.RawMessage) (*Response, error)

// Response is the returned body of a GUI API call
type Response struct {
	ID    int          `json:",omitempty"` // ID from the calling action is returned
	Name  string       // Name of the action that was called
	Error *Error       `json:",omitempty"`
	HTML  []HTMLUpdate `json:",omitempty"` // DOM updates to apply
	JS    []JSCall     `json:",omitempty"` // JS calls to execute
}

type Error struct {
	Code    string
	Message string
}

type HTMLOp int8

const (
	HTMLReplaceContent HTMLOp = 1
	HTMLReplaceElement HTMLOp = 2
	HTMLInsertBefore   HTMLOp = 3
)

type HTMLUpdate struct {
	Operation HTMLOp // how to apply this update
	Selector  string // jQuery style selector: #id .class
	Content   string `json:",omitempty"` // inner HTML
	// Init calls are executed after the HTML is added
	Init []JSCall `json:",omitempty"`
	// Destroy calls are executed before the HTML is removed
	Destroy []JSCall `json:",omitempty"`
}

type JSCall struct {
	Name string // name of the function to call
	// Args as object, gets encoded by the called function
	Args interface{} `json:",omitempty"`
}

func (r *Response) AddJSResponse(name string, args interface{}) {
	r.JS = append(r.JS, JSCall{Name: name, Args: args})
}
