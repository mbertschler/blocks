package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evanw/esbuild/pkg/cli"
)

// TODO:
// -[x] esbuild integration for more complex JS and importing
// -[ ] fake page switching

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
	err := buildBrowserAssets()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp()
	counter := &Counter{App: app}
	app.Server.RegisterComponent(counter)
	app.Server.RegisterPage("/counter", counter.RenderPage)

	registerTodoList(app.Server, app.DB)

	app.Server.Static("/dist/", "./dist")

	log.Println("listening on localhost:8000")
	err = http.ListenAndServe("localhost:8000", app.Server.Handler())
	if err != nil {
		log.Fatal(err)
	}
}

func buildBrowserAssets() error {
	log.Println("building browser assets")
	options := []string{
		"js/main.js",
		"--bundle",
		"--outfile=dist/bundle.js",
		"--minify",
		"--sourcemap",
	}
	returnCode := cli.Run(options)
	if returnCode != 0 {
		return fmt.Errorf("esbuild failed with code %d", returnCode)
	}
	return nil
}
