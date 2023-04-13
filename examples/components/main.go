package main

import (
	"log"
	"net/http"
)

// TODO:
// -[x] cookie based session
// -[x] guiapi
// -[x] make the storage swappable
// -[x] counter demo
// -[x] turn it into a reusable component
// -[ ] extract page registration from component
// -[ ] add a second page
// -[ ] TODO MVC example?
// -[ ] https://github.com/tastejs/todomvc/blob/master/app-spec.md

// maybe later:
// -[ ] offer cockroachdb as storage
// -[ ] text field for your name
// -[ ] esbuild integration for more complex JS and importing

type App struct {
	DB     *DB
	Server *Server
}

func NewApp() *App {
	app := &App{}
	app.DB = NewDB()
	app.Server = NewServer(app.DB)
	return app
}

func main() {
	app := NewApp()
	counter := &Counter{App: app}
	app.Server.RegisterComponent(counter)
	app.Server.RegisterPage("/counter", counter.RenderPage)

	registerTodoList(app.Server, app.DB)

	app.Server.Static("/js/", "./js")
	app.Server.Static("/css/", "./css")

	log.Println("listening on localhost:8000")
	err := http.ListenAndServe("localhost:8000", app.Server.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
